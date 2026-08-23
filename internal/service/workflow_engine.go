package service

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"go-waste-routes/internal/domain"
)

const (
	WorkflowPlan       = "plan"
	WorkflowRoute      = "route"
	WorkflowTask       = "task"
	WorkflowInvoice    = "invoice"
	WorkflowWeighing   = "weighing"
	WorkflowAbnormal   = "abnormality"
	WorkflowUser       = "user"
	WorkflowVehicle    = "vehicle"
)

type WorkflowEngine struct {
	flows map[string]map[string][]string
}

func NewWorkflowEngine() *WorkflowEngine {
	return &WorkflowEngine{
		flows: map[string]map[string][]string{
			WorkflowPlan: {
				domain.PlanDraft:    {domain.PlanActive, domain.PlanArchived},
				domain.PlanActive:   {domain.PlanDisabled, domain.PlanArchived},
				domain.PlanDisabled: {domain.PlanActive, domain.PlanArchived},
				domain.PlanArchived: {},
			},
			WorkflowRoute: {
				domain.RoutePending:   {domain.RouteRunning, domain.RouteCancelled, domain.RouteAbnormal},
				domain.RouteRunning:   {domain.RouteCompleted, domain.RouteCancelled, domain.RouteAbnormal},
				domain.RouteCompleted: {},
				domain.RouteCancelled: {},
				domain.RouteAbnormal:  {domain.RouteCancelled},
			},
			WorkflowTask: {
				domain.TaskPending:   {domain.TaskClaimed, domain.TaskSkipped, domain.TaskAbnormal},
				domain.TaskClaimed:   {domain.TaskRunning, domain.TaskSkipped, domain.TaskAbnormal},
				domain.TaskRunning:   {domain.TaskCompleted, domain.TaskSkipped, domain.TaskAbnormal},
				domain.TaskCompleted: {},
				domain.TaskSkipped:   {},
				domain.TaskAbnormal:  {domain.TaskCompleted, domain.TaskSkipped},
			},
			WorkflowInvoice: {
				domain.InvoicePending:    {domain.InvoiceConfirmed, domain.InvoiceWrittenOff},
				domain.InvoiceConfirmed:  {domain.InvoicePartial, domain.InvoicePaid, domain.InvoiceWrittenOff},
				domain.InvoicePartial:    {domain.InvoicePaid, domain.InvoiceWrittenOff},
				domain.InvoicePaid:       {domain.InvoiceWrittenOff},
				domain.InvoiceWrittenOff: {},
			},
			WorkflowWeighing: {
				"draft":    {"confirmed", "abnormal"},
				"confirmed": {"corrected", "abnormal"},
				"corrected": {},
				"abnormal":  {"reviewing", "rejected"},
				"reviewing": {"confirmed", "rejected"},
				"rejected":  {},
			},
			WorkflowAbnormal: {
				"pending": { "reviewing", "rejected" },
				"reviewing": { "confirmed", "rejected" },
				"confirmed": {},
				"rejected": {},
			},
			WorkflowUser: {
				domain.UserStatusEnabled:  {domain.UserStatusDisabled},
				domain.UserStatusDisabled: {domain.UserStatusEnabled},
			},
			WorkflowVehicle: {
				"available": {"repair", "scrap", "disabled"},
				"repair":     {"available", "disabled"},
				"scrap":      {},
				"disabled":   {"available"},
			},
		},
	}
}

func (e *WorkflowEngine) AllowedNext(workflow, current string) []string {
	current = NormalizeStatus(current)
	rule := e.flowFor(workflow)
	next := append([]string(nil), rule[current]...)
	sort.Strings(next)
	return next
}

func (e *WorkflowEngine) CanTransition(workflow, current, next string) bool {
	current = NormalizeStatus(current)
	next = NormalizeStatus(next)
	for _, candidate := range e.flowFor(workflow)[current] {
		if candidate == next {
			return true
		}
	}
	return false
}

func (e *WorkflowEngine) ValidateTransition(workflow, current, next string) error {
	if current == "" {
		return fmt.Errorf("current status required")
	}
	if next == "" {
		return fmt.Errorf("next status required")
	}
	if !e.CanTransition(workflow, current, next) {
		return fmt.Errorf("invalid %s transition from %s to %s", workflow, current, next)
	}
	return nil
}

func (e *WorkflowEngine) FlowFor(workflow string) map[string][]string {
	return e.flowFor(workflow)
}

func (e *WorkflowEngine) flowFor(workflow string) map[string][]string {
	workflow = strings.ToLower(strings.TrimSpace(workflow))
	if flow, ok := e.flows[workflow]; ok {
		return flow
	}
	return map[string][]string{}
}

func (e *WorkflowEngine) TerminalStates(workflow string) []string {
	flow := e.flowFor(workflow)
	terminal := make([]string, 0)
	for status, next := range flow {
		if len(next) == 0 {
			terminal = append(terminal, status)
		}
	}
	sort.Strings(terminal)
	return terminal
}

func (e *WorkflowEngine) IsTerminal(workflow, status string) bool {
	status = NormalizeStatus(status)
	for _, candidate := range e.TerminalStates(workflow) {
		if candidate == status {
			return true
		}
	}
	return false
}

func (e *WorkflowEngine) ApplyPlanStatus(plan *domain.CollectionPlan, next string) error {
	if plan == nil {
		return fmt.Errorf("plan required")
	}
	if err := e.ValidateTransition(WorkflowPlan, plan.Status, next); err != nil {
		return err
	}
	plan.Status = NormalizeStatus(next)
	return nil
}

func (e *WorkflowEngine) ApplyRouteStatus(route *domain.Route, next string) error {
	if route == nil {
		return fmt.Errorf("route required")
	}
	if err := e.ValidateTransition(WorkflowRoute, route.Status, next); err != nil {
		return err
	}
	route.Status = NormalizeStatus(next)
	return nil
}

func (e *WorkflowEngine) ApplyTaskStatus(task *domain.Task, next string, moment time.Time) error {
	if task == nil {
		return fmt.Errorf("task required")
	}
	if err := e.ValidateTransition(WorkflowTask, task.Status, next); err != nil {
		return err
	}
	task.Status = NormalizeStatus(next)
	switch task.Status {
	case domain.TaskClaimed:
		task.ClaimedAt = &moment
	case domain.TaskCompleted, domain.TaskSkipped:
		task.CompletedAt = &moment
	}
	return nil
}

func (e *WorkflowEngine) ValidateTaskStopsForCompletion(stops []domain.TaskStop) error {
	if len(stops) == 0 {
		return fmt.Errorf("task stops required")
	}
	for _, stop := range stops {
		switch NormalizeStatus(stop.Status) {
		case "completed":
			continue
		case "skipped":
			if strings.TrimSpace(stop.SkipReason) == "" {
				return fmt.Errorf("skipped task stop requires reason")
			}
		case "pending":
			return fmt.Errorf("task stop in pending status cannot complete task")
		default:
			return fmt.Errorf("task stop in %s status cannot complete task", stop.Status)
		}
	}
	return nil
}

func (e *WorkflowEngine) ApplyInvoiceStatus(invoice *domain.Invoice, next string) error {
	if invoice == nil {
		return fmt.Errorf("invoice required")
	}
	if err := e.ValidateTransition(WorkflowInvoice, invoice.Status, next); err != nil {
		return err
	}
	invoice.Status = NormalizeStatus(next)
	return nil
}

func (e *WorkflowEngine) ApplyWeighingStatus(record *domain.WeighingRecord, next string, confirmer *int64, confirmedAt *time.Time) error {
	if record == nil {
		return fmt.Errorf("weighing record required")
	}
	if err := e.ValidateTransition(WorkflowWeighing, record.Status, next); err != nil {
		return err
	}
	record.Status = NormalizeStatus(next)
	record.ConfirmedBy = confirmer
	record.ConfirmedAt = confirmedAt
	return nil
}

func (e *WorkflowEngine) ApplyUserStatus(user *domain.User, next string) error {
	if user == nil {
		return fmt.Errorf("user required")
	}
	if err := e.ValidateTransition(WorkflowUser, user.Status, next); err != nil {
		return err
	}
	user.Status = NormalizeStatus(next)
	return nil
}

func (e *WorkflowEngine) ApplyVehicleStatus(vehicle *domain.Vehicle, next string) error {
	if vehicle == nil {
		return fmt.Errorf("vehicle required")
	}
	if err := e.ValidateTransition(WorkflowVehicle, vehicle.Status, next); err != nil {
		return err
	}
	vehicle.Status = NormalizeStatus(next)
	return nil
}

func (e *WorkflowEngine) History(entity string, entityID int64, from, to, reason string, operatorID int64) domain.StatusHistory {
	return domain.StatusHistory{
		EntityID:   entityID,
		FromStatus: NormalizeStatus(from),
		ToStatus:   NormalizeStatus(to),
		OperatorID: operatorID,
		Reason:     strings.TrimSpace(reason),
	}
}

func (e *WorkflowEngine) BatchHistory(entity string, entityID int64, operatorID int64, transitions []StatusTransition) []domain.StatusHistory {
	histories := make([]domain.StatusHistory, 0, len(transitions))
	for _, transition := range transitions {
		histories = append(histories, e.History(entity, entityID, transition.From, transition.To, transition.Reason, operatorID))
	}
	return histories
}

type StatusTransition struct {
	From   string
	To     string
	Reason string
}

func NormalizeStatusList(values []string) []string {
	unique := make(map[string]struct{}, len(values))
	for _, value := range values {
		normalized := NormalizeStatus(value)
		if normalized == "" {
			continue
		}
		unique[normalized] = struct{}{}
	}
	ordered := make([]string, 0, len(unique))
	for value := range unique {
		ordered = append(ordered, value)
	}
	sort.Strings(ordered)
	return ordered
}

func MergeStatusHistory(existing []domain.StatusHistory, incoming []domain.StatusHistory) []domain.StatusHistory {
	merged := append([]domain.StatusHistory(nil), existing...)
	merged = append(merged, incoming...)
	sort.SliceStable(merged, func(i, j int) bool {
		if merged[i].EntityID == merged[j].EntityID {
			if merged[i].CreatedAt.Equal(merged[j].CreatedAt) {
				return merged[i].ToStatus < merged[j].ToStatus
			}
			return merged[i].CreatedAt.Before(merged[j].CreatedAt)
		}
		return merged[i].EntityID < merged[j].EntityID
	})
	return merged
}


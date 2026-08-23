package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
	"go-waste-routes/internal/transport/contextx"
	"go-waste-routes/internal/transport/http/response"
)

type OpsHandler struct {
	Dashboard     *service.DashboardService
	Exporter      *service.ExportService
	Reviews       *service.ReviewService
	Abnormalities *service.ResourceService[domain.WeighingAbnormality]
	Invoices      *service.ResourceService[domain.Invoice]
	Payments      *service.ResourceService[domain.PaymentRecord]
	Tasks         *service.ResourceService[domain.Task]
	TaskStops     *service.ResourceService[domain.TaskStop]
	TaskActions   *service.TaskService
}

type taskCompletionPayload struct {
	MileageKm  float64 `json:"mileage_km"`
	FuelLiters float64 `json:"fuel_liters"`
}

type paymentPayload struct {
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"`
	VoucherNumber string  `json:"voucher_number"`
	Notes         string  `json:"notes"`
	OperatorID    int64   `json:"operator_id"`
}

type reviewPayload struct {
	Approved bool   `json:"approved"`
	Reason   string `json:"reason"`
	Reviewer int64  `json:"reviewer_id"`
}

func (h *OpsHandler) DashboardView(w http.ResponseWriter, r *http.Request) {
	snapshot, err := h.Dashboard.Snapshot(r.Context(), time.Now())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, response.CodeSystemError, err.Error(), contextx.RequestID(r.Context()))
		return
	}
	response.OK(w, snapshot, contextx.RequestID(r.Context()))
}

func (h *OpsHandler) CompleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil { response.Error(w, http.StatusBadRequest, response.CodeValidation, "invalid id", contextx.RequestID(r.Context())); return }
	var payload taskCompletionPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil { response.Error(w, http.StatusBadRequest, response.CodeValidation, "invalid json", contextx.RequestID(r.Context())); return }
	task, err := h.Tasks.Get(r.Context(), id)
	if err != nil { response.Error(w, http.StatusNotFound, response.CodeNotFound, "not found", contextx.RequestID(r.Context())); return }
	allStops, _, err := h.TaskStops.List(r.Context(), 1, 10000)
	if err != nil { response.Error(w, http.StatusInternalServerError, response.CodeSystemError, err.Error(), contextx.RequestID(r.Context())); return }
	stops := make([]domain.TaskStop, 0)
	for _, stop := range allStops { if stop.TaskID == id { stops = append(stops, stop) } }
	if err := h.TaskActions.CompleteWithStops(&task, stops, time.Now(), payload.MileageKm, payload.FuelLiters); err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeBusiness, err.Error(), contextx.RequestID(r.Context()))
	}
	updated, err := h.Tasks.Update(r.Context(), id, task)
	if err != nil { response.Error(w, http.StatusInternalServerError, response.CodeSystemError, err.Error(), contextx.RequestID(r.Context())); return }
	response.OK(w, updated, contextx.RequestID(r.Context()))
}

func (h *OpsHandler) Export(w http.ResponseWriter, r *http.Request) {
	kind := r.URL.Query().Get("kind")
	format := r.URL.Query().Get("format")
	artifact, err := h.Exporter.Export(r.Context(), kind, format)
	if err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeBusiness, err.Error(), contextx.RequestID(r.Context()))
		return
	}
	w.Header().Set("Content-Type", artifact.ContentType)
	w.Header().Set("Content-Disposition", "attachment; filename="+artifact.FileName)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(artifact.Body)
}

func (h *OpsHandler) ReviewAbnormality(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeValidation, "invalid id", contextx.RequestID(r.Context()))
		return
	}
	var payload reviewPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeValidation, "invalid json", contextx.RequestID(r.Context()))
		return
	}
	abnormality, err := h.Abnormalities.Get(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, response.CodeNotFound, "not found", contextx.RequestID(r.Context()))
		return
	}
	if err := h.Reviews.ReviewAbnormality(&abnormality, service.ReviewDecision{
		Approved:   payload.Approved,
		Reason:     payload.Reason,
		ReviewerID: payload.Reviewer,
		ReviewedAt: time.Now(),
	}); err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeBusiness, err.Error(), contextx.RequestID(r.Context()))
		return
	}
	if _, err := h.Abnormalities.Update(r.Context(), id, abnormality); err != nil {
		response.Error(w, http.StatusInternalServerError, response.CodeSystemError, err.Error(), contextx.RequestID(r.Context()))
		return
	}
	response.OK(w, abnormality, contextx.RequestID(r.Context()))
}

func (h *OpsHandler) ConfirmInvoice(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeValidation, "invalid id", contextx.RequestID(r.Context()))
		return
	}
	invoice, err := h.Invoices.Get(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, response.CodeNotFound, "not found", contextx.RequestID(r.Context()))
		return
	}
	if err := h.Reviews.ConfirmInvoice(&invoice); err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeBusiness, err.Error(), contextx.RequestID(r.Context()))
		return
	}
	if _, err := h.Invoices.Update(r.Context(), id, invoice); err != nil {
		response.Error(w, http.StatusInternalServerError, response.CodeSystemError, err.Error(), contextx.RequestID(r.Context()))
		return
	}
	response.OK(w, invoice, contextx.RequestID(r.Context()))
}

func (h *OpsHandler) RecordPayment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeValidation, "invalid id", contextx.RequestID(r.Context()))
		return
	}
	var payload paymentPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeValidation, "invalid json", contextx.RequestID(r.Context()))
		return
	}
	invoice, err := h.Invoices.Get(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, response.CodeNotFound, "not found", contextx.RequestID(r.Context()))
		return
	}
	record := domain.PaymentRecord{}
	if err := h.Reviews.RecordPayment(&invoice, &record, payload.Amount, payload.OperatorID, payload.PaymentMethod, payload.VoucherNumber, time.Now()); err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeBusiness, err.Error(), contextx.RequestID(r.Context()))
		return
	}
	created, err := h.Payments.Create(r.Context(), record)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, response.CodeSystemError, err.Error(), contextx.RequestID(r.Context()))
		return
	}
	if _, err := h.Invoices.Update(r.Context(), id, invoice); err != nil {
		response.Error(w, http.StatusInternalServerError, response.CodeSystemError, err.Error(), contextx.RequestID(r.Context()))
		return
	}
	response.Created(w, map[string]any{"invoice": invoice, "payment": created}, contextx.RequestID(r.Context()))
}

func (h *OpsHandler) ExportKinds(w http.ResponseWriter, r *http.Request) {
	kinds := []string{
		"customers", "routes", "tasks", "invoices", "weighings", "abnormalities",
		"payments", "audit-logs", "plan-rules", "route-stops", "task-stops", "invoice-items",
		"weighbridges", "system-configs", "holiday-configs",
	}
	response.OK(w, map[string]any{"kinds": kinds, "formats": []string{"json", "csv"}}, contextx.RequestID(r.Context()))
}

func (h *OpsHandler) PreviewExportName(kind, format string) string {
	kind = strings.TrimSpace(kind)
	format = strings.TrimSpace(format)
	if format == "" {
		format = "json"
	}
	return kind + "." + format
}


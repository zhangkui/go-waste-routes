package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"go-waste-routes/internal/domain"
)

type ExportArtifact struct {
	FileName    string
	ContentType string
	Body        []byte
}

type ExportService struct {
	Customers     *ResourceService[domain.Customer]
	Routes        *ResourceService[domain.Route]
	Tasks         *ResourceService[domain.Task]
	Invoices      *ResourceService[domain.Invoice]
	Weighings     *ResourceService[domain.WeighingRecord]
	Abnormalities *ResourceService[domain.WeighingAbnormality]
	Payments      *ResourceService[domain.PaymentRecord]
	AuditLogs     *ResourceService[domain.AuditLog]
	PlanRules     *ResourceService[domain.PlanRule]
	RouteStops    *ResourceService[domain.RouteStop]
	TaskStops     *ResourceService[domain.TaskStop]
	InvoiceItems  *ResourceService[domain.InvoiceItem]
	Weighbridges  *ResourceService[domain.Weighbridge]
	SystemConfigs *ResourceService[domain.SystemConfig]
	Holidays      *ResourceService[domain.HolidayConfig]
}

func NewExportService(deps ExportService) *ExportService {
	return &deps
}

func (s *ExportService) Export(ctx context.Context, kind string, format string) (ExportArtifact, error) {
	kind = strings.ToLower(strings.TrimSpace(kind))
	format = strings.ToLower(strings.TrimSpace(format))
	switch kind {
	case "customers":
		return exportResource(ctx, "customers", format, "customers", s.Customers)
	case "routes":
		return exportResource(ctx, "routes", format, "routes", s.Routes)
	case "tasks":
		return exportResource(ctx, "tasks", format, "tasks", s.Tasks)
	case "invoices":
		return exportResource(ctx, "invoices", format, "invoices", s.Invoices)
	case "weighings":
		return exportResource(ctx, "weighings", format, "weighings", s.Weighings)
	case "abnormalities":
		return exportResource(ctx, "abnormalities", format, "abnormalities", s.Abnormalities)
	case "payments":
		return exportResource(ctx, "payments", format, "payments", s.Payments)
	case "audit-logs":
		return exportResource(ctx, "audit-logs", format, "audit_logs", s.AuditLogs)
	case "plan-rules":
		return exportResource(ctx, "plan-rules", format, "plan_rules", s.PlanRules)
	case "route-stops":
		return exportResource(ctx, "route-stops", format, "route_stops", s.RouteStops)
	case "task-stops":
		return exportResource(ctx, "task-stops", format, "task_stops", s.TaskStops)
	case "invoice-items":
		return exportResource(ctx, "invoice-items", format, "invoice_items", s.InvoiceItems)
	case "weighbridges":
		return exportResource(ctx, "weighbridges", format, "weighbridges", s.Weighbridges)
	case "system-configs":
		return exportResource(ctx, "system-configs", format, "system_configs", s.SystemConfigs)
	case "holiday-configs":
		return exportResource(ctx, "holiday-configs", format, "holiday_configs", s.Holidays)
	default:
		return ExportArtifact{}, fmt.Errorf("unsupported export kind %s", kind)
	}
}

func exportResource[T any](ctx context.Context, kind, format, prefix string, resource *ResourceService[T]) (ExportArtifact, error) {
	if resource == nil {
		return ExportArtifact{}, fmt.Errorf("export resource %s unavailable", kind)
	}
	items, err := LoadAll(ctx, resource)
	if err != nil {
		return ExportArtifact{}, err
	}
	switch format {
	case "csv":
		body, err := encodeCSV(items)
		if err != nil {
			return ExportArtifact{}, err
		}
		return ExportArtifact{
			FileName:    fmt.Sprintf("%s-%s.csv", prefix, time.Now().Format("20060102-150405")),
			ContentType: "text/csv; charset=utf-8",
			Body:        body,
		}, nil
	default:
		body, err := json.MarshalIndent(items, "", "  ")
		if err != nil {
			return ExportArtifact{}, err
		}
		return ExportArtifact{
			FileName:    fmt.Sprintf("%s-%s.json", prefix, time.Now().Format("20060102-150405")),
			ContentType: "application/json; charset=utf-8",
			Body:        body,
		}, nil
	}
}

func encodeCSV[T any](items []T) ([]byte, error) {
	buffer := &bytes.Buffer{}
	writer := csv.NewWriter(buffer)
	rows := make([]map[string]string, 0, len(items))
	columns := make(map[string]struct{})
	for _, item := range items {
		row, err := structToStringMap(item)
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
		for key := range row {
			columns[key] = struct{}{}
		}
	}
	headers := make([]string, 0, len(columns))
	for key := range columns {
		headers = append(headers, key)
	}
	sort.Strings(headers)
	if err := writer.Write(headers); err != nil {
		return nil, err
	}
	for _, row := range rows {
		values := make([]string, len(headers))
		for index, header := range headers {
			values[index] = row[header]
		}
		if err := writer.Write(values); err != nil {
			return nil, err
		}
	}
	writer.Flush()
	return buffer.Bytes(), writer.Error()
}

func structToStringMap[T any](item T) (map[string]string, error) {
	raw, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	var data map[string]any
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, err
	}
	result := make(map[string]string, len(data))
	for key, value := range data {
		result[key] = stringify(value)
	}
	return result, nil
}

func stringify(value any) string {
	if value == nil {
		return ""
	}
	switch typed := value.(type) {
	case string:
		return typed
	case float64:
		return fmt.Sprintf("%.2f", typed)
	case bool:
		if typed {
			return "true"
		}
		return "false"
	case map[string]any, []any:
		raw, _ := json.Marshal(typed)
		return string(raw)
	default:
		if reflect.TypeOf(value).Kind() == reflect.Slice {
			raw, _ := json.Marshal(value)
			return string(raw)
		}
		return fmt.Sprintf("%v", value)
	}
}

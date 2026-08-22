package service

import (
	"fmt"
	"time"

	"go-waste-routes/internal/domain"
)

type ReviewService struct {
	workflow *WorkflowEngine
	billing  *BillingEngine
	weighing *WeighingEngine
}

type ReviewDecision struct {
	Approved   bool
	Reason     string
	ReviewerID int64
	ReviewedAt time.Time
}

func NewReviewService(workflow *WorkflowEngine, billing *BillingEngine, weighing *WeighingEngine) *ReviewService {
	return &ReviewService{workflow: workflow, billing: billing, weighing: weighing}
}

func (s *ReviewService) ReviewAbnormality(abnormality *domain.WeighingAbnormality, decision ReviewDecision) error {
	if abnormality == nil {
		return fmt.Errorf("abnormality required")
	}
	result := "rejected"
	nextStatus := "rejected"
	if decision.Approved {
		result = "confirmed"
		nextStatus = "confirmed"
	}
	abnormality.ReviewerID = &decision.ReviewerID
	abnormality.ReviewedAt = &decision.ReviewedAt
	abnormality.ReviewResult = result
	abnormality.RejectReason = decision.Reason
	abnormality.Status = nextStatus
	return nil
}

func (s *ReviewService) ConfirmWeighing(record *domain.WeighingRecord, confirmerID int64, confirmedAt time.Time) error {
	if record == nil {
		return fmt.Errorf("weighing record required")
	}
	return s.workflow.ApplyWeighingStatus(record, "confirmed", &confirmerID, &confirmedAt)
}

func (s *ReviewService) RejectWeighing(record *domain.WeighingRecord, reason string) error {
	if record == nil {
		return fmt.Errorf("weighing record required")
	}
	record.Status = "rejected"
	record.ManualReason = reason
	return nil
}

func (s *ReviewService) ConfirmInvoice(invoice *domain.Invoice) error {
	if invoice == nil {
		return fmt.Errorf("invoice required")
	}
	return s.workflow.ApplyInvoiceStatus(invoice, domain.InvoiceConfirmed)
}

func (s *ReviewService) WriteOffInvoice(invoice *domain.Invoice) error {
	if invoice == nil {
		return fmt.Errorf("invoice required")
	}
	s.billing.WriteOff(invoice)
	return nil
}

func (s *ReviewService) RecordPayment(invoice *domain.Invoice, record *domain.PaymentRecord, amount float64, operatorID int64, method, voucher string, when time.Time) error {
	if invoice == nil {
		return fmt.Errorf("invoice required")
	}
	if record == nil {
		return fmt.Errorf("payment record required")
	}
	record.InvoiceID = invoice.ID
	record.Amount = amount
	record.OperatorID = operatorID
	record.PaymentMethod = method
	record.VoucherNumber = voucher
	record.PaymentTime = when
	s.billing.ApplyPayment(invoice, amount)
	return nil
}


package service

import (
	"encoding/json"
	"fmt"
	"time"

	"go-waste-routes/internal/domain"
)

type AuditService struct{}

func NewAuditService() *AuditService { return &AuditService{} }

func (s *AuditService) Record(username, action, resourceType, resourceID, requestID string, before, after any) domain.AuditLog {
	return domain.AuditLog{
		Username:     username,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		BeforeValue:  marshalJSON(before),
		AfterValue:   marshalJSON(after),
		RequestID:    requestID,
		CreatedAt:    time.Now().UTC(),
	}
}

func (s *AuditService) StatusChange(username, resourceType, resourceID, from, to string) domain.AuditLog {
	return s.Record(username, "status_change", resourceType, resourceID, "", map[string]string{"from": from}, map[string]string{"to": to})
}

func marshalJSON(value any) string {
	if value == nil {
		return ""
	}
	content, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprintf("%v", value)
	}
	return string(content)
}

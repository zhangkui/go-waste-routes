package service

import (
	"fmt"
	"time"
)

type Notification struct {
	Target   string    `json:"target"`
	Subject  string    `json:"subject"`
	Body     string    `json:"body"`
	SendAt   time.Time `json:"send_at"`
	Channel  string    `json:"channel"`
	Disabled bool      `json:"disabled"`
}

type NotificationService struct{}

func NewNotificationService() *NotificationService { return &NotificationService{} }

func (s *NotificationService) BuildOverdue(target string, invoiceNumber string, dueDate time.Time) Notification {
	return Notification{
		Target:  target,
		Subject: "账单逾期提醒",
		Body:    fmt.Sprintf("账单 %s 已逾期，截止日期 %s", invoiceNumber, dueDate.Format("2006-01-02")),
		SendAt:  time.Now().UTC(),
		Channel: "in-app",
	}
}

func (s *NotificationService) BuildTaskSummary(target string, taskNumber string, count int) Notification {
	return Notification{
		Target:  target,
		Subject: "任务汇总",
		Body:    fmt.Sprintf("任务 %s 共 %d 个站点", taskNumber, count),
		SendAt:  time.Now().UTC(),
		Channel: "in-app",
	}
}


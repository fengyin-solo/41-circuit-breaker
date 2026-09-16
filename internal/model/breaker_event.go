package model

import (
	"strings"
	"time"
)

const (
	EventTypeOpened   = "opened"
	EventTypeClosed   = "closed"
	EventTypeHalfOpen = "half_open"
	EventTypeRejected = "rejected"
)

type BreakerEvent struct {
	ID         string    `json:"id"`
	BreakerID  string    `json:"breaker_id"`
	ServiceID  string    `json:"service_id"`
	EventType  string    `json:"event_type"`
	Reason     string    `json:"reason"`
	OccurredAt time.Time `json:"occurred_at"`
}

func (e *BreakerEvent) Validate() error {
	if e.BreakerID == "" {
		return NewValidationError("breaker_id", "熔断器 ID 不能为空")
	}
	if e.ServiceID == "" {
		return NewValidationError("service_id", "服务 ID 不能为空")
	}
	if e.EventType == "" {
		return NewValidationError("event_type", "事件类型不能为空")
	}
	if e.EventType != EventTypeOpened && e.EventType != EventTypeClosed && e.EventType != EventTypeHalfOpen && e.EventType != EventTypeRejected {
		return NewValidationError("event_type", "事件类型不合法")
	}
	if e.OccurredAt.IsZero() {
		e.OccurredAt = time.Now()
	}
	return nil
}

type BreakerEventFilter struct {
	BreakerID string
	ServiceID string
	EventType string
	Keyword   string
	StartTime *time.Time
	EndTime   *time.Time
}

func (f BreakerEventFilter) Match(e *BreakerEvent) bool {
	if f.BreakerID != "" && e.BreakerID != f.BreakerID {
		return false
	}
	if f.ServiceID != "" && e.ServiceID != f.ServiceID {
		return false
	}
	if f.EventType != "" && e.EventType != f.EventType {
		return false
	}
	if f.StartTime != nil && e.OccurredAt.Before(*f.StartTime) {
		return false
	}
	if f.EndTime != nil && e.OccurredAt.After(*f.EndTime) {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(e.Reason), k) {
			return false
		}
	}
	return true
}

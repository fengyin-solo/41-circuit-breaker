package model

import (
	"strings"
	"time"
)

const (
	HealthCheckStatusHealthy   = "healthy"
	HealthCheckStatusUnhealthy = "unhealthy"
)

type HealthCheck struct {
	ID                  string    `json:"id"`
	ServiceID           string    `json:"service_id"`
	IntervalSeconds     int       `json:"interval_seconds"`
	LastCheckedAt       *time.Time `json:"last_checked_at,omitempty"`
	LastStatus          string    `json:"last_status"`
	ConsecutiveFailures int       `json:"consecutive_failures"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (h *HealthCheck) Validate() error {
	if h.ServiceID == "" {
		return NewValidationError("service_id", "服务 ID 不能为空")
	}
	if h.IntervalSeconds <= 0 {
		h.IntervalSeconds = 10
	}
	if h.LastStatus == "" {
		h.LastStatus = HealthCheckStatusHealthy
	}
	if h.LastStatus != HealthCheckStatusHealthy && h.LastStatus != HealthCheckStatusUnhealthy {
		return NewValidationError("last_status", "状态必须为 healthy 或 unhealthy")
	}
	if h.ConsecutiveFailures < 0 {
		h.ConsecutiveFailures = 0
	}
	return nil
}

type HealthCheckFilter struct {
	ServiceID string
	LastStatus string
	Keyword   string
}

func (f HealthCheckFilter) Match(h *HealthCheck) bool {
	if f.ServiceID != "" && h.ServiceID != f.ServiceID {
		return false
	}
	if f.LastStatus != "" && h.LastStatus != f.LastStatus {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(h.ServiceID), k) {
			return false
		}
	}
	return true
}

package model

import (
	"strings"
	"time"
)

const (
	BreakerStateClosed    = "closed"
	BreakerStateOpen      = "open"
	BreakerStateHalfOpen  = "half_open"
)

type CircuitBreaker struct {
	ID            string    `json:"id"`
	ServiceID     string    `json:"service_id"`
	RuleID        string    `json:"rule_id"`
	State         string    `json:"state"`
	FailureCount  int       `json:"failure_count"`
	SuccessCount  int       `json:"success_count"`
	TotalCalls    int       `json:"total_calls"`
	LastOpenedAt  *time.Time `json:"last_opened_at,omitempty"`
	LastClosedAt  *time.Time `json:"last_closed_at,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (b *CircuitBreaker) Validate() error {
	if b.ServiceID == "" {
		return NewValidationError("service_id", "服务 ID 不能为空")
	}
	if b.RuleID == "" {
		return NewValidationError("rule_id", "规则 ID 不能为空")
	}
	if b.State == "" {
		b.State = BreakerStateClosed
	}
	if b.State != BreakerStateClosed && b.State != BreakerStateOpen && b.State != BreakerStateHalfOpen {
		return NewValidationError("state", "状态必须为 closed、open 或 half_open")
	}
	if b.FailureCount < 0 {
		b.FailureCount = 0
	}
	if b.SuccessCount < 0 {
		b.SuccessCount = 0
	}
	if b.TotalCalls < 0 {
		b.TotalCalls = 0
	}
	return nil
}

var breakerTransitions = map[string]map[string]bool{
	BreakerStateClosed:   {BreakerStateOpen: true},
	BreakerStateOpen:     {BreakerStateHalfOpen: true},
	BreakerStateHalfOpen: {BreakerStateClosed: true, BreakerStateOpen: true},
}

func CanTransitionBreaker(from, to string) bool {
	if m, ok := breakerTransitions[from]; ok {
		return m[to]
	}
	return false
}

type CircuitBreakerFilter struct {
	ServiceID string
	RuleID    string
	State     string
	Keyword   string
}

func (f CircuitBreakerFilter) Match(b *CircuitBreaker) bool {
	if f.ServiceID != "" && b.ServiceID != f.ServiceID {
		return false
	}
	if f.RuleID != "" && b.RuleID != f.RuleID {
		return false
	}
	if f.State != "" && b.State != f.State {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(b.ServiceID), k) {
			return false
		}
	}
	return true
}

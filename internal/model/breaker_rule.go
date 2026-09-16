package model

import (
	"strings"
	"time"
)

type BreakerRule struct {
	ID                      string    `json:"id"`
	Name                    string    `json:"name"`
	ServiceID               string    `json:"service_id"`
	FailureRatioThreshold   float64   `json:"failure_ratio_threshold"`
	SlowCallRatioThreshold  float64   `json:"slow_call_ratio_threshold"`
	SlowCallMs              int       `json:"slow_call_ms"`
	WindowSeconds           int       `json:"window_seconds"`
	MinRequestCount         int       `json:"min_request_count"`
	MaxHalfOpenRequests     int       `json:"max_half_open_requests"`
	Enabled                 bool      `json:"enabled"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

func (r *BreakerRule) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	if r.Name == "" {
		return NewValidationError("name", "规则名称不能为空")
	}
	if r.ServiceID == "" {
		return NewValidationError("service_id", "关联服务 ID 不能为空")
	}
	if r.FailureRatioThreshold < 0 || r.FailureRatioThreshold > 1 {
		return NewValidationError("failure_ratio_threshold", "失败率阈值必须在 0~1 之间")
	}
	if r.SlowCallRatioThreshold < 0 || r.SlowCallRatioThreshold > 1 {
		return NewValidationError("slow_call_ratio_threshold", "慢调用率阈值必须在 0~1 之间")
	}
	if r.SlowCallMs <= 0 {
		r.SlowCallMs = 1000
	}
	if r.WindowSeconds <= 0 {
		r.WindowSeconds = 60
	}
	if r.MinRequestCount <= 0 {
		r.MinRequestCount = 10
	}
	if r.MaxHalfOpenRequests <= 0 {
		r.MaxHalfOpenRequests = 5
	}
	return nil
}

type BreakerRuleFilter struct {
	Name      string
	ServiceID string
	Enabled   *bool
	Keyword   string
}

func (f BreakerRuleFilter) Match(r *BreakerRule) bool {
	if f.Name != "" && r.Name != f.Name {
		return false
	}
	if f.ServiceID != "" && r.ServiceID != f.ServiceID {
		return false
	}
	if f.Enabled != nil && r.Enabled != *f.Enabled {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(r.Name), k) {
			return false
		}
	}
	return true
}

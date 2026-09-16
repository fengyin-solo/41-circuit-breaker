package model

import (
	"strings"
	"time"
)

type RecoveryPolicy struct {
	ID                    string    `json:"id"`
	Name                  string    `json:"name"`
	ServiceID             string    `json:"service_id"`
	HalfOpenProbeRatio    float64   `json:"half_open_probe_ratio"`
	RecoveryWindowSeconds int       `json:"recovery_window_seconds"`
	MaxRetry              int       `json:"max_retry"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

func (p *RecoveryPolicy) Validate() error {
	p.Name = strings.TrimSpace(p.Name)
	if p.Name == "" {
		return NewValidationError("name", "策略名称不能为空")
	}
	if p.ServiceID == "" {
		return NewValidationError("service_id", "服务 ID 不能为空")
	}
	if p.HalfOpenProbeRatio <= 0 || p.HalfOpenProbeRatio > 1 {
		return NewValidationError("half_open_probe_ratio", "探测比例必须在 0~1 之间")
	}
	if p.RecoveryWindowSeconds <= 0 {
		p.RecoveryWindowSeconds = 30
	}
	if p.MaxRetry < 0 {
		p.MaxRetry = 3
	}
	if p.MaxRetry > 10 {
		return NewValidationError("max_retry", "最大重试次数不能超过 10")
	}
	return nil
}

type RecoveryPolicyFilter struct {
	Name      string
	ServiceID string
	Keyword   string
}

func (f RecoveryPolicyFilter) Match(p *RecoveryPolicy) bool {
	if f.Name != "" && p.Name != f.Name {
		return false
	}
	if f.ServiceID != "" && p.ServiceID != f.ServiceID {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(p.Name), k) {
			return false
		}
	}
	return true
}

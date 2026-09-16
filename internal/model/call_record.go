package model

import (
	"strings"
	"time"
)

const (
	CallOutcomeSuccess = "success"
	CallOutcomeFailure = "failure"
	CallOutcomeTimeout = "timeout"
)

type CallRecord struct {
	ID        string    `json:"id"`
	RequestID string    `json:"request_id"`
	ServiceID string    `json:"service_id"`
	Outcome   string    `json:"outcome"`
	LatencyMs int       `json:"latency_ms"`
	CalledAt  time.Time `json:"called_at"`
}

func (c *CallRecord) Validate() error {
	if c.RequestID == "" {
		return NewValidationError("request_id", "请求 ID 不能为空")
	}
	if c.ServiceID == "" {
		return NewValidationError("service_id", "服务 ID 不能为空")
	}
	if c.Outcome == "" {
		return NewValidationError("outcome", "调用结果不能为空")
	}
	if c.Outcome != CallOutcomeSuccess && c.Outcome != CallOutcomeFailure && c.Outcome != CallOutcomeTimeout {
		return NewValidationError("outcome", "结果必须为 success、failure 或 timeout")
	}
	if c.LatencyMs < 0 {
		c.LatencyMs = 0
	}
	if c.CalledAt.IsZero() {
		c.CalledAt = time.Now()
	}
	return nil
}

type CallRecordFilter struct {
	ServiceID string
	Outcome   string
	Keyword   string
	StartTime *time.Time
	EndTime   *time.Time
}

func (f CallRecordFilter) Match(c *CallRecord) bool {
	if f.ServiceID != "" && c.ServiceID != f.ServiceID {
		return false
	}
	if f.Outcome != "" && c.Outcome != f.Outcome {
		return false
	}
	if f.StartTime != nil && c.CalledAt.Before(*f.StartTime) {
		return false
	}
	if f.EndTime != nil && c.CalledAt.After(*f.EndTime) {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(c.RequestID), k) {
			return false
		}
	}
	return true
}

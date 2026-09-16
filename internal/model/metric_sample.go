package model

import (
	"strings"
	"time"
)

type MetricSample struct {
	ID            string    `json:"id"`
	ServiceID     string    `json:"service_id"`
	WindowStart   time.Time `json:"window_start"`
	WindowEnd     time.Time `json:"window_end"`
	TotalCalls    int       `json:"total_calls"`
	SuccessCalls  int       `json:"success_calls"`
	FailureCalls  int       `json:"failure_calls"`
	SlowCalls     int       `json:"slow_calls"`
	FailureRatio  float64   `json:"failure_ratio"`
	SlowRatio     float64   `json:"slow_ratio"`
	CreatedAt     time.Time `json:"created_at"`
}

func (m *MetricSample) Validate() error {
	if m.ServiceID == "" {
		return NewValidationError("service_id", "服务 ID 不能为空")
	}
	if m.WindowStart.IsZero() || m.WindowEnd.IsZero() {
		return NewValidationError("window", "时间窗口不能为空")
	}
	if m.WindowEnd.Before(m.WindowStart) {
		return NewValidationError("window", "窗口结束时间不能早于开始时间")
	}
	if m.TotalCalls < 0 {
		m.TotalCalls = 0
	}
	if m.SuccessCalls < 0 {
		m.SuccessCalls = 0
	}
	if m.FailureCalls < 0 {
		m.FailureCalls = 0
	}
	if m.SlowCalls < 0 {
		m.SlowCalls = 0
	}
	if m.FailureRatio < 0 || m.FailureRatio > 1 {
		return NewValidationError("failure_ratio", "失败率必须在 0~1 之间")
	}
	if m.SlowRatio < 0 || m.SlowRatio > 1 {
		return NewValidationError("slow_ratio", "慢调用率必须在 0~1 之间")
	}
	return nil
}

type MetricSampleFilter struct {
	ServiceID string
	Keyword   string
	StartTime *time.Time
	EndTime   *time.Time
}

func (f MetricSampleFilter) Match(m *MetricSample) bool {
	if f.ServiceID != "" && m.ServiceID != f.ServiceID {
		return false
	}
	if f.StartTime != nil && m.WindowStart.Before(*f.StartTime) {
		return false
	}
	if f.EndTime != nil && m.WindowEnd.After(*f.EndTime) {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(m.ServiceID), k) {
			return false
		}
	}
	return true
}

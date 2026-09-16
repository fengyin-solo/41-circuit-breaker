package model

import (
	"strings"
	"time"
)

type BreakerSnapshot struct {
	ID              string    `json:"id"`
	ServiceID       string    `json:"service_id"`
	State           string    `json:"state"`
	TotalCalls      int       `json:"total_calls"`
	SuccessCalls    int       `json:"success_calls"`
	FailureCalls    int       `json:"failure_calls"`
	SlowCalls       int       `json:"slow_calls"`
	FailureRatio    float64   `json:"failure_ratio"`
	SlowRatio       float64   `json:"slow_ratio"`
	SnapshotVersion int       `json:"snapshot_version"`
	CreatedAt       time.Time `json:"created_at"`
}

func (s *BreakerSnapshot) Validate() error {
	if s.ServiceID == "" {
		return NewValidationError("service_id", "服务 ID 不能为空")
	}
	if s.State == "" {
		return NewValidationError("state", "状态不能为空")
	}
	if s.State != BreakerStateClosed && s.State != BreakerStateOpen && s.State != BreakerStateHalfOpen {
		return NewValidationError("state", "状态不合法")
	}
	if s.TotalCalls < 0 {
		s.TotalCalls = 0
	}
	if s.SuccessCalls < 0 {
		s.SuccessCalls = 0
	}
	if s.FailureCalls < 0 {
		s.FailureCalls = 0
	}
	if s.SlowCalls < 0 {
		s.SlowCalls = 0
	}
	if s.FailureRatio < 0 || s.FailureRatio > 1 {
		return NewValidationError("failure_ratio", "失败率必须在 0~1 之间")
	}
	if s.SlowRatio < 0 || s.SlowRatio > 1 {
		return NewValidationError("slow_ratio", "慢调用率必须在 0~1 之间")
	}
	if s.SnapshotVersion <= 0 {
		s.SnapshotVersion = 1
	}
	return nil
}

type BreakerSnapshotFilter struct {
	ServiceID string
	State     string
	Keyword   string
}

func (f BreakerSnapshotFilter) Match(s *BreakerSnapshot) bool {
	if f.ServiceID != "" && s.ServiceID != f.ServiceID {
		return false
	}
	if f.State != "" && s.State != f.State {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(s.ServiceID), k) {
			return false
		}
	}
	return true
}

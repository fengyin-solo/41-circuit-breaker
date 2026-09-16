package model

import (
	"strings"
	"time"
)

const (
	AlertMetricStateChanged     = "state_changed"
	AlertMetricFailureRateHigh  = "failure_rate_high"
	AlertMetricConsecutiveFailure = "consecutive_failure"

	AlertSeverityInfo     = "info"
	AlertSeverityWarn     = "warn"
	AlertSeverityCritical = "critical"

	AlertChannelWebhook = "webhook"
	AlertChannelEmail   = "email"
	AlertChannelSMS     = "sms"
)

type AlertRule struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	ServiceID     string    `json:"service_id"`
	Metric        string    `json:"metric"`
	Threshold     float64   `json:"threshold"`
	Severity      string    `json:"severity"`
	NotifyChannel string    `json:"notify_channel"`
	Enabled       bool      `json:"enabled"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (a *AlertRule) Validate() error {
	a.Name = strings.TrimSpace(a.Name)
	if a.Name == "" {
		return NewValidationError("name", "规则名称不能为空")
	}
	if a.ServiceID == "" {
		return NewValidationError("service_id", "服务 ID 不能为空")
	}
	if a.Metric == "" {
		return NewValidationError("metric", "告警指标不能为空")
	}
	if a.Metric != AlertMetricStateChanged && a.Metric != AlertMetricFailureRateHigh && a.Metric != AlertMetricConsecutiveFailure {
		return NewValidationError("metric", "指标类型不合法")
	}
	if a.Threshold < 0 {
		return NewValidationError("threshold", "阈值不能为负数")
	}
	if a.Severity == "" {
		a.Severity = AlertSeverityInfo
	}
	if a.Severity != AlertSeverityInfo && a.Severity != AlertSeverityWarn && a.Severity != AlertSeverityCritical {
		return NewValidationError("severity", "严重级别不合法")
	}
	if a.NotifyChannel == "" {
		a.NotifyChannel = AlertChannelWebhook
	}
	if a.NotifyChannel != AlertChannelWebhook && a.NotifyChannel != AlertChannelEmail && a.NotifyChannel != AlertChannelSMS {
		return NewValidationError("notify_channel", "通知渠道不合法")
	}
	return nil
}

type AlertRuleFilter struct {
	Name      string
	ServiceID string
	Metric    string
	Severity  string
	Enabled   *bool
	Keyword   string
}

func (f AlertRuleFilter) Match(a *AlertRule) bool {
	if f.Name != "" && a.Name != f.Name {
		return false
	}
	if f.ServiceID != "" && a.ServiceID != f.ServiceID {
		return false
	}
	if f.Metric != "" && a.Metric != f.Metric {
		return false
	}
	if f.Severity != "" && a.Severity != f.Severity {
		return false
	}
	if f.Enabled != nil && a.Enabled != *f.Enabled {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(a.Name), k) {
			return false
		}
	}
	return true
}

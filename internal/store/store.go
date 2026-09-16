// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"circuitbreaker/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// DownstreamService
	CreateDownstreamService(s *model.DownstreamService) error
	GetDownstreamService(id string) (*model.DownstreamService, error)
	GetDownstreamServiceByName(name string) (*model.DownstreamService, error)
	ListDownstreamServices() []*model.DownstreamService
	UpdateDownstreamService(s *model.DownstreamService) error
	DeleteDownstreamService(id string) error

	// BreakerRule
	CreateBreakerRule(r *model.BreakerRule) error
	GetBreakerRule(id string) (*model.BreakerRule, error)
	ListBreakerRules() []*model.BreakerRule
	UpdateBreakerRule(r *model.BreakerRule) error
	DeleteBreakerRule(id string) error

	// CircuitBreaker
	CreateCircuitBreaker(b *model.CircuitBreaker) error
	GetCircuitBreaker(id string) (*model.CircuitBreaker, error)
	GetCircuitBreakerByServiceAndRule(serviceID, ruleID string) (*model.CircuitBreaker, error)
	ListCircuitBreakers() []*model.CircuitBreaker
	UpdateCircuitBreaker(b *model.CircuitBreaker) error
	DeleteCircuitBreaker(id string) error

	// CallRecord
	CreateCallRecord(c *model.CallRecord) error
	GetCallRecord(id string) (*model.CallRecord, error)
	ListCallRecords() []*model.CallRecord
	DeleteCallRecord(id string) error

	// HealthCheck
	CreateHealthCheck(h *model.HealthCheck) error
	GetHealthCheck(id string) (*model.HealthCheck, error)
	GetHealthCheckByServiceID(serviceID string) (*model.HealthCheck, error)
	ListHealthChecks() []*model.HealthCheck
	UpdateHealthCheck(h *model.HealthCheck) error
	DeleteHealthCheck(id string) error

	// AlertRule
	CreateAlertRule(a *model.AlertRule) error
	GetAlertRule(id string) (*model.AlertRule, error)
	ListAlertRules() []*model.AlertRule
	UpdateAlertRule(a *model.AlertRule) error
	DeleteAlertRule(id string) error

	// RecoveryPolicy
	CreateRecoveryPolicy(p *model.RecoveryPolicy) error
	GetRecoveryPolicy(id string) (*model.RecoveryPolicy, error)
	GetRecoveryPolicyByServiceID(serviceID string) (*model.RecoveryPolicy, error)
	ListRecoveryPolicies() []*model.RecoveryPolicy
	UpdateRecoveryPolicy(p *model.RecoveryPolicy) error
	DeleteRecoveryPolicy(id string) error

	// MetricSample
	CreateMetricSample(m *model.MetricSample) error
	GetMetricSample(id string) (*model.MetricSample, error)
	ListMetricSamples() []*model.MetricSample
	DeleteMetricSample(id string) error

	// BreakerEvent
	CreateBreakerEvent(e *model.BreakerEvent) error
	GetBreakerEvent(id string) (*model.BreakerEvent, error)
	ListBreakerEvents() []*model.BreakerEvent
	DeleteBreakerEvent(id string) error

	// BreakerSnapshot
	CreateBreakerSnapshot(s *model.BreakerSnapshot) error
	GetBreakerSnapshot(id string) (*model.BreakerSnapshot, error)
	ListBreakerSnapshots() []*model.BreakerSnapshot
	UpdateBreakerSnapshot(s *model.BreakerSnapshot) error
	DeleteBreakerSnapshot(id string) error
}

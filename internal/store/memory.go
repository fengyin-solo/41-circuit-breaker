package store

import (
	"sync"

	"circuitbreaker/internal/model"
)

type MemoryStore struct {
	mu                sync.RWMutex
	downstreamServices map[string]*model.DownstreamService
	breakerRules       map[string]*model.BreakerRule
	circuitBreakers    map[string]*model.CircuitBreaker
	callRecords        map[string]*model.CallRecord
	healthChecks       map[string]*model.HealthCheck
	alertRules         map[string]*model.AlertRule
	recoveryPolicies   map[string]*model.RecoveryPolicy
	metricSamples      map[string]*model.MetricSample
	breakerEvents      map[string]*model.BreakerEvent
	breakerSnapshots   map[string]*model.BreakerSnapshot
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		downstreamServices: make(map[string]*model.DownstreamService),
		breakerRules:       make(map[string]*model.BreakerRule),
		circuitBreakers:    make(map[string]*model.CircuitBreaker),
		callRecords:        make(map[string]*model.CallRecord),
		healthChecks:       make(map[string]*model.HealthCheck),
		alertRules:         make(map[string]*model.AlertRule),
		recoveryPolicies:   make(map[string]*model.RecoveryPolicy),
		metricSamples:      make(map[string]*model.MetricSample),
		breakerEvents:      make(map[string]*model.BreakerEvent),
		breakerSnapshots:   make(map[string]*model.BreakerSnapshot),
	}
}

var _ Store = (*MemoryStore)(nil)

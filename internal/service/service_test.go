package service

import (
	"testing"
	"time"

	"circuitbreaker/internal/config"
	"circuitbreaker/internal/model"
	"circuitbreaker/internal/store"
	"circuitbreaker/pkg/logger"
)

func newTestService() *Service {
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	return New(store.NewMemoryStore(), log, cfg)
}

func TestService_CreateDownstreamService(t *testing.T) {
	svc := newTestService()
	s, err := svc.CreateDownstreamService(model.DownstreamService{Name: "svc1", Address: "127.0.0.1:8080"})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if s.ID == "" {
		t.Fatalf("id not generated")
	}
	if _, err := svc.CreateDownstreamService(model.DownstreamService{Name: "svc1", Address: "127.0.0.1:8081"}); err == nil {
		t.Fatalf("expected conflict")
	}
}

func TestService_ListDownstreamServices(t *testing.T) {
	svc := newTestService()
	_, _ = svc.CreateDownstreamService(model.DownstreamService{Name: "svc1", Address: "127.0.0.1:8080"})
	_, _ = svc.CreateDownstreamService(model.DownstreamService{Name: "svc2", Address: "127.0.0.1:8081"})
	items, total, err := svc.ListDownstreamServices(model.DownstreamServiceFilter{}, 1, 10)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if total != 2 {
		t.Fatalf("total mismatch")
	}
	if len(items) != 2 {
		t.Fatalf("items length mismatch")
	}
}

func TestService_CreateBreakerRule_ForeignKey(t *testing.T) {
	svc := newTestService()
	if _, err := svc.CreateBreakerRule(model.BreakerRule{Name: "rule1", ServiceID: "none"}); err == nil {
		t.Fatalf("expected validation error for missing service")
	}
	s, _ := svc.CreateDownstreamService(model.DownstreamService{Name: "svc1", Address: "127.0.0.1:8080"})
	r, err := svc.CreateBreakerRule(model.BreakerRule{Name: "rule1", ServiceID: s.ID})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if r.ServiceID != s.ID {
		t.Fatalf("service id mismatch")
	}
}

func TestService_CircuitBreakerStateMachine(t *testing.T) {
	svc := newTestService()
	s, _ := svc.CreateDownstreamService(model.DownstreamService{Name: "svc1", Address: "127.0.0.1:8080"})
	rule, _ := svc.CreateBreakerRule(model.BreakerRule{Name: "rule1", ServiceID: s.ID})
	b, err := svc.CreateCircuitBreaker(model.CircuitBreaker{ServiceID: s.ID, RuleID: rule.ID})
	if err != nil {
		t.Fatalf("create breaker failed: %v", err)
	}
	if b.State != model.BreakerStateClosed {
		t.Fatalf("initial state should be closed")
	}

	_, err = svc.UpdateCircuitBreaker(b.ID, model.CircuitBreaker{ServiceID: s.ID, RuleID: rule.ID, State: model.BreakerStateOpen})
	if err != nil {
		t.Fatalf("closed->open should be allowed: %v", err)
	}

	_, err = svc.UpdateCircuitBreaker(b.ID, model.CircuitBreaker{ServiceID: s.ID, RuleID: rule.ID, State: model.BreakerStateClosed})
	if err == nil {
		t.Fatalf("open->closed should be disallowed")
	}

	_, err = svc.UpdateCircuitBreaker(b.ID, model.CircuitBreaker{ServiceID: s.ID, RuleID: rule.ID, State: model.BreakerStateHalfOpen})
	if err != nil {
		t.Fatalf("open->half_open should be allowed: %v", err)
	}

	_, err = svc.UpdateCircuitBreaker(b.ID, model.CircuitBreaker{ServiceID: s.ID, RuleID: rule.ID, State: model.BreakerStateClosed})
	if err != nil {
		t.Fatalf("half_open->closed should be allowed: %v", err)
	}
}

func TestService_CreateCallRecord(t *testing.T) {
	svc := newTestService()
	if _, err := svc.CreateCallRecord(model.CallRecord{RequestID: "req1", ServiceID: "none", Outcome: model.CallOutcomeSuccess}); err == nil {
		t.Fatalf("expected validation error for missing service")
	}
	s, _ := svc.CreateDownstreamService(model.DownstreamService{Name: "svc1", Address: "127.0.0.1:8080"})
	rec, err := svc.CreateCallRecord(model.CallRecord{RequestID: "req1", ServiceID: s.ID, Outcome: model.CallOutcomeSuccess, LatencyMs: 100})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if rec.ServiceID != s.ID {
		t.Fatalf("service id mismatch")
	}
}

func TestService_CreateHealthCheck(t *testing.T) {
	svc := newTestService()
	if _, err := svc.CreateHealthCheck(model.HealthCheck{ServiceID: "none"}); err == nil {
		t.Fatalf("expected validation error")
	}
	s, _ := svc.CreateDownstreamService(model.DownstreamService{Name: "svc1", Address: "127.0.0.1:8080"})
	hc, err := svc.CreateHealthCheck(model.HealthCheck{ServiceID: s.ID, IntervalSeconds: 10})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if hc.ServiceID != s.ID {
		t.Fatalf("service id mismatch")
	}
}

func TestService_CreateAlertRule(t *testing.T) {
	svc := newTestService()
	if _, err := svc.CreateAlertRule(model.AlertRule{Name: "alert1", ServiceID: "none", Metric: model.AlertMetricStateChanged}); err == nil {
		t.Fatalf("expected validation error")
	}
	s, _ := svc.CreateDownstreamService(model.DownstreamService{Name: "svc1", Address: "127.0.0.1:8080"})
	ar, err := svc.CreateAlertRule(model.AlertRule{Name: "alert1", ServiceID: s.ID, Metric: model.AlertMetricStateChanged, Threshold: 1, Severity: model.AlertSeverityWarn, NotifyChannel: model.AlertChannelWebhook})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if ar.ServiceID != s.ID {
		t.Fatalf("service id mismatch")
	}
}

func TestService_CreateRecoveryPolicy(t *testing.T) {
	svc := newTestService()
	if _, err := svc.CreateRecoveryPolicy(model.RecoveryPolicy{Name: "policy1", ServiceID: "none"}); err == nil {
		t.Fatalf("expected validation error")
	}
	s, _ := svc.CreateDownstreamService(model.DownstreamService{Name: "svc1", Address: "127.0.0.1:8080"})
	p, err := svc.CreateRecoveryPolicy(model.RecoveryPolicy{Name: "policy1", ServiceID: s.ID, HalfOpenProbeRatio: 0.1})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if p.ServiceID != s.ID {
		t.Fatalf("service id mismatch")
	}
}

func TestService_CreateMetricSample(t *testing.T) {
	svc := newTestService()
	if _, err := svc.CreateMetricSample(model.MetricSample{ServiceID: "none", WindowStart: time.Now(), WindowEnd: time.Now().Add(time.Minute)}); err == nil {
		t.Fatalf("expected validation error")
	}
	s, _ := svc.CreateDownstreamService(model.DownstreamService{Name: "svc1", Address: "127.0.0.1:8080"})
	m, err := svc.CreateMetricSample(model.MetricSample{ServiceID: s.ID, WindowStart: time.Now(), WindowEnd: time.Now().Add(time.Minute), TotalCalls: 100})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if m.ServiceID != s.ID {
		t.Fatalf("service id mismatch")
	}
}

func TestService_CreateBreakerEvent(t *testing.T) {
	svc := newTestService()
	if _, err := svc.CreateBreakerEvent(model.BreakerEvent{BreakerID: "none", ServiceID: "s1", EventType: model.EventTypeOpened}); err == nil {
		t.Fatalf("expected validation error")
	}
	s, _ := svc.CreateDownstreamService(model.DownstreamService{Name: "svc1", Address: "127.0.0.1:8080"})
	rule, _ := svc.CreateBreakerRule(model.BreakerRule{Name: "rule1", ServiceID: s.ID})
	b, _ := svc.CreateCircuitBreaker(model.CircuitBreaker{ServiceID: s.ID, RuleID: rule.ID})
	e, err := svc.CreateBreakerEvent(model.BreakerEvent{BreakerID: b.ID, ServiceID: s.ID, EventType: model.EventTypeOpened, Reason: "test"})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if e.BreakerID != b.ID {
		t.Fatalf("breaker id mismatch")
	}
}

func TestService_CreateBreakerSnapshot(t *testing.T) {
	svc := newTestService()
	if _, err := svc.CreateBreakerSnapshot(model.BreakerSnapshot{ServiceID: "none", State: model.BreakerStateClosed}); err == nil {
		t.Fatalf("expected validation error")
	}
	s, _ := svc.CreateDownstreamService(model.DownstreamService{Name: "svc1", Address: "127.0.0.1:8080"})
	snap, err := svc.CreateBreakerSnapshot(model.BreakerSnapshot{ServiceID: s.ID, State: model.BreakerStateClosed, TotalCalls: 100})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if snap.ServiceID != s.ID {
		t.Fatalf("service id mismatch")
	}
}

func TestService_StatsOverview(t *testing.T) {
	svc := newTestService()
	s1, _ := svc.CreateDownstreamService(model.DownstreamService{Name: "svc1", Address: "127.0.0.1:8080"})
	s2, _ := svc.CreateDownstreamService(model.DownstreamService{Name: "svc2", Address: "127.0.0.1:8081"})
	rule, _ := svc.CreateBreakerRule(model.BreakerRule{Name: "rule1", ServiceID: s1.ID})
	_, _ = svc.CreateCircuitBreaker(model.CircuitBreaker{ServiceID: s1.ID, RuleID: rule.ID, State: model.BreakerStateOpen})
	_, _ = svc.CreateCallRecord(model.CallRecord{RequestID: "r1", ServiceID: s1.ID, Outcome: model.CallOutcomeFailure, LatencyMs: 100})
	_, _ = svc.CreateCallRecord(model.CallRecord{RequestID: "r2", ServiceID: s2.ID, Outcome: model.CallOutcomeSuccess, LatencyMs: 50})

	stats := svc.GetStatsOverview()
	if stats.TotalServices != 2 {
		t.Fatalf("total services mismatch")
	}
	if stats.TotalCalls != 2 {
		t.Fatalf("total calls mismatch")
	}
	if stats.OpenBreakers != 1 {
		t.Fatalf("open breakers mismatch")
	}
}

func TestService_TopFailureServices(t *testing.T) {
	svc := newTestService()
	s1, _ := svc.CreateDownstreamService(model.DownstreamService{Name: "svc1", Address: "127.0.0.1:8080"})
	_, _ = svc.CreateCallRecord(model.CallRecord{RequestID: "r1", ServiceID: s1.ID, Outcome: model.CallOutcomeFailure, LatencyMs: 100})
	_, _ = svc.CreateCallRecord(model.CallRecord{RequestID: "r2", ServiceID: s1.ID, Outcome: model.CallOutcomeFailure, LatencyMs: 100})

	top := svc.GetTopFailureServices(5)
	if len(top) != 1 {
		t.Fatalf("top length mismatch")
	}
	if top[0].FailureCount != 2 {
		t.Fatalf("failure count mismatch")
	}
}

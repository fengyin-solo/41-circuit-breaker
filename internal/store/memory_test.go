package store

import (
	"testing"
	"time"

	"circuitbreaker/internal/model"
)

func TestMemoryStore_DownstreamService(t *testing.T) {
	st := NewMemoryStore()
	svc := &model.DownstreamService{ID: "s1", Name: "svc1", Address: "127.0.0.1:8080", Protocol: model.ServiceProtocolHTTP, TimeoutMs: 3000, Status: model.ServiceStatusUp, Weight: 10, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := st.CreateDownstreamService(svc); err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if err := st.CreateDownstreamService(&model.DownstreamService{ID: "s2", Name: "svc1", Address: "127.0.0.1:8081", Protocol: model.ServiceProtocolHTTP, TimeoutMs: 3000, Status: model.ServiceStatusUp, Weight: 10, CreatedAt: time.Now(), UpdatedAt: time.Now()}); err != ErrConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
	got, err := st.GetDownstreamService("s1")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if got.Name != "svc1" {
		t.Fatalf("name mismatch")
	}
	if _, err := st.GetDownstreamService("none"); err != ErrNotFound {
		t.Fatalf("expected not found")
	}
	list := st.ListDownstreamServices()
	if len(list) != 1 {
		t.Fatalf("list length mismatch")
	}
	svc.Name = "svc1-updated"
	if err := st.UpdateDownstreamService(svc); err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if err := st.DeleteDownstreamService("s1"); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	if _, err := st.GetDownstreamService("s1"); err != ErrNotFound {
		t.Fatalf("expected not found after delete")
	}
}

func TestMemoryStore_BreakerRule(t *testing.T) {
	st := NewMemoryStore()
	r := &model.BreakerRule{ID: "r1", Name: "rule1", ServiceID: "s1", FailureRatioThreshold: 0.5, SlowCallRatioThreshold: 0.5, SlowCallMs: 1000, WindowSeconds: 60, MinRequestCount: 10, MaxHalfOpenRequests: 5, Enabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := st.CreateBreakerRule(r); err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if err := st.CreateBreakerRule(&model.BreakerRule{ID: "r2", Name: "rule1", ServiceID: "s1"}); err != ErrConflict {
		t.Fatalf("expected conflict")
	}
	got, err := st.GetBreakerRule("r1")
	if err != nil || got.Name != "rule1" {
		t.Fatalf("get failed")
	}
	if len(st.ListBreakerRules()) != 1 {
		t.Fatalf("list mismatch")
	}
	r.Name = "rule1-updated"
	if err := st.UpdateBreakerRule(r); err != nil {
		t.Fatalf("update failed")
	}
	if err := st.DeleteBreakerRule("r1"); err != nil {
		t.Fatalf("delete failed")
	}
	if _, err := st.GetBreakerRule("r1"); err != ErrNotFound {
		t.Fatalf("expected not found")
	}
}

func TestMemoryStore_CircuitBreaker(t *testing.T) {
	st := NewMemoryStore()
	b := &model.CircuitBreaker{ID: "b1", ServiceID: "s1", RuleID: "r1", State: model.BreakerStateClosed, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := st.CreateCircuitBreaker(b); err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if err := st.CreateCircuitBreaker(&model.CircuitBreaker{ID: "b2", ServiceID: "s1", RuleID: "r1"}); err != ErrConflict {
		t.Fatalf("expected conflict")
	}
	got, err := st.GetCircuitBreaker("b1")
	if err != nil || got.State != model.BreakerStateClosed {
		t.Fatalf("get failed")
	}
	if _, err := st.GetCircuitBreakerByServiceAndRule("s1", "r1"); err != nil {
		t.Fatalf("get by service and rule failed")
	}
	b.State = model.BreakerStateOpen
	if err := st.UpdateCircuitBreaker(b); err != nil {
		t.Fatalf("update failed")
	}
	if err := st.DeleteCircuitBreaker("b1"); err != nil {
		t.Fatalf("delete failed")
	}
}

func TestMemoryStore_CallRecord(t *testing.T) {
	st := NewMemoryStore()
	c := &model.CallRecord{ID: "c1", RequestID: "req1", ServiceID: "s1", Outcome: model.CallOutcomeSuccess, LatencyMs: 100, CalledAt: time.Now()}
	if err := st.CreateCallRecord(c); err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if err := st.CreateCallRecord(&model.CallRecord{ID: "c2", RequestID: "req1", ServiceID: "s1"}); err != ErrConflict {
		t.Fatalf("expected conflict")
	}
	got, err := st.GetCallRecord("c1")
	if err != nil || got.RequestID != "req1" {
		t.Fatalf("get failed")
	}
	if len(st.ListCallRecords()) != 1 {
		t.Fatalf("list mismatch")
	}
	if err := st.DeleteCallRecord("c1"); err != nil {
		t.Fatalf("delete failed")
	}
}

func TestMemoryStore_HealthCheck(t *testing.T) {
	st := NewMemoryStore()
	h := &model.HealthCheck{ID: "h1", ServiceID: "s1", IntervalSeconds: 10, LastStatus: model.HealthCheckStatusHealthy, ConsecutiveFailures: 0, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := st.CreateHealthCheck(h); err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if err := st.CreateHealthCheck(&model.HealthCheck{ID: "h2", ServiceID: "s1"}); err != ErrConflict {
		t.Fatalf("expected conflict")
	}
	got, err := st.GetHealthCheck("h1")
	if err != nil || got.ServiceID != "s1" {
		t.Fatalf("get failed")
	}
	if _, err := st.GetHealthCheckByServiceID("s1"); err != nil {
		t.Fatalf("get by service id failed")
	}
	h.LastStatus = model.HealthCheckStatusUnhealthy
	if err := st.UpdateHealthCheck(h); err != nil {
		t.Fatalf("update failed")
	}
	if err := st.DeleteHealthCheck("h1"); err != nil {
		t.Fatalf("delete failed")
	}
}

func TestMemoryStore_AlertRule(t *testing.T) {
	st := NewMemoryStore()
	a := &model.AlertRule{ID: "a1", Name: "alert1", ServiceID: "s1", Metric: model.AlertMetricStateChanged, Threshold: 1, Severity: model.AlertSeverityWarn, NotifyChannel: model.AlertChannelWebhook, Enabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := st.CreateAlertRule(a); err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if err := st.CreateAlertRule(&model.AlertRule{ID: "a2", Name: "alert1", ServiceID: "s1"}); err != ErrConflict {
		t.Fatalf("expected conflict")
	}
	got, err := st.GetAlertRule("a1")
	if err != nil || got.Name != "alert1" {
		t.Fatalf("get failed")
	}
	if len(st.ListAlertRules()) != 1 {
		t.Fatalf("list mismatch")
	}
	a.Enabled = false
	if err := st.UpdateAlertRule(a); err != nil {
		t.Fatalf("update failed")
	}
	if err := st.DeleteAlertRule("a1"); err != nil {
		t.Fatalf("delete failed")
	}
}

func TestMemoryStore_RecoveryPolicy(t *testing.T) {
	st := NewMemoryStore()
	p := &model.RecoveryPolicy{ID: "p1", Name: "policy1", ServiceID: "s1", HalfOpenProbeRatio: 0.1, RecoveryWindowSeconds: 30, MaxRetry: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := st.CreateRecoveryPolicy(p); err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if err := st.CreateRecoveryPolicy(&model.RecoveryPolicy{ID: "p2", Name: "policy1", ServiceID: "s1"}); err != ErrConflict {
		t.Fatalf("expected conflict")
	}
	got, err := st.GetRecoveryPolicy("p1")
	if err != nil || got.Name != "policy1" {
		t.Fatalf("get failed")
	}
	if _, err := st.GetRecoveryPolicyByServiceID("s1"); err != nil {
		t.Fatalf("get by service id failed")
	}
	p.MaxRetry = 5
	if err := st.UpdateRecoveryPolicy(p); err != nil {
		t.Fatalf("update failed")
	}
	if err := st.DeleteRecoveryPolicy("p1"); err != nil {
		t.Fatalf("delete failed")
	}
}

func TestMemoryStore_MetricSample(t *testing.T) {
	st := NewMemoryStore()
	m := &model.MetricSample{ID: "m1", ServiceID: "s1", WindowStart: time.Now(), WindowEnd: time.Now().Add(time.Minute), TotalCalls: 100, SuccessCalls: 90, FailureCalls: 5, SlowCalls: 5, FailureRatio: 0.05, SlowRatio: 0.05, CreatedAt: time.Now()}
	if err := st.CreateMetricSample(m); err != nil {
		t.Fatalf("create failed: %v", err)
	}
	got, err := st.GetMetricSample("m1")
	if err != nil || got.ServiceID != "s1" {
		t.Fatalf("get failed")
	}
	if len(st.ListMetricSamples()) != 1 {
		t.Fatalf("list mismatch")
	}
	if err := st.DeleteMetricSample("m1"); err != nil {
		t.Fatalf("delete failed")
	}
}

func TestMemoryStore_BreakerEvent(t *testing.T) {
	st := NewMemoryStore()
	e := &model.BreakerEvent{ID: "e1", BreakerID: "b1", ServiceID: "s1", EventType: model.EventTypeOpened, Reason: "test", OccurredAt: time.Now()}
	if err := st.CreateBreakerEvent(e); err != nil {
		t.Fatalf("create failed: %v", err)
	}
	got, err := st.GetBreakerEvent("e1")
	if err != nil || got.EventType != model.EventTypeOpened {
		t.Fatalf("get failed")
	}
	if len(st.ListBreakerEvents()) != 1 {
		t.Fatalf("list mismatch")
	}
	if err := st.DeleteBreakerEvent("e1"); err != nil {
		t.Fatalf("delete failed")
	}
}

func TestMemoryStore_BreakerSnapshot(t *testing.T) {
	st := NewMemoryStore()
	snap := &model.BreakerSnapshot{ID: "snap1", ServiceID: "s1", State: model.BreakerStateClosed, TotalCalls: 100, SuccessCalls: 95, FailureCalls: 5, SlowCalls: 0, FailureRatio: 0.05, SlowRatio: 0, SnapshotVersion: 1, CreatedAt: time.Now()}
	if err := st.CreateBreakerSnapshot(snap); err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if err := st.CreateBreakerSnapshot(&model.BreakerSnapshot{ID: "snap2", ServiceID: "s1", SnapshotVersion: 1}); err != ErrConflict {
		t.Fatalf("expected conflict")
	}
	got, err := st.GetBreakerSnapshot("snap1")
	if err != nil || got.ServiceID != "s1" {
		t.Fatalf("get failed")
	}
	if len(st.ListBreakerSnapshots()) != 1 {
		t.Fatalf("list mismatch")
	}
	snap.SnapshotVersion = 2
	if err := st.UpdateBreakerSnapshot(snap); err != nil {
		t.Fatalf("update failed")
	}
	if err := st.DeleteBreakerSnapshot("snap1"); err != nil {
		t.Fatalf("delete failed")
	}
}

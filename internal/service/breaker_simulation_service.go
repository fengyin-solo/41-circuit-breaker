package service

import (
	"math/rand"
	"time"

	"circuitbreaker/internal/model"
	"circuitbreaker/pkg/idgen"
)

type SimulationResult struct {
	RequestID     string `json:"request_id"`
	ServiceID     string `json:"service_id"`
	Allowed       bool   `json:"allowed"`
	BreakerState  string `json:"breaker_state"`
	Outcome       string `json:"outcome"`
	LatencyMs     int    `json:"latency_ms"`
	Reason        string `json:"reason"`
}

func (s *Service) SimulateCall(serviceID string, forcedOutcome string, forcedLatencyMs int) (*SimulationResult, error) {
	service, err := s.store.GetDownstreamService(serviceID)
	if err != nil {
		return nil, model.NewValidationError("service_id", "下游服务不存在")
	}
	breakers := s.store.ListCircuitBreakers()
	var breaker *model.CircuitBreaker
	for _, b := range breakers {
		if b.ServiceID == serviceID {
			breaker = b
			break
		}
	}
	res := &SimulationResult{
		RequestID: idgen.HexN(8),
		ServiceID: serviceID,
	}
	if breaker == nil {
		res.Allowed = true
		res.BreakerState = "none"
		res.Reason = "无熔断器，直接放行"
	} else {
		res.BreakerState = breaker.State
		switch breaker.State {
		case model.BreakerStateOpen:
			res.Allowed = false
			res.Reason = "熔断器已打开，请求被拒绝"
			res.Outcome = model.CallOutcomeFailure
		case model.BreakerStateHalfOpen:
			policy, perr := s.store.GetRecoveryPolicyByServiceID(serviceID)
			if perr == nil && policy != nil {
				if rand.Float64() > policy.HalfOpenProbeRatio {
					res.Allowed = false
					res.Reason = "半开状态，探测比例未命中"
					res.Outcome = model.CallOutcomeFailure
					break
				}
			}
			res.Allowed = true
			res.Reason = "半开状态，允许探测请求"
		case model.BreakerStateClosed:
			res.Allowed = true
			res.Reason = "熔断器关闭，请求正常放行"
		}
	}
	if !res.Allowed {
		s.recordCallAndEvent(serviceID, res.RequestID, res.Outcome, 0, breaker)
		return res, nil
	}
	outcome := forcedOutcome
	latency := forcedLatencyMs
	if outcome == "" {
		outcome = s.randomOutcome(service)
	}
	if latency <= 0 {
		latency = s.randomLatency(service.TimeoutMs)
	}
	res.Outcome = outcome
	res.LatencyMs = latency
	s.recordCallAndEvent(serviceID, res.RequestID, outcome, latency, breaker)
	if breaker != nil {
		s.updateBreakerStats(breaker, outcome, latency)
	}
	return res, nil
}

func (s *Service) randomOutcome(service *model.DownstreamService) string {
	if service.Status == model.ServiceStatusDown {
		if rand.Float64() < 0.8 {
			return model.CallOutcomeFailure
		}
		return model.CallOutcomeTimeout
	}
	r := rand.Float64()
	if r < 0.85 {
		return model.CallOutcomeSuccess
	}
	if r < 0.95 {
		return model.CallOutcomeFailure
	}
	return model.CallOutcomeTimeout
}

func (s *Service) randomLatency(timeoutMs int) int {
	if timeoutMs <= 0 {
		timeoutMs = 5000
	}
	base := rand.Intn(timeoutMs / 2)
	if rand.Float64() < 0.1 {
		base += timeoutMs / 2
	}
	return base + rand.Intn(50)
}

func (s *Service) recordCallAndEvent(serviceID, requestID, outcome string, latencyMs int, breaker *model.CircuitBreaker) {
	rec := model.CallRecord{
		ID:        idgen.Hex(),
		RequestID: requestID,
		ServiceID: serviceID,
		Outcome:   outcome,
		LatencyMs: latencyMs,
		CalledAt:  time.Now(),
	}
	_ = s.store.CreateCallRecord(&rec)
	if breaker != nil && (outcome == model.CallOutcomeFailure || outcome == model.CallOutcomeTimeout) {
		if breaker.FailureCount >= 5 {
			event := model.BreakerEvent{
				ID:         idgen.Hex(),
				BreakerID:  breaker.ID,
				ServiceID:  serviceID,
				EventType:  model.EventTypeRejected,
				Reason:     "模拟调用失败，请求被拒绝",
				OccurredAt: time.Now(),
			}
			_ = s.store.CreateBreakerEvent(&event)
		}
	}
}

func (s *Service) updateBreakerStats(breaker *model.CircuitBreaker, outcome string, latencyMs int) {
	breaker.TotalCalls++
	switch outcome {
	case model.CallOutcomeSuccess:
		breaker.SuccessCount++
	case model.CallOutcomeFailure:
		breaker.FailureCount++
	case model.CallOutcomeTimeout:
		breaker.FailureCount++
	}
	rule, err := s.store.GetBreakerRule(breaker.RuleID)
	if err == nil && rule != nil {
		if latencyMs >= rule.SlowCallMs {
			if breaker.TotalCalls > 0 {
			}
		}
	}
	breaker.UpdatedAt = time.Now()
	_ = s.store.UpdateCircuitBreaker(breaker)
}

func (s *Service) SimulateBulkCalls(serviceID string, count int) []*SimulationResult {
	results := make([]*SimulationResult, 0, count)
	for i := 0; i < count; i++ {
		res, err := s.SimulateCall(serviceID, "", 0)
		if err != nil {
			continue
		}
		results = append(results, res)
	}
	return results
}

type SimulationScenarioResult struct {
	TotalRequests   int     `json:"total_requests"`
	AllowedRequests int     `json:"allowed_requests"`
	RejectedRequests int    `json:"rejected_requests"`
	SuccessCalls    int     `json:"success_calls"`
	FailureCalls    int     `json:"failure_calls"`
	TimeoutCalls    int     `json:"timeout_calls"`
	AvgLatencyMs    int     `json:"avg_latency_ms"`
	FinalBreakerState string `json:"final_breaker_state"`
}

func (s *Service) RunScenario(serviceID string, requestCount int, failureRate float64, avgLatencyMs int) (*SimulationScenarioResult, error) {
	service, err := s.store.GetDownstreamService(serviceID)
	if err != nil {
		return nil, model.NewValidationError("service_id", "下游服务不存在")
	}
	breakers := s.store.ListCircuitBreakers()
	var breaker *model.CircuitBreaker
	for _, b := range breakers {
		if b.ServiceID == serviceID {
			breaker = b
			break
		}
	}
	if breaker == nil {
		return nil, model.NewValidationError("breaker", "该服务未配置熔断器")
	}
	result := &SimulationScenarioResult{TotalRequests: requestCount}
	totalLatency := 0
	for i := 0; i < requestCount; i++ {
		if breaker.State == model.BreakerStateOpen {
			result.RejectedRequests++
			continue
		}
		outcome := model.CallOutcomeSuccess
		if rand.Float64() < failureRate {
			if rand.Float64() < 0.5 {
				outcome = model.CallOutcomeFailure
			} else {
				outcome = model.CallOutcomeTimeout
			}
		}
		latency := avgLatencyMs + rand.Intn(100) - 50
		if latency < 0 {
			latency = 0
		}
		res, _ := s.SimulateCall(serviceID, outcome, latency)
		if res.Allowed {
			result.AllowedRequests++
			totalLatency += res.LatencyMs
		} else {
			result.RejectedRequests++
		}
		switch res.Outcome {
		case model.CallOutcomeSuccess:
			result.SuccessCalls++
		case model.CallOutcomeFailure:
			result.FailureCalls++
		case model.CallOutcomeTimeout:
			result.TimeoutCalls++
		}
		breaker, _ = s.store.GetCircuitBreaker(breaker.ID)
	}
	if result.AllowedRequests > 0 {
		result.AvgLatencyMs = totalLatency / result.AllowedRequests
	}
	if breaker != nil {
		result.FinalBreakerState = breaker.State
	}
	_ = service
	return result, nil
}

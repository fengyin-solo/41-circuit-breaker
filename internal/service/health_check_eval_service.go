package service

import (
	"time"

	"circuitbreaker/internal/model"
	"circuitbreaker/pkg/idgen"
)

type HealthEvalResult struct {
	ServiceID           string `json:"service_id"`
	PreviousStatus      string `json:"previous_status"`
	CurrentStatus       string `json:"current_status"`
	ConsecutiveFailures int    `json:"consecutive_failures"`
	ActionTaken         string `json:"action_taken"`
}

func (s *Service) EvaluateHealthChecks() []*HealthEvalResult {
	checks := s.store.ListHealthChecks()
	result := make([]*HealthEvalResult, 0, len(checks))
	for _, hc := range checks {
		service, err := s.store.GetDownstreamService(hc.ServiceID)
		if err != nil {
			continue
		}
		res := &HealthEvalResult{
			ServiceID:      service.ID,
			PreviousStatus: hc.LastStatus,
			CurrentStatus:  hc.LastStatus,
		}
		if service.Status == model.ServiceStatusDown {
			hc.LastStatus = model.HealthCheckStatusUnhealthy
			hc.ConsecutiveFailures++
			res.CurrentStatus = model.HealthCheckStatusUnhealthy
			res.ConsecutiveFailures = hc.ConsecutiveFailures
			res.ActionTaken = "service_down"
		} else {
			hc.LastStatus = model.HealthCheckStatusHealthy
			if hc.ConsecutiveFailures > 0 {
				res.ActionTaken = "recovered"
			}
			hc.ConsecutiveFailures = 0
		}
		now := time.Now()
		hc.LastCheckedAt = &now
		hc.UpdatedAt = now
		_ = s.store.UpdateHealthCheck(hc)
		s.autoAdjustBreakerByHealth(service.ID, hc.LastStatus, hc.ConsecutiveFailures)
		result = append(result, res)
	}
	return result
}

func (s *Service) autoAdjustBreakerByHealth(serviceID, healthStatus string, consecutiveFailures int) {
	breakers := s.store.ListCircuitBreakers()
	for _, b := range breakers {
		if b.ServiceID != serviceID {
			continue
		}
		if healthStatus == model.HealthCheckStatusUnhealthy && consecutiveFailures >= 3 && b.State == model.BreakerStateClosed {
			_, _ = s.UpdateCircuitBreaker(b.ID, model.CircuitBreaker{
				ServiceID: b.ServiceID,
				RuleID:    b.RuleID,
				State:     model.BreakerStateOpen,
			})
			s.log.Warnf("健康检查连续失败 %d 次，自动熔断服务 %s", consecutiveFailures, serviceID)
		}
		if healthStatus == model.HealthCheckStatusHealthy && b.State == model.BreakerStateOpen {
			_, _ = s.UpdateCircuitBreaker(b.ID, model.CircuitBreaker{
				ServiceID: b.ServiceID,
				RuleID:    b.RuleID,
				State:     model.BreakerStateHalfOpen,
			})
			s.log.Infof("服务 %s 恢复健康，自动进入半开探测", serviceID)
		}
	}
}

func (s *Service) EvaluateSingleHealthCheck(serviceID string) (*HealthEvalResult, error) {
	hc, err := s.store.GetHealthCheckByServiceID(serviceID)
	if err != nil {
		return nil, err
	}
	service, err := s.store.GetDownstreamService(hc.ServiceID)
	if err != nil {
		return nil, err
	}
	res := &HealthEvalResult{
		ServiceID:      service.ID,
		PreviousStatus: hc.LastStatus,
		CurrentStatus:  hc.LastStatus,
	}
	if service.Status == model.ServiceStatusDown {
		hc.LastStatus = model.HealthCheckStatusUnhealthy
		hc.ConsecutiveFailures++
		res.CurrentStatus = model.HealthCheckStatusUnhealthy
		res.ConsecutiveFailures = hc.ConsecutiveFailures
		res.ActionTaken = "service_down"
	} else {
		hc.LastStatus = model.HealthCheckStatusHealthy
		if hc.ConsecutiveFailures > 0 {
			res.ActionTaken = "recovered"
		}
		hc.ConsecutiveFailures = 0
	}
	now := time.Now()
	hc.LastCheckedAt = &now
	hc.UpdatedAt = now
	_ = s.store.UpdateHealthCheck(hc)
	s.autoAdjustBreakerByHealth(service.ID, hc.LastStatus, hc.ConsecutiveFailures)
	return res, nil
}

func (s *Service) PerformManualHealthCheck(serviceID string, success bool) (*HealthEvalResult, error) {
	hc, err := s.store.GetHealthCheckByServiceID(serviceID)
	if err != nil {
		return nil, err
	}
	res := &HealthEvalResult{
		ServiceID:      serviceID,
		PreviousStatus: hc.LastStatus,
	}
	now := time.Now()
	hc.LastCheckedAt = &now
	hc.UpdatedAt = now
	if success {
		hc.LastStatus = model.HealthCheckStatusHealthy
		if hc.ConsecutiveFailures > 0 {
			res.ActionTaken = "manual_recovered"
		}
		hc.ConsecutiveFailures = 0
	} else {
		hc.LastStatus = model.HealthCheckStatusUnhealthy
		hc.ConsecutiveFailures++
		res.ActionTaken = "manual_failure"
	}
	res.CurrentStatus = hc.LastStatus
	res.ConsecutiveFailures = hc.ConsecutiveFailures
	_ = s.store.UpdateHealthCheck(hc)
	s.autoAdjustBreakerByHealth(serviceID, hc.LastStatus, hc.ConsecutiveFailures)
	s.recordHealthCheckEvent(serviceID, success, hc.ConsecutiveFailures)
	return res, nil
}

func (s *Service) recordHealthCheckEvent(serviceID string, success bool, consecutiveFailures int) {
	reason := "健康检查通过"
	if !success {
		reason = "健康检查失败"
	}
	if consecutiveFailures > 0 {
		reason += "，连续失败次数: " + string(rune('0'+consecutiveFailures))
	}
	breakers := s.store.ListCircuitBreakers()
	for _, b := range breakers {
		if b.ServiceID == serviceID {
			eventType := model.EventTypeClosed
			if !success {
				eventType = model.EventTypeOpened
			}
			event := model.BreakerEvent{
				ID:         idgen.Hex(),
				BreakerID:  b.ID,
				ServiceID:  serviceID,
				EventType:  eventType,
				Reason:     reason,
				OccurredAt: time.Now(),
			}
			_ = s.store.CreateBreakerEvent(&event)
			break
		}
	}
}

package service

import (
	"sort"
	"time"

	"circuitbreaker/internal/model"
	"circuitbreaker/pkg/idgen"
)

func (s *Service) CreateCircuitBreaker(input model.CircuitBreaker) (*model.CircuitBreaker, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetDownstreamService(input.ServiceID); err != nil {
		return nil, model.NewValidationError("service_id", "关联的下游服务不存在")
	}
	if _, err := s.store.GetBreakerRule(input.RuleID); err != nil {
		return nil, model.NewValidationError("rule_id", "关联的熔断规则不存在")
	}
	input.ID = idgen.Hex()
	input.CreatedAt = time.Now()
	input.UpdatedAt = input.CreatedAt
	if err := s.store.CreateCircuitBreaker(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetCircuitBreaker(id string) (*model.CircuitBreaker, error) {
	return s.store.GetCircuitBreaker(id)
}

func (s *Service) ListCircuitBreakers(filter model.CircuitBreakerFilter, page, size int) ([]*model.CircuitBreaker, int, error) {
	all := s.store.ListCircuitBreakers()
	matched := make([]*model.CircuitBreaker, 0, len(all))
	for _, b := range all {
		if filter.Match(b) {
			matched = append(matched, b)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.CircuitBreaker{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateCircuitBreaker(id string, input model.CircuitBreaker) (*model.CircuitBreaker, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	exist, err := s.store.GetCircuitBreaker(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionBreaker(exist.State, input.State) {
		return nil, model.NewValidationError("state", "非法的状态流转: "+exist.State+" -> "+input.State)
	}
	if input.ServiceID != exist.ServiceID {
		if _, err := s.store.GetDownstreamService(input.ServiceID); err != nil {
			return nil, model.NewValidationError("service_id", "关联的下游服务不存在")
		}
	}
	if input.RuleID != exist.RuleID {
		if _, err := s.store.GetBreakerRule(input.RuleID); err != nil {
			return nil, model.NewValidationError("rule_id", "关联的熔断规则不存在")
		}
	}
	oldState := exist.State
	exist.State = input.State
	exist.ServiceID = input.ServiceID
	exist.RuleID = input.RuleID
	exist.UpdatedAt = time.Now()
	now := time.Now()
	if input.State == model.BreakerStateOpen {
		exist.LastOpenedAt = &now
		s.log.Infof("熔断器 %s 进入 open 状态", exist.ID)
	}
	if input.State == model.BreakerStateClosed {
		exist.LastClosedAt = &now
		s.log.Infof("熔断器 %s 进入 closed 状态", exist.ID)
	}
	if err := s.store.UpdateCircuitBreaker(exist); err != nil {
		return nil, err
	}
	s.recordBreakerEvent(exist.ID, exist.ServiceID, input.State, oldState)
	s.checkAlertRules(exist.ServiceID, input.State, exist.FailureCount, exist.TotalCalls)
	return exist, nil
}

func (s *Service) recordBreakerEvent(breakerID, serviceID, newState, oldState string) {
	var eventType string
	var reason string
	switch newState {
	case model.BreakerStateOpen:
		eventType = model.EventTypeOpened
		reason = "触发熔断阈值"
	case model.BreakerStateClosed:
		eventType = model.EventTypeClosed
		reason = "服务恢复"
	case model.BreakerStateHalfOpen:
		eventType = model.EventTypeHalfOpen
		reason = "进入半开探测"
	default:
		return
	}
	event := model.BreakerEvent{
		ID:         idgen.Hex(),
		BreakerID:  breakerID,
		ServiceID:  serviceID,
		EventType:  eventType,
		Reason:     reason + " (" + oldState + " -> " + newState + ")",
		OccurredAt: time.Now(),
	}
	_ = s.store.CreateBreakerEvent(&event)
}

func (s *Service) checkAlertRules(serviceID, state string, failureCount, totalCalls int) {
	rules := s.store.ListAlertRules()
	for _, rule := range rules {
		if !rule.Enabled || rule.ServiceID != serviceID {
			continue
		}
		triggered := false
		switch rule.Metric {
		case model.AlertMetricStateChanged:
			if state == model.BreakerStateOpen || state == model.BreakerStateHalfOpen {
				triggered = true
			}
		case model.AlertMetricFailureRateHigh:
			if totalCalls > 0 && float64(failureCount)/float64(totalCalls) >= rule.Threshold {
				triggered = true
			}
		case model.AlertMetricConsecutiveFailure:
			if float64(failureCount) >= rule.Threshold {
				triggered = true
			}
		}
		if triggered {
			s.log.Warnf("告警触发: service=%s metric=%s severity=%s channel=%s", serviceID, rule.Metric, rule.Severity, rule.NotifyChannel)
		}
	}
}

func (s *Service) DeleteCircuitBreaker(id string) error {
	return s.store.DeleteCircuitBreaker(id)
}

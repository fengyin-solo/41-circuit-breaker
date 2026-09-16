package service

import (
	"sort"
	"time"

	"circuitbreaker/internal/model"
	"circuitbreaker/pkg/idgen"
)

func (s *Service) CreateBreakerRule(input model.BreakerRule) (*model.BreakerRule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetDownstreamService(input.ServiceID); err != nil {
		return nil, model.NewValidationError("service_id", "关联的下游服务不存在")
	}
	input.ID = idgen.Hex()
	input.CreatedAt = time.Now()
	input.UpdatedAt = input.CreatedAt
	if err := s.store.CreateBreakerRule(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetBreakerRule(id string) (*model.BreakerRule, error) {
	return s.store.GetBreakerRule(id)
}

func (s *Service) ListBreakerRules(filter model.BreakerRuleFilter, page, size int) ([]*model.BreakerRule, int, error) {
	all := s.store.ListBreakerRules()
	matched := make([]*model.BreakerRule, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.BreakerRule{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateBreakerRule(id string, input model.BreakerRule) (*model.BreakerRule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	exist, err := s.store.GetBreakerRule(id)
	if err != nil {
		return nil, err
	}
	if input.ServiceID != exist.ServiceID {
		if _, err := s.store.GetDownstreamService(input.ServiceID); err != nil {
			return nil, model.NewValidationError("service_id", "关联的下游服务不存在")
		}
	}
	exist.Name = input.Name
	exist.ServiceID = input.ServiceID
	exist.FailureRatioThreshold = input.FailureRatioThreshold
	exist.SlowCallRatioThreshold = input.SlowCallRatioThreshold
	exist.SlowCallMs = input.SlowCallMs
	exist.WindowSeconds = input.WindowSeconds
	exist.MinRequestCount = input.MinRequestCount
	exist.MaxHalfOpenRequests = input.MaxHalfOpenRequests
	exist.Enabled = input.Enabled
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateBreakerRule(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

func (s *Service) DeleteBreakerRule(id string) error {
	return s.store.DeleteBreakerRule(id)
}

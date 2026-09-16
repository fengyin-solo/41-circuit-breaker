package service

import (
	"sort"
	"time"

	"circuitbreaker/internal/model"
	"circuitbreaker/pkg/idgen"
)

func (s *Service) CreateHealthCheck(input model.HealthCheck) (*model.HealthCheck, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetDownstreamService(input.ServiceID); err != nil {
		return nil, model.NewValidationError("service_id", "关联的下游服务不存在")
	}
	input.ID = idgen.Hex()
	input.CreatedAt = time.Now()
	input.UpdatedAt = input.CreatedAt
	if err := s.store.CreateHealthCheck(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetHealthCheck(id string) (*model.HealthCheck, error) {
	return s.store.GetHealthCheck(id)
}

func (s *Service) ListHealthChecks(filter model.HealthCheckFilter, page, size int) ([]*model.HealthCheck, int, error) {
	all := s.store.ListHealthChecks()
	matched := make([]*model.HealthCheck, 0, len(all))
	for _, h := range all {
		if filter.Match(h) {
			matched = append(matched, h)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.HealthCheck{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateHealthCheck(id string, input model.HealthCheck) (*model.HealthCheck, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	exist, err := s.store.GetHealthCheck(id)
	if err != nil {
		return nil, err
	}
	if input.ServiceID != exist.ServiceID {
		if _, err := s.store.GetDownstreamService(input.ServiceID); err != nil {
			return nil, model.NewValidationError("service_id", "关联的下游服务不存在")
		}
	}
	exist.ServiceID = input.ServiceID
	exist.IntervalSeconds = input.IntervalSeconds
	exist.LastCheckedAt = input.LastCheckedAt
	exist.LastStatus = input.LastStatus
	exist.ConsecutiveFailures = input.ConsecutiveFailures
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateHealthCheck(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

func (s *Service) DeleteHealthCheck(id string) error {
	return s.store.DeleteHealthCheck(id)
}

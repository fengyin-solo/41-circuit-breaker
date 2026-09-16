package service

import (
	"sort"
	"time"

	"circuitbreaker/internal/model"
	"circuitbreaker/pkg/idgen"
)

func (s *Service) CreateRecoveryPolicy(input model.RecoveryPolicy) (*model.RecoveryPolicy, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetDownstreamService(input.ServiceID); err != nil {
		return nil, model.NewValidationError("service_id", "关联的下游服务不存在")
	}
	input.ID = idgen.Hex()
	input.CreatedAt = time.Now()
	input.UpdatedAt = input.CreatedAt
	if err := s.store.CreateRecoveryPolicy(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetRecoveryPolicy(id string) (*model.RecoveryPolicy, error) {
	return s.store.GetRecoveryPolicy(id)
}

func (s *Service) ListRecoveryPolicies(filter model.RecoveryPolicyFilter, page, size int) ([]*model.RecoveryPolicy, int, error) {
	all := s.store.ListRecoveryPolicies()
	matched := make([]*model.RecoveryPolicy, 0, len(all))
	for _, p := range all {
		if filter.Match(p) {
			matched = append(matched, p)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.RecoveryPolicy{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateRecoveryPolicy(id string, input model.RecoveryPolicy) (*model.RecoveryPolicy, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	exist, err := s.store.GetRecoveryPolicy(id)
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
	exist.HalfOpenProbeRatio = input.HalfOpenProbeRatio
	exist.RecoveryWindowSeconds = input.RecoveryWindowSeconds
	exist.MaxRetry = input.MaxRetry
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateRecoveryPolicy(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

func (s *Service) DeleteRecoveryPolicy(id string) error {
	return s.store.DeleteRecoveryPolicy(id)
}

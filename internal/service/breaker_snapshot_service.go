package service

import (
	"sort"
	"time"

	"circuitbreaker/internal/model"
	"circuitbreaker/pkg/idgen"
)

func (s *Service) CreateBreakerSnapshot(input model.BreakerSnapshot) (*model.BreakerSnapshot, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetDownstreamService(input.ServiceID); err != nil {
		return nil, model.NewValidationError("service_id", "关联的下游服务不存在")
	}
	input.ID = idgen.Hex()
	input.CreatedAt = time.Now()
	if err := s.store.CreateBreakerSnapshot(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetBreakerSnapshot(id string) (*model.BreakerSnapshot, error) {
	return s.store.GetBreakerSnapshot(id)
}

func (s *Service) ListBreakerSnapshots(filter model.BreakerSnapshotFilter, page, size int) ([]*model.BreakerSnapshot, int, error) {
	all := s.store.ListBreakerSnapshots()
	matched := make([]*model.BreakerSnapshot, 0, len(all))
	for _, snap := range all {
		if filter.Match(snap) {
			matched = append(matched, snap)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.BreakerSnapshot{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateBreakerSnapshot(id string, input model.BreakerSnapshot) (*model.BreakerSnapshot, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	exist, err := s.store.GetBreakerSnapshot(id)
	if err != nil {
		return nil, err
	}
	if input.ServiceID != exist.ServiceID {
		if _, err := s.store.GetDownstreamService(input.ServiceID); err != nil {
			return nil, model.NewValidationError("service_id", "关联的下游服务不存在")
		}
	}
	exist.ServiceID = input.ServiceID
	exist.State = input.State
	exist.TotalCalls = input.TotalCalls
	exist.SuccessCalls = input.SuccessCalls
	exist.FailureCalls = input.FailureCalls
	exist.SlowCalls = input.SlowCalls
	exist.FailureRatio = input.FailureRatio
	exist.SlowRatio = input.SlowRatio
	exist.SnapshotVersion = input.SnapshotVersion
	exist.CreatedAt = time.Now()
	if err := s.store.UpdateBreakerSnapshot(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

func (s *Service) DeleteBreakerSnapshot(id string) error {
	return s.store.DeleteBreakerSnapshot(id)
}

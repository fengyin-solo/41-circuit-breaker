package service

import (
	"sort"
	"time"

	"circuitbreaker/internal/model"
	"circuitbreaker/pkg/idgen"
)

func (s *Service) CreateBreakerEvent(input model.BreakerEvent) (*model.BreakerEvent, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetCircuitBreaker(input.BreakerID); err != nil {
		return nil, model.NewValidationError("breaker_id", "关联的熔断器不存在")
	}
	input.ID = idgen.Hex()
	if input.OccurredAt.IsZero() {
		input.OccurredAt = time.Now()
	}
	if err := s.store.CreateBreakerEvent(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetBreakerEvent(id string) (*model.BreakerEvent, error) {
	return s.store.GetBreakerEvent(id)
}

func (s *Service) ListBreakerEvents(filter model.BreakerEventFilter, page, size int) ([]*model.BreakerEvent, int, error) {
	all := s.store.ListBreakerEvents()
	matched := make([]*model.BreakerEvent, 0, len(all))
	for _, e := range all {
		if filter.Match(e) {
			matched = append(matched, e)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].OccurredAt.After(matched[j].OccurredAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.BreakerEvent{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) DeleteBreakerEvent(id string) error {
	return s.store.DeleteBreakerEvent(id)
}

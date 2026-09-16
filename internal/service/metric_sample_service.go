package service

import (
	"sort"
	"time"

	"circuitbreaker/internal/model"
	"circuitbreaker/pkg/idgen"
)

func (s *Service) CreateMetricSample(input model.MetricSample) (*model.MetricSample, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetDownstreamService(input.ServiceID); err != nil {
		return nil, model.NewValidationError("service_id", "关联的下游服务不存在")
	}
	input.ID = idgen.Hex()
	input.CreatedAt = time.Now()
	if err := s.store.CreateMetricSample(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetMetricSample(id string) (*model.MetricSample, error) {
	return s.store.GetMetricSample(id)
}

func (s *Service) ListMetricSamples(filter model.MetricSampleFilter, page, size int) ([]*model.MetricSample, int, error) {
	all := s.store.ListMetricSamples()
	matched := make([]*model.MetricSample, 0, len(all))
	for _, m := range all {
		if filter.Match(m) {
			matched = append(matched, m)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].WindowStart.After(matched[j].WindowStart)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.MetricSample{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) DeleteMetricSample(id string) error {
	return s.store.DeleteMetricSample(id)
}

package store

import (
	"circuitbreaker/internal/model"
)

func (s *MemoryStore) CreateMetricSample(m *model.MetricSample) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.metricSamples[m.ID] = m
	return nil
}

func (s *MemoryStore) GetMetricSample(id string) (*model.MetricSample, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.metricSamples[id]
	if !ok {
		return nil, ErrNotFound
	}
	return m, nil
}

func (s *MemoryStore) ListMetricSamples() []*model.MetricSample {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.MetricSample, 0, len(s.metricSamples))
	for _, m := range s.metricSamples {
		list = append(list, m)
	}
	return list
}

func (s *MemoryStore) DeleteMetricSample(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.metricSamples[id]; !ok {
		return ErrNotFound
	}
	delete(s.metricSamples, id)
	return nil
}

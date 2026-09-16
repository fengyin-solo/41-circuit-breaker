package store

import (
	"circuitbreaker/internal/model"
)

func (s *MemoryStore) CreateBreakerRule(r *model.BreakerRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.breakerRules {
		if exist.Name == r.Name && exist.ServiceID == r.ServiceID {
			return ErrConflict
		}
	}
	s.breakerRules[r.ID] = r
	return nil
}

func (s *MemoryStore) GetBreakerRule(id string) (*model.BreakerRule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.breakerRules[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) ListBreakerRules() []*model.BreakerRule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.BreakerRule, 0, len(s.breakerRules))
	for _, r := range s.breakerRules {
		list = append(list, r)
	}
	return list
}

func (s *MemoryStore) UpdateBreakerRule(r *model.BreakerRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.breakerRules[r.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.breakerRules {
		if exist.ID != r.ID && exist.Name == r.Name && exist.ServiceID == r.ServiceID {
			return ErrConflict
		}
	}
	s.breakerRules[r.ID] = r
	return nil
}

func (s *MemoryStore) DeleteBreakerRule(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.breakerRules[id]; !ok {
		return ErrNotFound
	}
	delete(s.breakerRules, id)
	return nil
}

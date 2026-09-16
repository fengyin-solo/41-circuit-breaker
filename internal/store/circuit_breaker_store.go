package store

import (
	"circuitbreaker/internal/model"
)

func (s *MemoryStore) CreateCircuitBreaker(b *model.CircuitBreaker) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.circuitBreakers {
		if exist.ServiceID == b.ServiceID && exist.RuleID == b.RuleID {
			return ErrConflict
		}
	}
	s.circuitBreakers[b.ID] = b
	return nil
}

func (s *MemoryStore) GetCircuitBreaker(id string) (*model.CircuitBreaker, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.circuitBreakers[id]
	if !ok {
		return nil, ErrNotFound
	}
	return b, nil
}

func (s *MemoryStore) GetCircuitBreakerByServiceAndRule(serviceID, ruleID string) (*model.CircuitBreaker, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, b := range s.circuitBreakers {
		if b.ServiceID == serviceID && b.RuleID == ruleID {
			return b, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListCircuitBreakers() []*model.CircuitBreaker {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.CircuitBreaker, 0, len(s.circuitBreakers))
	for _, b := range s.circuitBreakers {
		list = append(list, b)
	}
	return list
}

func (s *MemoryStore) UpdateCircuitBreaker(b *model.CircuitBreaker) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.circuitBreakers[b.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.circuitBreakers {
		if exist.ID != b.ID && exist.ServiceID == b.ServiceID && exist.RuleID == b.RuleID {
			return ErrConflict
		}
	}
	s.circuitBreakers[b.ID] = b
	return nil
}

func (s *MemoryStore) DeleteCircuitBreaker(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.circuitBreakers[id]; !ok {
		return ErrNotFound
	}
	delete(s.circuitBreakers, id)
	return nil
}

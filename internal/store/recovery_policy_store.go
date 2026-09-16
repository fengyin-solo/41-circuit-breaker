package store

import (
	"circuitbreaker/internal/model"
)

func (s *MemoryStore) CreateRecoveryPolicy(p *model.RecoveryPolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.recoveryPolicies {
		if exist.Name == p.Name && exist.ServiceID == p.ServiceID {
			return ErrConflict
		}
	}
	s.recoveryPolicies[p.ID] = p
	return nil
}

func (s *MemoryStore) GetRecoveryPolicy(id string) (*model.RecoveryPolicy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.recoveryPolicies[id]
	if !ok {
		return nil, ErrNotFound
	}
	return p, nil
}

func (s *MemoryStore) GetRecoveryPolicyByServiceID(serviceID string) (*model.RecoveryPolicy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, p := range s.recoveryPolicies {
		if p.ServiceID == serviceID {
			return p, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListRecoveryPolicies() []*model.RecoveryPolicy {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.RecoveryPolicy, 0, len(s.recoveryPolicies))
	for _, p := range s.recoveryPolicies {
		list = append(list, p)
	}
	return list
}

func (s *MemoryStore) UpdateRecoveryPolicy(p *model.RecoveryPolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.recoveryPolicies[p.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.recoveryPolicies {
		if exist.ID != p.ID && exist.Name == p.Name && exist.ServiceID == p.ServiceID {
			return ErrConflict
		}
	}
	s.recoveryPolicies[p.ID] = p
	return nil
}

func (s *MemoryStore) DeleteRecoveryPolicy(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.recoveryPolicies[id]; !ok {
		return ErrNotFound
	}
	delete(s.recoveryPolicies, id)
	return nil
}

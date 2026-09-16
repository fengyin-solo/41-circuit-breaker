package store

import (
	"circuitbreaker/internal/model"
)

func (s *MemoryStore) CreateAlertRule(a *model.AlertRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.alertRules {
		if exist.Name == a.Name && exist.ServiceID == a.ServiceID {
			return ErrConflict
		}
	}
	s.alertRules[a.ID] = a
	return nil
}

func (s *MemoryStore) GetAlertRule(id string) (*model.AlertRule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.alertRules[id]
	if !ok {
		return nil, ErrNotFound
	}
	return a, nil
}

func (s *MemoryStore) ListAlertRules() []*model.AlertRule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.AlertRule, 0, len(s.alertRules))
	for _, a := range s.alertRules {
		list = append(list, a)
	}
	return list
}

func (s *MemoryStore) UpdateAlertRule(a *model.AlertRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.alertRules[a.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.alertRules {
		if exist.ID != a.ID && exist.Name == a.Name && exist.ServiceID == a.ServiceID {
			return ErrConflict
		}
	}
	s.alertRules[a.ID] = a
	return nil
}

func (s *MemoryStore) DeleteAlertRule(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.alertRules[id]; !ok {
		return ErrNotFound
	}
	delete(s.alertRules, id)
	return nil
}

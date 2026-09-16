package store

import (
	"circuitbreaker/internal/model"
)

func (s *MemoryStore) CreateBreakerEvent(e *model.BreakerEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.breakerEvents[e.ID] = e
	return nil
}

func (s *MemoryStore) GetBreakerEvent(id string) (*model.BreakerEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.breakerEvents[id]
	if !ok {
		return nil, ErrNotFound
	}
	return e, nil
}

func (s *MemoryStore) ListBreakerEvents() []*model.BreakerEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.BreakerEvent, 0, len(s.breakerEvents))
	for _, e := range s.breakerEvents {
		list = append(list, e)
	}
	return list
}

func (s *MemoryStore) DeleteBreakerEvent(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.breakerEvents[id]; !ok {
		return ErrNotFound
	}
	delete(s.breakerEvents, id)
	return nil
}

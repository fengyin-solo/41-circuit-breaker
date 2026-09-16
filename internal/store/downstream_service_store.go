package store

import (
	"circuitbreaker/internal/model"
)

func (s *MemoryStore) CreateDownstreamService(svc *model.DownstreamService) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.downstreamServices {
		if exist.Name == svc.Name {
			return ErrConflict
		}
	}
	s.downstreamServices[svc.ID] = svc
	return nil
}

func (s *MemoryStore) GetDownstreamService(id string) (*model.DownstreamService, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	svc, ok := s.downstreamServices[id]
	if !ok {
		return nil, ErrNotFound
	}
	return svc, nil
}

func (s *MemoryStore) GetDownstreamServiceByName(name string) (*model.DownstreamService, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, svc := range s.downstreamServices {
		if svc.Name == name {
			return svc, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListDownstreamServices() []*model.DownstreamService {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.DownstreamService, 0, len(s.downstreamServices))
	for _, svc := range s.downstreamServices {
		list = append(list, svc)
	}
	return list
}

func (s *MemoryStore) UpdateDownstreamService(svc *model.DownstreamService) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.downstreamServices[svc.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.downstreamServices {
		if exist.ID != svc.ID && exist.Name == svc.Name {
			return ErrConflict
		}
	}
	s.downstreamServices[svc.ID] = svc
	return nil
}

func (s *MemoryStore) DeleteDownstreamService(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.downstreamServices[id]; !ok {
		return ErrNotFound
	}
	delete(s.downstreamServices, id)
	return nil
}

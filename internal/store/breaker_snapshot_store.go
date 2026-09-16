package store

import (
	"circuitbreaker/internal/model"
)

func (s *MemoryStore) CreateBreakerSnapshot(snap *model.BreakerSnapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.breakerSnapshots {
		if exist.ServiceID == snap.ServiceID && exist.SnapshotVersion == snap.SnapshotVersion {
			return ErrConflict
		}
	}
	s.breakerSnapshots[snap.ID] = snap
	return nil
}

func (s *MemoryStore) GetBreakerSnapshot(id string) (*model.BreakerSnapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	snap, ok := s.breakerSnapshots[id]
	if !ok {
		return nil, ErrNotFound
	}
	return snap, nil
}

func (s *MemoryStore) ListBreakerSnapshots() []*model.BreakerSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.BreakerSnapshot, 0, len(s.breakerSnapshots))
	for _, snap := range s.breakerSnapshots {
		list = append(list, snap)
	}
	return list
}

func (s *MemoryStore) UpdateBreakerSnapshot(snap *model.BreakerSnapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.breakerSnapshots[snap.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.breakerSnapshots {
		if exist.ID != snap.ID && exist.ServiceID == snap.ServiceID && exist.SnapshotVersion == snap.SnapshotVersion {
			return ErrConflict
		}
	}
	s.breakerSnapshots[snap.ID] = snap
	return nil
}

func (s *MemoryStore) DeleteBreakerSnapshot(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.breakerSnapshots[id]; !ok {
		return ErrNotFound
	}
	delete(s.breakerSnapshots, id)
	return nil
}

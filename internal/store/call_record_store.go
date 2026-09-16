package store

import (
	"circuitbreaker/internal/model"
)

func (s *MemoryStore) CreateCallRecord(c *model.CallRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.callRecords {
		if exist.RequestID == c.RequestID {
			return ErrConflict
		}
	}
	s.callRecords[c.ID] = c
	return nil
}

func (s *MemoryStore) GetCallRecord(id string) (*model.CallRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.callRecords[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (s *MemoryStore) ListCallRecords() []*model.CallRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.CallRecord, 0, len(s.callRecords))
	for _, c := range s.callRecords {
		list = append(list, c)
	}
	return list
}

func (s *MemoryStore) DeleteCallRecord(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.callRecords[id]; !ok {
		return ErrNotFound
	}
	delete(s.callRecords, id)
	return nil
}

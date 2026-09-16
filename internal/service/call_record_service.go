package service

import (
	"sort"
	"time"

	"circuitbreaker/internal/model"
	"circuitbreaker/pkg/idgen"
)

func (s *Service) CreateCallRecord(input model.CallRecord) (*model.CallRecord, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetDownstreamService(input.ServiceID); err != nil {
		return nil, model.NewValidationError("service_id", "关联的下游服务不存在")
	}
	input.ID = idgen.Hex()
	if input.CalledAt.IsZero() {
		input.CalledAt = time.Now()
	}
	if err := s.store.CreateCallRecord(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetCallRecord(id string) (*model.CallRecord, error) {
	return s.store.GetCallRecord(id)
}

func (s *Service) ListCallRecords(filter model.CallRecordFilter, page, size int) ([]*model.CallRecord, int, error) {
	all := s.store.ListCallRecords()
	matched := make([]*model.CallRecord, 0, len(all))
	for _, c := range all {
		if filter.Match(c) {
			matched = append(matched, c)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CalledAt.After(matched[j].CalledAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.CallRecord{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) DeleteCallRecord(id string) error {
	return s.store.DeleteCallRecord(id)
}

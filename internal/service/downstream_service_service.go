package service

import (
	"sort"
	"time"

	"circuitbreaker/internal/model"
	"circuitbreaker/pkg/idgen"
)

func (s *Service) CreateDownstreamService(input model.DownstreamService) (*model.DownstreamService, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	input.ID = idgen.Hex()
	input.CreatedAt = time.Now()
	input.UpdatedAt = input.CreatedAt
	if err := s.store.CreateDownstreamService(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetDownstreamService(id string) (*model.DownstreamService, error) {
	return s.store.GetDownstreamService(id)
}

func (s *Service) ListDownstreamServices(filter model.DownstreamServiceFilter, page, size int) ([]*model.DownstreamService, int, error) {
	all := s.store.ListDownstreamServices()
	matched := make([]*model.DownstreamService, 0, len(all))
	for _, svc := range all {
		if filter.Match(svc) {
			matched = append(matched, svc)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.DownstreamService{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateDownstreamService(id string, input model.DownstreamService) (*model.DownstreamService, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	exist, err := s.store.GetDownstreamService(id)
	if err != nil {
		return nil, err
	}
	exist.Name = input.Name
	exist.Address = input.Address
	exist.Protocol = input.Protocol
	exist.TimeoutMs = input.TimeoutMs
	exist.Status = input.Status
	exist.Weight = input.Weight
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateDownstreamService(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

func (s *Service) DeleteDownstreamService(id string) error {
	return s.store.DeleteDownstreamService(id)
}

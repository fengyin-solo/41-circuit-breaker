package service

import (
	"sort"
	"time"

	"circuitbreaker/internal/model"
	"circuitbreaker/pkg/idgen"
)

func (s *Service) CreateAlertRule(input model.AlertRule) (*model.AlertRule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetDownstreamService(input.ServiceID); err != nil {
		return nil, model.NewValidationError("service_id", "关联的下游服务不存在")
	}
	input.ID = idgen.Hex()
	input.CreatedAt = time.Now()
	input.UpdatedAt = input.CreatedAt
	if err := s.store.CreateAlertRule(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetAlertRule(id string) (*model.AlertRule, error) {
	return s.store.GetAlertRule(id)
}

func (s *Service) ListAlertRules(filter model.AlertRuleFilter, page, size int) ([]*model.AlertRule, int, error) {
	all := s.store.ListAlertRules()
	matched := make([]*model.AlertRule, 0, len(all))
	for _, a := range all {
		if filter.Match(a) {
			matched = append(matched, a)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.AlertRule{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateAlertRule(id string, input model.AlertRule) (*model.AlertRule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	exist, err := s.store.GetAlertRule(id)
	if err != nil {
		return nil, err
	}
	if input.ServiceID != exist.ServiceID {
		if _, err := s.store.GetDownstreamService(input.ServiceID); err != nil {
			return nil, model.NewValidationError("service_id", "关联的下游服务不存在")
		}
	}
	exist.Name = input.Name
	exist.ServiceID = input.ServiceID
	exist.Metric = input.Metric
	exist.Threshold = input.Threshold
	exist.Severity = input.Severity
	exist.NotifyChannel = input.NotifyChannel
	exist.Enabled = input.Enabled
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateAlertRule(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

func (s *Service) DeleteAlertRule(id string) error {
	return s.store.DeleteAlertRule(id)
}

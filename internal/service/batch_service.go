package service

import (
	"time"

	"circuitbreaker/internal/model"
	"circuitbreaker/pkg/idgen"
)

type BatchCreateServicesResult struct {
	Created []string `json:"created"`
	Errors  []string `json:"errors"`
}

func (s *Service) BatchCreateServices(inputs []model.DownstreamService) *BatchCreateServicesResult {
	res := &BatchCreateServicesResult{Created: []string{}, Errors: []string{}}
	for _, input := range inputs {
		if err := input.Validate(); err != nil {
			res.Errors = append(res.Errors, input.Name+": "+err.Error())
			continue
		}
		input.ID = idgen.Hex()
		input.CreatedAt = time.Now()
		input.UpdatedAt = input.CreatedAt
		if err := s.store.CreateDownstreamService(&input); err != nil {
			res.Errors = append(res.Errors, input.Name+": "+err.Error())
			continue
		}
		res.Created = append(res.Created, input.ID)
	}
	return res
}

func (s *Service) BatchDeleteServices(ids []string) []string {
	errors := []string{}
	for _, id := range ids {
		if err := s.store.DeleteDownstreamService(id); err != nil {
			errors = append(errors, id+": "+err.Error())
		}
	}
	return errors
}

type BatchTransitionBreakersResult struct {
	Success []string `json:"success"`
	Errors  []string `json:"errors"`
}

func (s *Service) BatchTransitionBreakers(ids []string, targetState string) *BatchTransitionBreakersResult {
	res := &BatchTransitionBreakersResult{Success: []string{}, Errors: []string{}}
	for _, id := range ids {
		exist, err := s.store.GetCircuitBreaker(id)
		if err != nil {
			res.Errors = append(res.Errors, id+": "+err.Error())
			continue
		}
		if !model.CanTransitionBreaker(exist.State, targetState) {
			res.Errors = append(res.Errors, id+": 非法状态流转 "+exist.State+" -> "+targetState)
			continue
		}
		_, err = s.UpdateCircuitBreaker(id, model.CircuitBreaker{
			ServiceID: exist.ServiceID,
			RuleID:    exist.RuleID,
			State:     targetState,
		})
		if err != nil {
			res.Errors = append(res.Errors, id+": "+err.Error())
			continue
		}
		res.Success = append(res.Success, id)
	}
	return res
}

type BatchCreateCallRecordsResult struct {
	Created []string `json:"created"`
	Errors  []string `json:"errors"`
}

func (s *Service) BatchCreateCallRecords(inputs []model.CallRecord) *BatchCreateCallRecordsResult {
	res := &BatchCreateCallRecordsResult{Created: []string{}, Errors: []string{}}
	for _, input := range inputs {
		if err := input.Validate(); err != nil {
			res.Errors = append(res.Errors, input.RequestID+": "+err.Error())
			continue
		}
		if _, err := s.store.GetDownstreamService(input.ServiceID); err != nil {
			res.Errors = append(res.Errors, input.RequestID+": 服务不存在")
			continue
		}
		input.ID = idgen.Hex()
		if input.CalledAt.IsZero() {
			input.CalledAt = time.Now()
		}
		if err := s.store.CreateCallRecord(&input); err != nil {
			res.Errors = append(res.Errors, input.RequestID+": "+err.Error())
			continue
		}
		res.Created = append(res.Created, input.ID)
	}
	return res
}

type BatchCreateSnapshotsResult struct {
	Created []string `json:"created"`
	Errors  []string `json:"errors"`
}

func (s *Service) BatchCreateSnapshots(inputs []model.BreakerSnapshot) *BatchCreateSnapshotsResult {
	res := &BatchCreateSnapshotsResult{Created: []string{}, Errors: []string{}}
	for _, input := range inputs {
		if err := input.Validate(); err != nil {
			res.Errors = append(res.Errors, input.ServiceID+": "+err.Error())
			continue
		}
		if _, err := s.store.GetDownstreamService(input.ServiceID); err != nil {
			res.Errors = append(res.Errors, input.ServiceID+": 服务不存在")
			continue
		}
		input.ID = idgen.Hex()
		input.CreatedAt = time.Now()
		if err := s.store.CreateBreakerSnapshot(&input); err != nil {
			res.Errors = append(res.Errors, input.ServiceID+": "+err.Error())
			continue
		}
		res.Created = append(res.Created, input.ID)
	}
	return res
}

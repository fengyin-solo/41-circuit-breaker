package service

import (
	"sort"

	"circuitbreaker/internal/model"
)

type StatsOverview struct {
	TotalServices      int     `json:"total_services"`
	TotalBreakers      int     `json:"total_breakers"`
	OpenBreakers       int     `json:"open_breakers"`
	HalfOpenBreakers   int     `json:"half_open_breakers"`
	ClosedBreakers     int     `json:"closed_breakers"`
	TotalCalls         int     `json:"total_calls"`
	TotalFailures      int     `json:"total_failures"`
	AvgFailureRatio    float64 `json:"avg_failure_ratio"`
	TotalAlertRules    int     `json:"total_alert_rules"`
	EnabledAlertRules  int     `json:"enabled_alert_rules"`
}

func (s *Service) GetStatsOverview() *StatsOverview {
	services := s.store.ListDownstreamServices()
	breakers := s.store.ListCircuitBreakers()
	records := s.store.ListCallRecords()
	rules := s.store.ListAlertRules()

	openCount := 0
	halfOpenCount := 0
	closedCount := 0
	for _, b := range breakers {
		switch b.State {
		case model.BreakerStateOpen:
			openCount++
		case model.BreakerStateHalfOpen:
			halfOpenCount++
		case model.BreakerStateClosed:
			closedCount++
		}
	}

	totalCalls := 0
	totalFailures := 0
	for _, r := range records {
		totalCalls++
		if r.Outcome == model.CallOutcomeFailure || r.Outcome == model.CallOutcomeTimeout {
			totalFailures++
		}
	}

	avgFailure := 0.0
	if totalCalls > 0 {
		avgFailure = float64(totalFailures) / float64(totalCalls)
	}

	enabledRules := 0
	for _, r := range rules {
		if r.Enabled {
			enabledRules++
		}
	}

	return &StatsOverview{
		TotalServices:     len(services),
		TotalBreakers:     len(breakers),
		OpenBreakers:      openCount,
		HalfOpenBreakers:  halfOpenCount,
		ClosedBreakers:    closedCount,
		TotalCalls:        totalCalls,
		TotalFailures:     totalFailures,
		AvgFailureRatio:   avgFailure,
		TotalAlertRules:   len(rules),
		EnabledAlertRules: enabledRules,
	}
}

type ServiceStats struct {
	ServiceID       string  `json:"service_id"`
	ServiceName     string  `json:"service_name"`
	TotalCalls      int     `json:"total_calls"`
	SuccessCalls    int     `json:"success_calls"`
	FailureCalls    int     `json:"failure_calls"`
	TimeoutCalls    int     `json:"timeout_calls"`
	FailureRatio    float64 `json:"failure_ratio"`
	BreakerState    string  `json:"breaker_state"`
}

func (s *Service) GetServiceStats() []*ServiceStats {
	services := s.store.ListDownstreamServices()
	records := s.store.ListCallRecords()
	breakers := s.store.ListCircuitBreakers()

	result := make([]*ServiceStats, 0, len(services))
	for _, svc := range services {
		st := &ServiceStats{
			ServiceID:   svc.ID,
			ServiceName: svc.Name,
		}
		for _, r := range records {
			if r.ServiceID == svc.ID {
				st.TotalCalls++
				switch r.Outcome {
				case model.CallOutcomeSuccess:
					st.SuccessCalls++
				case model.CallOutcomeFailure:
					st.FailureCalls++
				case model.CallOutcomeTimeout:
					st.TimeoutCalls++
				}
			}
		}
		if st.TotalCalls > 0 {
			st.FailureRatio = float64(st.FailureCalls+st.TimeoutCalls) / float64(st.TotalCalls)
		}
		for _, b := range breakers {
			if b.ServiceID == svc.ID {
				st.BreakerState = b.State
				break
			}
		}
		result = append(result, st)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].TotalCalls > result[j].TotalCalls
	})
	return result
}

type BreakerStateGroup struct {
	State string `json:"state"`
	Count int    `json:"count"`
}

func (s *Service) GetBreakerStateGroups() []*BreakerStateGroup {
	breakers := s.store.ListCircuitBreakers()
	counts := make(map[string]int)
	for _, b := range breakers {
		counts[b.State]++
	}
	result := make([]*BreakerStateGroup, 0, len(counts))
	for state, count := range counts {
		result = append(result, &BreakerStateGroup{State: state, Count: count})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})
	return result
}

type TopFailureService struct {
	ServiceID    string  `json:"service_id"`
	ServiceName  string  `json:"service_name"`
	FailureCount int     `json:"failure_count"`
	FailureRatio float64 `json:"failure_ratio"`
}

func (s *Service) GetTopFailureServices(n int) []*TopFailureService {
	stats := s.GetServiceStats()
	filtered := make([]*TopFailureService, 0, len(stats))
	nameMap := make(map[string]string)
	services := s.store.ListDownstreamServices()
	for _, svc := range services {
		nameMap[svc.ID] = svc.Name
	}
	for _, st := range stats {
		if st.FailureCalls > 0 {
			filtered = append(filtered, &TopFailureService{
				ServiceID:    st.ServiceID,
				ServiceName:  nameMap[st.ServiceID],
				FailureCount: st.FailureCalls + st.TimeoutCalls,
				FailureRatio: st.FailureRatio,
			})
		}
	}
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].FailureCount > filtered[j].FailureCount
	})
	if n > len(filtered) {
		n = len(filtered)
	}
	return filtered[:n]
}

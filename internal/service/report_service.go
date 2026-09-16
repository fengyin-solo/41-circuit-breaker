package service

import (
	"sort"
	"time"

	"circuitbreaker/internal/model"
)

type TimeWindowReport struct {
	WindowStart  string  `json:"window_start"`
	WindowEnd    string  `json:"window_end"`
	TotalCalls   int     `json:"total_calls"`
	SuccessCalls int     `json:"success_calls"`
	FailureCalls int     `json:"failure_calls"`
	TimeoutCalls int     `json:"timeout_calls"`
	FailureRatio float64 `json:"failure_ratio"`
}

func (s *Service) GetCallRecordsTimeWindowReport(serviceID string, windowMinutes int) []*TimeWindowReport {
	if windowMinutes <= 0 {
		windowMinutes = 60
	}
	records := s.store.ListCallRecords()
	windows := make(map[string]*TimeWindowReport)
	for _, r := range records {
		if serviceID != "" && r.ServiceID != serviceID {
			continue
		}
		bucket := r.CalledAt.Truncate(time.Duration(windowMinutes) * time.Minute)
		key := bucket.Format(time.RFC3339)
		if _, ok := windows[key]; !ok {
			windows[key] = &TimeWindowReport{
				WindowStart: bucket.Format(time.RFC3339),
				WindowEnd:   bucket.Add(time.Duration(windowMinutes) * time.Minute).Format(time.RFC3339),
			}
		}
		w := windows[key]
		w.TotalCalls++
		switch r.Outcome {
		case model.CallOutcomeSuccess:
			w.SuccessCalls++
		case model.CallOutcomeFailure:
			w.FailureCalls++
		case model.CallOutcomeTimeout:
			w.TimeoutCalls++
		}
	}
	result := make([]*TimeWindowReport, 0, len(windows))
	for _, w := range windows {
		if w.TotalCalls > 0 {
			w.FailureRatio = float64(w.FailureCalls+w.TimeoutCalls) / float64(w.TotalCalls)
		}
		result = append(result, w)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].WindowStart > result[j].WindowStart
	})
	if len(result) > 24 {
		result = result[:24]
	}
	return result
}

type RuleGroupStats struct {
	RuleID       string  `json:"rule_id"`
	RuleName     string  `json:"rule_name"`
	BreakerCount int     `json:"breaker_count"`
	OpenCount    int     `json:"open_count"`
	ClosedCount  int     `json:"closed_count"`
	HalfOpenCount int    `json:"half_open_count"`
}

func (s *Service) GetBreakerRuleGroupStats() []*RuleGroupStats {
	rules := s.store.ListBreakerRules()
	breakers := s.store.ListCircuitBreakers()
	result := make([]*RuleGroupStats, 0, len(rules))
	for _, rule := range rules {
		st := &RuleGroupStats{RuleID: rule.ID, RuleName: rule.Name}
		for _, b := range breakers {
			if b.RuleID == rule.ID {
				st.BreakerCount++
				switch b.State {
				case model.BreakerStateOpen:
					st.OpenCount++
				case model.BreakerStateClosed:
					st.ClosedCount++
				case model.BreakerStateHalfOpen:
					st.HalfOpenCount++
				}
			}
		}
		result = append(result, st)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].BreakerCount > result[j].BreakerCount
	})
	return result
}

type EventTrendItem struct {
	Time       string `json:"time"`
	Opened     int    `json:"opened"`
	Closed     int    `json:"closed"`
	HalfOpened int    `json:"half_opened"`
	Rejected   int    `json:"rejected"`
}

func (s *Service) GetBreakerEventTrend(windowMinutes int) []*EventTrendItem {
	if windowMinutes <= 0 {
		windowMinutes = 60
	}
	events := s.store.ListBreakerEvents()
	windows := make(map[string]*EventTrendItem)
	for _, e := range events {
		bucket := e.OccurredAt.Truncate(time.Duration(windowMinutes) * time.Minute)
		key := bucket.Format(time.RFC3339)
		if _, ok := windows[key]; !ok {
			windows[key] = &EventTrendItem{Time: bucket.Format(time.RFC3339)}
		}
		w := windows[key]
		switch e.EventType {
		case model.EventTypeOpened:
			w.Opened++
		case model.EventTypeClosed:
			w.Closed++
		case model.EventTypeHalfOpen:
			w.HalfOpened++
		case model.EventTypeRejected:
			w.Rejected++
		}
	}
	result := make([]*EventTrendItem, 0, len(windows))
	for _, w := range windows {
		result = append(result, w)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Time > result[j].Time
	})
	if len(result) > 24 {
		result = result[:24]
	}
	return result
}

type ServiceHealthSummary struct {
	ServiceID           string `json:"service_id"`
	ServiceName         string `json:"service_name"`
	LastStatus          string `json:"last_status"`
	ConsecutiveFailures int    `json:"consecutive_failures"`
	IntervalSeconds     int    `json:"interval_seconds"`
}

func (s *Service) GetServiceHealthSummaries() []*ServiceHealthSummary {
	services := s.store.ListDownstreamServices()
	checks := s.store.ListHealthChecks()
	checkMap := make(map[string]*model.HealthCheck)
	for _, c := range checks {
		checkMap[c.ServiceID] = c
	}
	result := make([]*ServiceHealthSummary, 0, len(services))
	for _, svc := range services {
		summary := &ServiceHealthSummary{
			ServiceID:   svc.ID,
			ServiceName: svc.Name,
		}
		if c, ok := checkMap[svc.ID]; ok {
			summary.LastStatus = c.LastStatus
			summary.ConsecutiveFailures = c.ConsecutiveFailures
			summary.IntervalSeconds = c.IntervalSeconds
		} else {
			summary.LastStatus = "unknown"
		}
		result = append(result, summary)
	}
	return result
}

type LatencyPercentileStats struct {
	ServiceID   string `json:"service_id"`
	ServiceName string `json:"service_name"`
	P50         int    `json:"p50"`
	P90         int    `json:"p90"`
	P99         int    `json:"p99"`
	Max         int    `json:"max"`
	Min         int    `json:"min"`
	Avg         int    `json:"avg"`
}

func (s *Service) GetLatencyPercentiles() []*LatencyPercentileStats {
	services := s.store.ListDownstreamServices()
	records := s.store.ListCallRecords()
	result := make([]*LatencyPercentileStats, 0, len(services))
	for _, svc := range services {
		latencies := []int{}
		sum := 0
		for _, r := range records {
			if r.ServiceID == svc.ID {
				latencies = append(latencies, r.LatencyMs)
				sum += r.LatencyMs
			}
		}
		if len(latencies) == 0 {
			continue
		}
		sort.Ints(latencies)
		st := &LatencyPercentileStats{
			ServiceID:   svc.ID,
			ServiceName: svc.Name,
			Min:         latencies[0],
			Max:         latencies[len(latencies)-1],
			Avg:         sum / len(latencies),
		}
		st.P50 = percentileInt(latencies, 0.5)
		st.P90 = percentileInt(latencies, 0.9)
		st.P99 = percentileInt(latencies, 0.99)
		result = append(result, st)
	}
	return result
}

func percentileInt(sorted []int, p float64) int {
	if len(sorted) == 0 {
		return 0
	}
	idx := int(float64(len(sorted)-1) * p)
	if idx < 0 {
		idx = 0
	}
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}

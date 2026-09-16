package handler

import (
	"net/http"
	"strconv"

	"circuitbreaker/pkg/httpx"
)

func (s *Server) registerReportRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/reports/call-records", s.getCallRecordsTimeWindowReport)
	mux.HandleFunc("GET /api/reports/breaker-rules", s.getBreakerRuleGroupStats)
	mux.HandleFunc("GET /api/reports/event-trend", s.getBreakerEventTrend)
	mux.HandleFunc("GET /api/reports/health-summary", s.getServiceHealthSummaries)
	mux.HandleFunc("GET /api/reports/latency-percentiles", s.getLatencyPercentiles)
}

func (s *Server) getCallRecordsTimeWindowReport(w http.ResponseWriter, r *http.Request) {
	serviceID := r.URL.Query().Get("service_id")
	windowStr := r.URL.Query().Get("window_minutes")
	window, _ := strconv.Atoi(windowStr)
	report := s.svc.GetCallRecordsTimeWindowReport(serviceID, window)
	httpx.OK(w, report)
}

func (s *Server) getBreakerRuleGroupStats(w http.ResponseWriter, r *http.Request) {
	stats := s.svc.GetBreakerRuleGroupStats()
	httpx.OK(w, stats)
}

func (s *Server) getBreakerEventTrend(w http.ResponseWriter, r *http.Request) {
	windowStr := r.URL.Query().Get("window_minutes")
	window, _ := strconv.Atoi(windowStr)
	trend := s.svc.GetBreakerEventTrend(window)
	httpx.OK(w, trend)
}

func (s *Server) getServiceHealthSummaries(w http.ResponseWriter, r *http.Request) {
	summaries := s.svc.GetServiceHealthSummaries()
	httpx.OK(w, summaries)
}

func (s *Server) getLatencyPercentiles(w http.ResponseWriter, r *http.Request) {
	stats := s.svc.GetLatencyPercentiles()
	httpx.OK(w, stats)
}

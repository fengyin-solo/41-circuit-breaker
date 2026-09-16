package handler

import (
	"net/http"
	"strconv"

	"circuitbreaker/pkg/httpx"
)

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/overview", s.getStatsOverview)
	mux.HandleFunc("GET /api/stats/services", s.getServiceStats)
	mux.HandleFunc("GET /api/stats/breaker-states", s.getBreakerStateGroups)
	mux.HandleFunc("GET /api/stats/top-failures", s.getTopFailureServices)
}

func (s *Server) getStatsOverview(w http.ResponseWriter, r *http.Request) {
	stats := s.svc.GetStatsOverview()
	httpx.OK(w, stats)
}

func (s *Server) getServiceStats(w http.ResponseWriter, r *http.Request) {
	stats := s.svc.GetServiceStats()
	httpx.OK(w, stats)
}

func (s *Server) getBreakerStateGroups(w http.ResponseWriter, r *http.Request) {
	groups := s.svc.GetBreakerStateGroups()
	httpx.OK(w, groups)
}

func (s *Server) getTopFailureServices(w http.ResponseWriter, r *http.Request) {
	nStr := r.URL.Query().Get("n")
	n, _ := strconv.Atoi(nStr)
	if n <= 0 {
		n = 5
	}
	if n > 20 {
		n = 20
	}
	top := s.svc.GetTopFailureServices(n)
	httpx.OK(w, top)
}

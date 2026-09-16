package handler

import (
	"net/http"

	"circuitbreaker/internal/model"
	"circuitbreaker/pkg/httpx"
)

func (s *Server) registerHealthCheckRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/health-checks", s.createHealthCheck)
	mux.HandleFunc("GET /api/health-checks", s.listHealthChecks)
	mux.HandleFunc("GET /api/health-checks/{id}", s.getHealthCheck)
	mux.HandleFunc("PUT /api/health-checks/{id}", s.updateHealthCheck)
	mux.HandleFunc("DELETE /api/health-checks/{id}", s.deleteHealthCheck)
}

type createHealthCheckRequest struct {
	ServiceID           string `json:"service_id"`
	IntervalSeconds     int    `json:"interval_seconds"`
	LastStatus          string `json:"last_status"`
	ConsecutiveFailures int    `json:"consecutive_failures"`
}

func (s *Server) createHealthCheck(w http.ResponseWriter, r *http.Request) {
	var req createHealthCheckRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	hc, err := s.svc.CreateHealthCheck(model.HealthCheck{
		ServiceID:           req.ServiceID,
		IntervalSeconds:     req.IntervalSeconds,
		LastStatus:          req.LastStatus,
		ConsecutiveFailures: req.ConsecutiveFailures,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, hc)
}

func (s *Server) listHealthChecks(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.HealthCheckFilter{
		ServiceID:  r.URL.Query().Get("service_id"),
		LastStatus: r.URL.Query().Get("last_status"),
		Keyword:    r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListHealthChecks(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getHealthCheck(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	hc, err := s.svc.GetHealthCheck(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, hc)
}

type updateHealthCheckRequest struct {
	ServiceID           string `json:"service_id"`
	IntervalSeconds     int    `json:"interval_seconds"`
	LastStatus          string `json:"last_status"`
	ConsecutiveFailures int    `json:"consecutive_failures"`
}

func (s *Server) updateHealthCheck(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateHealthCheckRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	hc, err := s.svc.UpdateHealthCheck(id, model.HealthCheck{
		ServiceID:           req.ServiceID,
		IntervalSeconds:     req.IntervalSeconds,
		LastStatus:          req.LastStatus,
		ConsecutiveFailures: req.ConsecutiveFailures,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, hc)
}

func (s *Server) deleteHealthCheck(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteHealthCheck(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

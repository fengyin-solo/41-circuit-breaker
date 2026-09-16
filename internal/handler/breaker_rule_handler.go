package handler

import (
	"net/http"
	"strconv"

	"circuitbreaker/internal/model"
	"circuitbreaker/pkg/httpx"
)

func (s *Server) registerBreakerRuleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/rules", s.createBreakerRule)
	mux.HandleFunc("GET /api/rules", s.listBreakerRules)
	mux.HandleFunc("GET /api/rules/{id}", s.getBreakerRule)
	mux.HandleFunc("PUT /api/rules/{id}", s.updateBreakerRule)
	mux.HandleFunc("DELETE /api/rules/{id}", s.deleteBreakerRule)
}

type createBreakerRuleRequest struct {
	Name                   string  `json:"name"`
	ServiceID              string  `json:"service_id"`
	FailureRatioThreshold  float64 `json:"failure_ratio_threshold"`
	SlowCallRatioThreshold float64 `json:"slow_call_ratio_threshold"`
	SlowCallMs             int     `json:"slow_call_ms"`
	WindowSeconds          int     `json:"window_seconds"`
	MinRequestCount        int     `json:"min_request_count"`
	MaxHalfOpenRequests    int     `json:"max_half_open_requests"`
	Enabled                bool    `json:"enabled"`
}

func (s *Server) createBreakerRule(w http.ResponseWriter, r *http.Request) {
	var req createBreakerRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rule, err := s.svc.CreateBreakerRule(model.BreakerRule{
		Name:                   req.Name,
		ServiceID:              req.ServiceID,
		FailureRatioThreshold:  req.FailureRatioThreshold,
		SlowCallRatioThreshold: req.SlowCallRatioThreshold,
		SlowCallMs:             req.SlowCallMs,
		WindowSeconds:          req.WindowSeconds,
		MinRequestCount:        req.MinRequestCount,
		MaxHalfOpenRequests:    req.MaxHalfOpenRequests,
		Enabled:                req.Enabled,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, rule)
}

func (s *Server) listBreakerRules(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.BreakerRuleFilter{
		Name:      r.URL.Query().Get("name"),
		ServiceID: r.URL.Query().Get("service_id"),
		Keyword:   r.URL.Query().Get("keyword"),
	}
	if enabledStr := r.URL.Query().Get("enabled"); enabledStr != "" {
		v, _ := strconv.ParseBool(enabledStr)
		filter.Enabled = &v
	}
	items, total, err := s.svc.ListBreakerRules(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getBreakerRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rule, err := s.svc.GetBreakerRule(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rule)
}

type updateBreakerRuleRequest struct {
	Name                   string  `json:"name"`
	ServiceID              string  `json:"service_id"`
	FailureRatioThreshold  float64 `json:"failure_ratio_threshold"`
	SlowCallRatioThreshold float64 `json:"slow_call_ratio_threshold"`
	SlowCallMs             int     `json:"slow_call_ms"`
	WindowSeconds          int     `json:"window_seconds"`
	MinRequestCount        int     `json:"min_request_count"`
	MaxHalfOpenRequests    int     `json:"max_half_open_requests"`
	Enabled                bool    `json:"enabled"`
}

func (s *Server) updateBreakerRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateBreakerRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rule, err := s.svc.UpdateBreakerRule(id, model.BreakerRule{
		Name:                   req.Name,
		ServiceID:              req.ServiceID,
		FailureRatioThreshold:  req.FailureRatioThreshold,
		SlowCallRatioThreshold: req.SlowCallRatioThreshold,
		SlowCallMs:             req.SlowCallMs,
		WindowSeconds:          req.WindowSeconds,
		MinRequestCount:        req.MinRequestCount,
		MaxHalfOpenRequests:    req.MaxHalfOpenRequests,
		Enabled:                req.Enabled,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rule)
}

func (s *Server) deleteBreakerRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteBreakerRule(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

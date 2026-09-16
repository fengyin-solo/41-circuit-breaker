package handler

import (
	"net/http"

	"circuitbreaker/internal/model"
	"circuitbreaker/pkg/httpx"
)

func (s *Server) registerRecoveryPolicyRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/policies", s.createRecoveryPolicy)
	mux.HandleFunc("GET /api/policies", s.listRecoveryPolicies)
	mux.HandleFunc("GET /api/policies/{id}", s.getRecoveryPolicy)
	mux.HandleFunc("PUT /api/policies/{id}", s.updateRecoveryPolicy)
	mux.HandleFunc("DELETE /api/policies/{id}", s.deleteRecoveryPolicy)
}

type createRecoveryPolicyRequest struct {
	Name                  string  `json:"name"`
	ServiceID             string  `json:"service_id"`
	HalfOpenProbeRatio    float64 `json:"half_open_probe_ratio"`
	RecoveryWindowSeconds int     `json:"recovery_window_seconds"`
	MaxRetry              int     `json:"max_retry"`
}

func (s *Server) createRecoveryPolicy(w http.ResponseWriter, r *http.Request) {
	var req createRecoveryPolicyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.CreateRecoveryPolicy(model.RecoveryPolicy{
		Name:                  req.Name,
		ServiceID:             req.ServiceID,
		HalfOpenProbeRatio:    req.HalfOpenProbeRatio,
		RecoveryWindowSeconds: req.RecoveryWindowSeconds,
		MaxRetry:              req.MaxRetry,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, p)
}

func (s *Server) listRecoveryPolicies(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.RecoveryPolicyFilter{
		Name:      r.URL.Query().Get("name"),
		ServiceID: r.URL.Query().Get("service_id"),
		Keyword:   r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListRecoveryPolicies(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRecoveryPolicy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	p, err := s.svc.GetRecoveryPolicy(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

type updateRecoveryPolicyRequest struct {
	Name                  string  `json:"name"`
	ServiceID             string  `json:"service_id"`
	HalfOpenProbeRatio    float64 `json:"half_open_probe_ratio"`
	RecoveryWindowSeconds int     `json:"recovery_window_seconds"`
	MaxRetry              int     `json:"max_retry"`
}

func (s *Server) updateRecoveryPolicy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateRecoveryPolicyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.UpdateRecoveryPolicy(id, model.RecoveryPolicy{
		Name:                  req.Name,
		ServiceID:             req.ServiceID,
		HalfOpenProbeRatio:    req.HalfOpenProbeRatio,
		RecoveryWindowSeconds: req.RecoveryWindowSeconds,
		MaxRetry:              req.MaxRetry,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

func (s *Server) deleteRecoveryPolicy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteRecoveryPolicy(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

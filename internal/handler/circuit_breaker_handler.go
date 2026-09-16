package handler

import (
	"net/http"

	"circuitbreaker/internal/model"
	"circuitbreaker/pkg/httpx"
)

func (s *Server) registerCircuitBreakerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/breakers", s.createCircuitBreaker)
	mux.HandleFunc("GET /api/breakers", s.listCircuitBreakers)
	mux.HandleFunc("GET /api/breakers/{id}", s.getCircuitBreaker)
	mux.HandleFunc("PUT /api/breakers/{id}", s.updateCircuitBreaker)
	mux.HandleFunc("DELETE /api/breakers/{id}", s.deleteCircuitBreaker)
	mux.HandleFunc("POST /api/breakers/{id}/transition", s.transitionCircuitBreaker)
}

type createCircuitBreakerRequest struct {
	ServiceID string `json:"service_id"`
	RuleID    string `json:"rule_id"`
	State     string `json:"state"`
}

func (s *Server) createCircuitBreaker(w http.ResponseWriter, r *http.Request) {
	var req createCircuitBreakerRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	b, err := s.svc.CreateCircuitBreaker(model.CircuitBreaker{
		ServiceID: req.ServiceID,
		RuleID:    req.RuleID,
		State:     req.State,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, b)
}

func (s *Server) listCircuitBreakers(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.CircuitBreakerFilter{
		ServiceID: r.URL.Query().Get("service_id"),
		RuleID:    r.URL.Query().Get("rule_id"),
		State:     r.URL.Query().Get("state"),
		Keyword:   r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListCircuitBreakers(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getCircuitBreaker(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	b, err := s.svc.GetCircuitBreaker(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, b)
}

type updateCircuitBreakerRequest struct {
	ServiceID string `json:"service_id"`
	RuleID    string `json:"rule_id"`
	State     string `json:"state"`
}

func (s *Server) updateCircuitBreaker(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateCircuitBreakerRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	b, err := s.svc.UpdateCircuitBreaker(id, model.CircuitBreaker{
		ServiceID: req.ServiceID,
		RuleID:    req.RuleID,
		State:     req.State,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, b)
}

func (s *Server) deleteCircuitBreaker(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteCircuitBreaker(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type transitionCircuitBreakerRequest struct {
	TargetState string `json:"target_state"`
}

func (s *Server) transitionCircuitBreaker(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req transitionCircuitBreakerRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	exist, err := s.svc.GetCircuitBreaker(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	b, err := s.svc.UpdateCircuitBreaker(id, model.CircuitBreaker{
		ServiceID: exist.ServiceID,
		RuleID:    exist.RuleID,
		State:     req.TargetState,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, b)
}

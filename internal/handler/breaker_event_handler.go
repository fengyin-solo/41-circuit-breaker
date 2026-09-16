package handler

import (
	"net/http"

	"circuitbreaker/internal/model"
	"circuitbreaker/pkg/httpx"
)

func (s *Server) registerBreakerEventRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/events", s.createBreakerEvent)
	mux.HandleFunc("GET /api/events", s.listBreakerEvents)
	mux.HandleFunc("GET /api/events/{id}", s.getBreakerEvent)
	mux.HandleFunc("DELETE /api/events/{id}", s.deleteBreakerEvent)
}

type createBreakerEventRequest struct {
	BreakerID string `json:"breaker_id"`
	ServiceID string `json:"service_id"`
	EventType string `json:"event_type"`
	Reason    string `json:"reason"`
}

func (s *Server) createBreakerEvent(w http.ResponseWriter, r *http.Request) {
	var req createBreakerEventRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	e, err := s.svc.CreateBreakerEvent(model.BreakerEvent{
		BreakerID: req.BreakerID,
		ServiceID: req.ServiceID,
		EventType: req.EventType,
		Reason:    req.Reason,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, e)
}

func (s *Server) listBreakerEvents(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.BreakerEventFilter{
		BreakerID: r.URL.Query().Get("breaker_id"),
		ServiceID: r.URL.Query().Get("service_id"),
		EventType: r.URL.Query().Get("event_type"),
		Keyword:   r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListBreakerEvents(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getBreakerEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	e, err := s.svc.GetBreakerEvent(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, e)
}

func (s *Server) deleteBreakerEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteBreakerEvent(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

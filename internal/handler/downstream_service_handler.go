package handler

import (
	"net/http"

	"circuitbreaker/internal/model"
	"circuitbreaker/pkg/httpx"
)

func (s *Server) registerDownstreamServiceRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/services", s.createDownstreamService)
	mux.HandleFunc("GET /api/services", s.listDownstreamServices)
	mux.HandleFunc("GET /api/services/{id}", s.getDownstreamService)
	mux.HandleFunc("PUT /api/services/{id}", s.updateDownstreamService)
	mux.HandleFunc("DELETE /api/services/{id}", s.deleteDownstreamService)
}

type createDownstreamServiceRequest struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	Protocol  string `json:"protocol"`
	TimeoutMs int    `json:"timeout_ms"`
	Status    string `json:"status"`
	Weight    int    `json:"weight"`
}

func (s *Server) createDownstreamService(w http.ResponseWriter, r *http.Request) {
	var req createDownstreamServiceRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	svc, err := s.svc.CreateDownstreamService(model.DownstreamService{
		Name:      req.Name,
		Address:   req.Address,
		Protocol:  req.Protocol,
		TimeoutMs: req.TimeoutMs,
		Status:    req.Status,
		Weight:    req.Weight,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, svc)
}

func (s *Server) listDownstreamServices(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.DownstreamServiceFilter{
		Name:     r.URL.Query().Get("name"),
		Protocol: r.URL.Query().Get("protocol"),
		Status:   r.URL.Query().Get("status"),
		Keyword:  r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListDownstreamServices(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getDownstreamService(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	svc, err := s.svc.GetDownstreamService(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, svc)
}

type updateDownstreamServiceRequest struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	Protocol  string `json:"protocol"`
	TimeoutMs int    `json:"timeout_ms"`
	Status    string `json:"status"`
	Weight    int    `json:"weight"`
}

func (s *Server) updateDownstreamService(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateDownstreamServiceRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	svc, err := s.svc.UpdateDownstreamService(id, model.DownstreamService{
		Name:      req.Name,
		Address:   req.Address,
		Protocol:  req.Protocol,
		TimeoutMs: req.TimeoutMs,
		Status:    req.Status,
		Weight:    req.Weight,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, svc)
}

func (s *Server) deleteDownstreamService(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteDownstreamService(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

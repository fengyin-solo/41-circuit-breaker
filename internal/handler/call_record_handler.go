package handler

import (
	"net/http"

	"circuitbreaker/internal/model"
	"circuitbreaker/pkg/httpx"
)

func (s *Server) registerCallRecordRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/records", s.createCallRecord)
	mux.HandleFunc("GET /api/records", s.listCallRecords)
	mux.HandleFunc("GET /api/records/{id}", s.getCallRecord)
	mux.HandleFunc("DELETE /api/records/{id}", s.deleteCallRecord)
}

type createCallRecordRequest struct {
	RequestID string `json:"request_id"`
	ServiceID string `json:"service_id"`
	Outcome   string `json:"outcome"`
	LatencyMs int    `json:"latency_ms"`
}

func (s *Server) createCallRecord(w http.ResponseWriter, r *http.Request) {
	var req createCallRecordRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rec, err := s.svc.CreateCallRecord(model.CallRecord{
		RequestID: req.RequestID,
		ServiceID: req.ServiceID,
		Outcome:   req.Outcome,
		LatencyMs: req.LatencyMs,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, rec)
}

func (s *Server) listCallRecords(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.CallRecordFilter{
		ServiceID: r.URL.Query().Get("service_id"),
		Outcome:   r.URL.Query().Get("outcome"),
		Keyword:   r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListCallRecords(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getCallRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rec, err := s.svc.GetCallRecord(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rec)
}

func (s *Server) deleteCallRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteCallRecord(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

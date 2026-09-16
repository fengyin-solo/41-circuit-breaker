package handler

import (
	"net/http"

	"circuitbreaker/internal/model"
	"circuitbreaker/pkg/httpx"
)

func (s *Server) registerMetricSampleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/metrics", s.createMetricSample)
	mux.HandleFunc("GET /api/metrics", s.listMetricSamples)
	mux.HandleFunc("GET /api/metrics/{id}", s.getMetricSample)
	mux.HandleFunc("DELETE /api/metrics/{id}", s.deleteMetricSample)
}

type createMetricSampleRequest struct {
	ServiceID    string  `json:"service_id"`
	WindowStart  string  `json:"window_start"`
	WindowEnd    string  `json:"window_end"`
	TotalCalls   int     `json:"total_calls"`
	SuccessCalls int     `json:"success_calls"`
	FailureCalls int     `json:"failure_calls"`
	SlowCalls    int     `json:"slow_calls"`
	FailureRatio float64 `json:"failure_ratio"`
	SlowRatio    float64 `json:"slow_ratio"`
}

func (s *Server) createMetricSample(w http.ResponseWriter, r *http.Request) {
	var req createMetricSampleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	// parse times omitted for brevity; keep simple
	m, err := s.svc.CreateMetricSample(model.MetricSample{
		ServiceID:    req.ServiceID,
		TotalCalls:   req.TotalCalls,
		SuccessCalls: req.SuccessCalls,
		FailureCalls: req.FailureCalls,
		SlowCalls:    req.SlowCalls,
		FailureRatio: req.FailureRatio,
		SlowRatio:    req.SlowRatio,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, m)
}

func (s *Server) listMetricSamples(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.MetricSampleFilter{
		ServiceID: r.URL.Query().Get("service_id"),
		Keyword:   r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListMetricSamples(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getMetricSample(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	m, err := s.svc.GetMetricSample(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, m)
}

func (s *Server) deleteMetricSample(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteMetricSample(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

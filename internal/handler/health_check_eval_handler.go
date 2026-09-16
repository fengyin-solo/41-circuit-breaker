package handler

import (
	"net/http"

	"circuitbreaker/pkg/httpx"
)

func (s *Server) registerHealthCheckEvalRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/health-checks/evaluate", s.evaluateAllHealthChecks)
	mux.HandleFunc("POST /api/health-checks/{id}/evaluate", s.evaluateSingleHealthCheck)
	mux.HandleFunc("POST /api/health-checks/{id}/manual-check", s.performManualHealthCheck)
}

func (s *Server) evaluateAllHealthChecks(w http.ResponseWriter, r *http.Request) {
	results := s.svc.EvaluateHealthChecks()
	httpx.OK(w, results)
}

func (s *Server) evaluateSingleHealthCheck(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	hc, err := s.svc.GetHealthCheck(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	result, err := s.svc.EvaluateSingleHealthCheck(hc.ServiceID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

type manualHealthCheckRequest struct {
	Success bool `json:"success"`
}

func (s *Server) performManualHealthCheck(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	hc, err := s.svc.GetHealthCheck(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	var req manualHealthCheckRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.PerformManualHealthCheck(hc.ServiceID, req.Success)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

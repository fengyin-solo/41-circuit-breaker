package handler

import (
	"net/http"

	"circuitbreaker/pkg/httpx"
)

func (s *Server) registerBreakerSimulationRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/simulate/call", s.simulateCall)
	mux.HandleFunc("POST /api/simulate/bulk", s.simulateBulkCalls)
	mux.HandleFunc("POST /api/simulate/scenario", s.runScenario)
}

type simulateCallRequest struct {
	ServiceID       string `json:"service_id"`
	ForcedOutcome   string `json:"forced_outcome,omitempty"`
	ForcedLatencyMs int    `json:"forced_latency_ms,omitempty"`
}

func (s *Server) simulateCall(w http.ResponseWriter, r *http.Request) {
	var req simulateCallRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	res, err := s.svc.SimulateCall(req.ServiceID, req.ForcedOutcome, req.ForcedLatencyMs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, res)
}

type simulateBulkCallsRequest struct {
	ServiceID string `json:"service_id"`
	Count     int    `json:"count"`
}

func (s *Server) simulateBulkCalls(w http.ResponseWriter, r *http.Request) {
	var req simulateBulkCallsRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if req.Count <= 0 || req.Count > 1000 {
		httpx.BadRequest(w, "count 必须在 1~1000 之间")
		return
	}
	results := s.svc.SimulateBulkCalls(req.ServiceID, req.Count)
	httpx.OK(w, results)
}

type runScenarioRequest struct {
	ServiceID     string  `json:"service_id"`
	RequestCount  int     `json:"request_count"`
	FailureRate   float64 `json:"failure_rate"`
	AvgLatencyMs  int     `json:"avg_latency_ms"`
}

func (s *Server) runScenario(w http.ResponseWriter, r *http.Request) {
	var req runScenarioRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if req.RequestCount <= 0 || req.RequestCount > 10000 {
		httpx.BadRequest(w, "request_count 必须在 1~10000 之间")
		return
	}
	if req.FailureRate < 0 || req.FailureRate > 1 {
		httpx.BadRequest(w, "failure_rate 必须在 0~1 之间")
		return
	}
	result, err := s.svc.RunScenario(req.ServiceID, req.RequestCount, req.FailureRate, req.AvgLatencyMs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

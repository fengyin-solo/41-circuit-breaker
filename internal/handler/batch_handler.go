package handler

import (
	"net/http"

	"circuitbreaker/internal/model"
	"circuitbreaker/pkg/httpx"
)

func (s *Server) registerBatchRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/batch/services", s.batchCreateServices)
	mux.HandleFunc("POST /api/batch/services/delete", s.batchDeleteServices)
	mux.HandleFunc("POST /api/batch/breakers/transition", s.batchTransitionBreakers)
	mux.HandleFunc("POST /api/batch/records", s.batchCreateCallRecords)
	mux.HandleFunc("POST /api/batch/snapshots", s.batchCreateSnapshots)
}

type batchCreateServicesRequest struct {
	Services []model.DownstreamService `json:"services"`
}

func (s *Server) batchCreateServices(w http.ResponseWriter, r *http.Request) {
	var req batchCreateServicesRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	res := s.svc.BatchCreateServices(req.Services)
	httpx.OK(w, res)
}

type batchDeleteServicesRequest struct {
	IDs []string `json:"ids"`
}

func (s *Server) batchDeleteServices(w http.ResponseWriter, r *http.Request) {
	var req batchDeleteServicesRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	errors := s.svc.BatchDeleteServices(req.IDs)
	httpx.OK(w, map[string]interface{}{"errors": errors})
}

type batchTransitionBreakersRequest struct {
	IDs         []string `json:"ids"`
	TargetState string   `json:"target_state"`
}

func (s *Server) batchTransitionBreakers(w http.ResponseWriter, r *http.Request) {
	var req batchTransitionBreakersRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	res := s.svc.BatchTransitionBreakers(req.IDs, req.TargetState)
	httpx.OK(w, res)
}

type batchCreateCallRecordsRequest struct {
	Records []model.CallRecord `json:"records"`
}

func (s *Server) batchCreateCallRecords(w http.ResponseWriter, r *http.Request) {
	var req batchCreateCallRecordsRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	res := s.svc.BatchCreateCallRecords(req.Records)
	httpx.OK(w, res)
}

type batchCreateSnapshotsRequest struct {
	Snapshots []model.BreakerSnapshot `json:"snapshots"`
}

func (s *Server) batchCreateSnapshots(w http.ResponseWriter, r *http.Request) {
	var req batchCreateSnapshotsRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	res := s.svc.BatchCreateSnapshots(req.Snapshots)
	httpx.OK(w, res)
}

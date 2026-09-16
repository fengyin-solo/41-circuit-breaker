package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"circuitbreaker/internal/model"
	"circuitbreaker/pkg/httpx"
)

func (s *Server) registerBreakerSnapshotRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/snapshots", s.createBreakerSnapshot)
	mux.HandleFunc("GET /api/snapshots", s.listBreakerSnapshots)
	mux.HandleFunc("GET /api/snapshots/{id}", s.getBreakerSnapshot)
	mux.HandleFunc("PUT /api/snapshots/{id}", s.updateBreakerSnapshot)
	mux.HandleFunc("DELETE /api/snapshots/{id}", s.deleteBreakerSnapshot)
	mux.HandleFunc("GET /api/snapshots/export", s.exportSnapshots)
	mux.HandleFunc("POST /api/snapshots/import", s.importSnapshots)
}

type createBreakerSnapshotRequest struct {
	ServiceID       string  `json:"service_id"`
	State           string  `json:"state"`
	TotalCalls      int     `json:"total_calls"`
	SuccessCalls    int     `json:"success_calls"`
	FailureCalls    int     `json:"failure_calls"`
	SlowCalls       int     `json:"slow_calls"`
	FailureRatio    float64 `json:"failure_ratio"`
	SlowRatio       float64 `json:"slow_ratio"`
	SnapshotVersion int     `json:"snapshot_version"`
}

func (s *Server) createBreakerSnapshot(w http.ResponseWriter, r *http.Request) {
	var req createBreakerSnapshotRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	snap, err := s.svc.CreateBreakerSnapshot(model.BreakerSnapshot{
		ServiceID:       req.ServiceID,
		State:           req.State,
		TotalCalls:      req.TotalCalls,
		SuccessCalls:    req.SuccessCalls,
		FailureCalls:    req.FailureCalls,
		SlowCalls:       req.SlowCalls,
		FailureRatio:    req.FailureRatio,
		SlowRatio:       req.SlowRatio,
		SnapshotVersion: req.SnapshotVersion,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, snap)
}

func (s *Server) listBreakerSnapshots(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.BreakerSnapshotFilter{
		ServiceID: r.URL.Query().Get("service_id"),
		State:     r.URL.Query().Get("state"),
		Keyword:   r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListBreakerSnapshots(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getBreakerSnapshot(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	snap, err := s.svc.GetBreakerSnapshot(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, snap)
}

type updateBreakerSnapshotRequest struct {
	ServiceID       string  `json:"service_id"`
	State           string  `json:"state"`
	TotalCalls      int     `json:"total_calls"`
	SuccessCalls    int     `json:"success_calls"`
	FailureCalls    int     `json:"failure_calls"`
	SlowCalls       int     `json:"slow_calls"`
	FailureRatio    float64 `json:"failure_ratio"`
	SlowRatio       float64 `json:"slow_ratio"`
	SnapshotVersion int     `json:"snapshot_version"`
}

func (s *Server) updateBreakerSnapshot(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateBreakerSnapshotRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	snap, err := s.svc.UpdateBreakerSnapshot(id, model.BreakerSnapshot{
		ServiceID:       req.ServiceID,
		State:           req.State,
		TotalCalls:      req.TotalCalls,
		SuccessCalls:    req.SuccessCalls,
		FailureCalls:    req.FailureCalls,
		SlowCalls:       req.SlowCalls,
		FailureRatio:    req.FailureRatio,
		SlowRatio:       req.SlowRatio,
		SnapshotVersion: req.SnapshotVersion,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, snap)
}

func (s *Server) deleteBreakerSnapshot(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteBreakerSnapshot(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) exportSnapshots(w http.ResponseWriter, r *http.Request) {
	snaps := s.svc.GetServiceStats()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=snapshots.json")
	_ = json.NewEncoder(w).Encode(snaps)
}

type importSnapshotsRequest struct {
	Snapshots []model.BreakerSnapshot `json:"snapshots"`
}

func (s *Server) importSnapshots(w http.ResponseWriter, r *http.Request) {
	var req importSnapshotsRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	created := 0
	for _, snap := range req.Snapshots {
		snap.CreatedAt = time.Now()
		_, err := s.svc.CreateBreakerSnapshot(snap)
		if err == nil {
			created++
		}
	}
	httpx.OK(w, map[string]int{"created": created})
}

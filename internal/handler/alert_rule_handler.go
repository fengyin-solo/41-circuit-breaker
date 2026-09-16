package handler

import (
	"net/http"
	"strconv"

	"circuitbreaker/internal/model"
	"circuitbreaker/pkg/httpx"
)

func (s *Server) registerAlertRuleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/alert-rules", s.createAlertRule)
	mux.HandleFunc("GET /api/alert-rules", s.listAlertRules)
	mux.HandleFunc("GET /api/alert-rules/{id}", s.getAlertRule)
	mux.HandleFunc("PUT /api/alert-rules/{id}", s.updateAlertRule)
	mux.HandleFunc("DELETE /api/alert-rules/{id}", s.deleteAlertRule)
}

type createAlertRuleRequest struct {
	Name          string  `json:"name"`
	ServiceID     string  `json:"service_id"`
	Metric        string  `json:"metric"`
	Threshold     float64 `json:"threshold"`
	Severity      string  `json:"severity"`
	NotifyChannel string  `json:"notify_channel"`
	Enabled       bool    `json:"enabled"`
}

func (s *Server) createAlertRule(w http.ResponseWriter, r *http.Request) {
	var req createAlertRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	ar, err := s.svc.CreateAlertRule(model.AlertRule{
		Name:          req.Name,
		ServiceID:     req.ServiceID,
		Metric:        req.Metric,
		Threshold:     req.Threshold,
		Severity:      req.Severity,
		NotifyChannel: req.NotifyChannel,
		Enabled:       req.Enabled,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, ar)
}

func (s *Server) listAlertRules(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.AlertRuleFilter{
		Name:      r.URL.Query().Get("name"),
		ServiceID: r.URL.Query().Get("service_id"),
		Metric:    r.URL.Query().Get("metric"),
		Severity:  r.URL.Query().Get("severity"),
		Keyword:   r.URL.Query().Get("keyword"),
	}
	if enabledStr := r.URL.Query().Get("enabled"); enabledStr != "" {
		v, _ := strconv.ParseBool(enabledStr)
		filter.Enabled = &v
	}
	items, total, err := s.svc.ListAlertRules(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getAlertRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ar, err := s.svc.GetAlertRule(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, ar)
}

type updateAlertRuleRequest struct {
	Name          string  `json:"name"`
	ServiceID     string  `json:"service_id"`
	Metric        string  `json:"metric"`
	Threshold     float64 `json:"threshold"`
	Severity      string  `json:"severity"`
	NotifyChannel string  `json:"notify_channel"`
	Enabled       bool    `json:"enabled"`
}

func (s *Server) updateAlertRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateAlertRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	ar, err := s.svc.UpdateAlertRule(id, model.AlertRule{
		Name:          req.Name,
		ServiceID:     req.ServiceID,
		Metric:        req.Metric,
		Threshold:     req.Threshold,
		Severity:      req.Severity,
		NotifyChannel: req.NotifyChannel,
		Enabled:       req.Enabled,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, ar)
}

func (s *Server) deleteAlertRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteAlertRule(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

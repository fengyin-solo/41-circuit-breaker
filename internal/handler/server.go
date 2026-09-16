// Package handler 实现 HTTP 处理器层。
package handler

import (
	"errors"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"circuitbreaker/internal/config"
	"circuitbreaker/internal/model"
	"circuitbreaker/internal/service"
	"circuitbreaker/internal/store"
	"circuitbreaker/pkg/httpx"
	"circuitbreaker/pkg/logger"
)

type Server struct {
	svc *service.Service
	log *logger.Logger
	cfg *config.Config
}

func NewServer(svc *service.Service, log *logger.Logger, cfg *config.Config) *Server {
	return &Server{svc: svc, log: log, cfg: cfg}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	s.registerDownstreamServiceRoutes(mux)
	s.registerBreakerRuleRoutes(mux)
	s.registerCircuitBreakerRoutes(mux)
	s.registerCallRecordRoutes(mux)
	s.registerHealthCheckRoutes(mux)
	s.registerAlertRuleRoutes(mux)
	s.registerRecoveryPolicyRoutes(mux)
	s.registerMetricSampleRoutes(mux)
	s.registerBreakerEventRoutes(mux)
	s.registerBreakerSnapshotRoutes(mux)
	s.registerStatsRoutes(mux)
	s.registerBatchRoutes(mux)
	s.registerReportRoutes(mux)
	s.registerHealthCheckEvalRoutes(mux)
	s.registerBreakerSimulationRoutes(mux)
	mux.Handle("GET /", http.FileServer(http.Dir("web")))
	return s.authMiddleware(s.rateLimitMiddleware(s.loggingMiddleware(s.recoveryMiddleware(mux))))
}

func (s *Server) maxPageSize() int {
	if s.cfg != nil && s.cfg.MaxPageSize > 0 {
		return s.cfg.MaxPageSize
	}
	return 100
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		s.log.Infof("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func (s *Server) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				s.log.Errorf("panic: %v\n%s", rec, debug.Stack())
				httpx.InternalError(w, "服务器内部错误")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path == "/index.html" || r.URL.Path == "/style.css" || r.URL.Path == "/app.js" {
			next.ServeHTTP(w, r)
			return
		}
		key := r.Header.Get("X-Api-Key")
		if key == "" {
			key = r.URL.Query().Get("api_key")
		}
		if key != s.cfg.ApiKey {
			httpx.Unauthorized(w, "无效的 API Key")
			return
		}
		next.ServeHTTP(w, r)
	})
}

type rateLimitBucket struct {
	tokens float64
	last   time.Time
}

var rateLimitMu sync.RWMutex
var rateLimitBuckets = make(map[string]*rateLimitBucket)

const rateLimitRate = 10.0
const rateLimitBurst = 20

func (s *Server) rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		if xf := r.Header.Get("X-Forwarded-For"); xf != "" {
			clientIP = xf
		}
		rateLimitMu.Lock()
		bucket, ok := rateLimitBuckets[clientIP]
		if !ok {
			bucket = &rateLimitBucket{tokens: rateLimitBurst, last: time.Now()}
			rateLimitBuckets[clientIP] = bucket
		}
		now := time.Now()
		elapsed := now.Sub(bucket.last).Seconds()
		bucket.tokens += elapsed * rateLimitRate
		if bucket.tokens > rateLimitBurst {
			bucket.tokens = rateLimitBurst
		}
		bucket.last = now
		if bucket.tokens < 1 {
			rateLimitMu.Unlock()
			httpx.Error(w, http.StatusTooManyRequests, 429, "请求过于频繁")
			return
		}
		bucket.tokens--
		rateLimitMu.Unlock()
		next.ServeHTTP(w, r)
	})
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case model.IsValidationError(err):
		httpx.BadRequest(w, err.Error())
	case errors.Is(err, store.ErrNotFound):
		httpx.NotFound(w, err.Error())
	case errors.Is(err, store.ErrConflict):
		httpx.Conflict(w, err.Error())
	default:
		httpx.InternalError(w, err.Error())
	}
}

package model

import (
	"strings"
	"time"
)

const (
	ServiceProtocolHTTP = "http"
	ServiceProtocolGRPC = "grpc"
	ServiceStatusUp     = "up"
	ServiceStatusDown   = "down"
)

type DownstreamService struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Protocol  string    `json:"protocol"`
	TimeoutMs int       `json:"timeout_ms"`
	Status    string    `json:"status"`
	Weight    int       `json:"weight"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *DownstreamService) Validate() error {
	s.Name = strings.TrimSpace(s.Name)
	s.Address = strings.TrimSpace(s.Address)
	if s.Name == "" {
		return NewValidationError("name", "服务名称不能为空")
	}
	if s.Address == "" {
		return NewValidationError("address", "服务地址不能为空")
	}
	if s.Protocol == "" {
		s.Protocol = ServiceProtocolHTTP
	}
	if s.Protocol != ServiceProtocolHTTP && s.Protocol != ServiceProtocolGRPC {
		return NewValidationError("protocol", "协议必须为 http 或 grpc")
	}
	if s.TimeoutMs <= 0 {
		s.TimeoutMs = 5000
	}
	if s.Status == "" {
		s.Status = ServiceStatusUp
	}
	if s.Status != ServiceStatusUp && s.Status != ServiceStatusDown {
		return NewValidationError("status", "状态必须为 up 或 down")
	}
	if s.Weight < 0 {
		s.Weight = 1
	}
	if s.Weight > 100 {
		return NewValidationError("weight", "权重不能超过 100")
	}
	return nil
}

type DownstreamServiceFilter struct {
	Name     string
	Protocol string
	Status   string
	Keyword  string
}

func (f DownstreamServiceFilter) Match(s *DownstreamService) bool {
	if f.Name != "" && s.Name != f.Name {
		return false
	}
	if f.Protocol != "" && s.Protocol != f.Protocol {
		return false
	}
	if f.Status != "" && s.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(s.Name), k) &&
			!strings.Contains(strings.ToLower(s.Address), k) {
			return false
		}
	}
	return true
}

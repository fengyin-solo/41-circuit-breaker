# 熔断器管理系统 (Circuit Breaker)

纯 Go 标准库实现的熔断器后端服务，零第三方依赖。

## 运行说明

```bash
cd origin
go run ./cmd/server
```

默认监听 `:8080`，API Key 为 `circuitbreaker-secret-key`（通过 Header `X-Api-Key` 或 Query `api_key` 传递）。

访问 `http://localhost:8080/` 可查看前端看板页面。

## API 表格

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/services | 创建下游服务 |
| GET | /api/services | 列出下游服务（支持 name/protocol/status/keyword 筛选） |
| GET | /api/services/{id} | 获取下游服务 |
| PUT | /api/services/{id} | 更新下游服务 |
| DELETE | /api/services/{id} | 删除下游服务 |
| POST | /api/rules | 创建熔断规则 |
| GET | /api/rules | 列出熔断规则（支持 name/service_id/enabled/keyword 筛选） |
| GET | /api/rules/{id} | 获取熔断规则 |
| PUT | /api/rules/{id} | 更新熔断规则 |
| DELETE | /api/rules/{id} | 删除熔断规则 |
| POST | /api/breakers | 创建熔断器实例 |
| GET | /api/breakers | 列出熔断器（支持 service_id/rule_id/state/keyword 筛选） |
| GET | /api/breakers/{id} | 获取熔断器 |
| PUT | /api/breakers/{id} | 更新熔断器（含状态机校验） |
| DELETE | /api/breakers/{id} | 删除熔断器 |
| POST | /api/breakers/{id}/transition | 状态流转（closed->open->half_open->closed） |
| POST | /api/records | 创建调用记录 |
| GET | /api/records | 列出调用记录（支持 service_id/outcome/keyword 筛选） |
| GET | /api/records/{id} | 获取调用记录 |
| DELETE | /api/records/{id} | 删除调用记录 |
| POST | /api/health-checks | 创建健康检查 |
| GET | /api/health-checks | 列出健康检查（支持 service_id/last_status/keyword 筛选） |
| GET | /api/health-checks/{id} | 获取健康检查 |
| PUT | /api/health-checks/{id} | 更新健康检查 |
| DELETE | /api/health-checks/{id} | 删除健康检查 |
| POST | /api/alert-rules | 创建告警规则 |
| GET | /api/alert-rules | 列出告警规则（支持 name/service_id/metric/severity/enabled/keyword 筛选） |
| GET | /api/alert-rules/{id} | 获取告警规则 |
| PUT | /api/alert-rules/{id} | 更新告警规则 |
| DELETE | /api/alert-rules/{id} | 删除告警规则 |
| POST | /api/policies | 创建恢复策略 |
| GET | /api/policies | 列出恢复策略（支持 name/service_id/keyword 筛选） |
| GET | /api/policies/{id} | 获取恢复策略 |
| PUT | /api/policies/{id} | 更新恢复策略 |
| DELETE | /api/policies/{id} | 删除恢复策略 |
| POST | /api/metrics | 创建指标采样 |
| GET | /api/metrics | 列出指标采样（支持 service_id/keyword 筛选） |
| GET | /api/metrics/{id} | 获取指标采样 |
| DELETE | /api/metrics/{id} | 删除指标采样 |
| POST | /api/events | 创建熔断事件 |
| GET | /api/events | 列出熔断事件（支持 breaker_id/service_id/event_type/keyword 筛选） |
| GET | /api/events/{id} | 获取熔断事件 |
| DELETE | /api/events/{id} | 删除熔断事件 |
| POST | /api/snapshots | 创建熔断快照 |
| GET | /api/snapshots | 列出熔断快照（支持 service_id/state/keyword 筛选） |
| GET | /api/snapshots/{id} | 获取熔断快照 |
| PUT | /api/snapshots/{id} | 更新熔断快照 |
| DELETE | /api/snapshots/{id} | 删除熔断快照 |
| GET | /api/snapshots/export | 导出快照 JSON |
| POST | /api/snapshots/import | 导入快照 JSON |
| GET | /api/stats/overview | 全局统计概览 |
| GET | /api/stats/services | 按服务统计 |
| GET | /api/stats/breaker-states | 按熔断器状态分组统计 |
| GET | /api/stats/top-failures?n=5 | TOP N 失败服务 |
| GET | / | 前端看板页面 |

## 实体列表

1. DownstreamService（下游服务）
2. BreakerRule（熔断规则）
3. CircuitBreaker（熔断器实例）
4. CallRecord（调用记录）
5. HealthCheck（健康检查）
6. AlertRule（告警规则）
7. RecoveryPolicy（恢复策略）
8. MetricSample（指标采样）
9. BreakerEvent（熔断事件）
10. BreakerSnapshot（熔断快照）

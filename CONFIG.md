# 中间件配置说明

## Temporal 工作流引擎

### 本地开发环境

#### 服务地址
| 服务 | 地址 | 用途 |
|------|------|------|
| **Temporal Server** | `localhost:7233` | gRPC API 端口，用于连接Temporal |
| **Temporal UI** | `http://localhost:8080` | Web控制台，可视化工作流执行状态 |
| **Temporal Metrics** | `http://localhost:60471/metrics` | Prometheus指标端点，用于监控 |

#### 连接配置

```go
// Go 客户端连接配置
package temporal

import (
    "go.temporal.io/sdk/client"
)

const (
    TemporalAddress = "localhost:7233"
)

func NewTemporalClient() (client.Client, error) {
    return client.Dial(client.Options{
        HostPort: TemporalAddress,
    })
}
```

```yaml
# temporal.yaml (配置文件示例)
temporal:
  address: localhost:7233
  namespace: default
  task_queue: develop-workflow
  timeout: 30s
```

#### 环境变量

```bash
# .env 文件
TEMPORAL_ADDRESS=localhost:7233
TEMPORAL_NAMESPACE=default
TEMPORAL_TASK_QUEUE=develop-workflow
TEMPORAL_UI_URL=http://localhost:8080
```

#### Docker Compose 配置

```yaml
# docker-compose.yml
version: '3.8'

services:
  temporal:
    image: temporalio/auto-setup:1.30.0
    ports:
      - "7233:7233"
      - "8080:8080"  # UI
      - "61012:61012"  # Metrics
    environment:
      - TEMPORAL_ADDRESS=localhost:7233
      - TEMPORAL_NAMESPACE=default

  temporal-ui:
    image: temporalio/ui:2.45.0
    ports:
      - "8080:8080"
    environment:
      - TEMPORAL_ADDRESS=temporal:7233
    depends_on:
      - temporal
```

### 生产环境配置

```yaml
# production-temporal.yaml
temporal:
  address: temporal.example.com:7233
  namespace: production
  tls:
    enabled: true
    cert_file: /path/to/cert.pem
    key_file: /path/to/key.pem
  metrics:
    enabled: true
    address: :8080
```

---

## Redis 状态存储

### 本地开发环境

#### 服务地址
| 服务 | 地址 | 用途 |
|------|------|------|
| **Redis Server** | `localhost:6379` | 主服务端口 |
| **Redis Sentinel** | `localhost:26379` | 高可用监控（可选） |

#### 连接配置

```go
// Go Redis 客户端配置
package storage

import (
    "github.com/redis/go-redis/v9"
)

const (
    RedisAddr     = "localhost:6379"
    RedisPassword = ""
    RedisDB       = 0
)

func NewRedisClient() *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr:     RedisAddr,
        Password: RedisPassword,
        DB:       RedisDB,
    })
}
```

```yaml
# redis.yaml
redis:
  address: localhost:6379
  password: ""
  db: 0
  pool_size: 10
  min_idle_conns: 5
```

#### 环境变量

```bash
# .env 文件
REDIS_ADDRESS=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
```

#### Docker Compose 配置

```yaml
# docker-compose.yml
services:
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    command: redis-server --appendonly yes

volumes:
  redis_data:
```

---

## NATS 消息系统

### 本地开发环境

#### 服务地址
| 服务 | 地址 | 用途 |
|------|------|------|
| **NATS Server** | `localhost:4222` | 客户端连接端口 |
| **NATS Monitoring** | `http://localhost:8222` | 监控端点 |

#### 连接配置

```go
// Go NATS 客户端配置
package messaging

import (
    "github.com/nats-io/nats.go"
)

const (
    NATSAddress = "nats://localhost:4222"
)

func NewNATSConnection() (*nats.Conn, error) {
    return nats.Connect(NATSAddress)
}
```

```yaml
# nats.yaml
nats:
  address: nats://localhost:4222
  cluster_id: develop-workflow
  client_id: workflow-client-1
```

#### 环境变量

```bash
# .env 文件
NATS_ADDRESS=nats://localhost:4222
NATS_CLUSTER_ID=develop-workflow
NATS_CLIENT_ID=workflow-client-1
```

#### Docker Compose 配置

```yaml
# docker-compose.yml
services:
  nats:
    image: nats:2.10
    ports:
      - "4222:4222"
      - "8222:8222"  # Monitoring
    command: "--http_port 8222"
```

---

## GitHub API 配置

### 认证配置

```yaml
# github.yaml
github:
  token: ${GITHUB_TOKEN}
  owner: your-username
  private: true
  webhook_secret: ${GITHUB_WEBHOOK_SECRET}
```

#### 环境变量

```bash
# .env 文件
GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxx
GITHUB_WEBHOOK_SECRET=your-webhook-secret
GITHUB_OWNER=your-username
```

---

## AI 模型配置

### Claude API

```yaml
# claude.yaml
claude:
  api_key: ${ANTHROPIC_API_KEY}
  model: claude-sonnet-4-20250514
  max_tokens: 4096
  temperature: 0.7
```

### OpenAI API

```yaml
# openai.yaml
openai:
  api_key: ${OPENAI_API_KEY}
  model: gpt-4
  max_tokens: 4096
  temperature: 0.7
```

#### 环境变量

```bash
# .env 文件
ANTHROPIC_API_KEY=sk-ant-xxxxxxxxxxxx
OPENAI_API_KEY=sk-xxxxxxxxxxxx
```

---

## 完整环境变量配置

### .env 文件模板

```bash
# ============================================
# Temporal 配置
# ============================================
TEMPORAL_ADDRESS=localhost:7233
TEMPORAL_NAMESPACE=default
TEMPORAL_TASK_QUEUE=develop-workflow
TEMPORAL_UI_URL=http://localhost:8080

# ============================================
# Redis 配置
# ============================================
REDIS_ADDRESS=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# ============================================
# NATS 配置
# ============================================
NATS_ADDRESS=nats://localhost:4222
NATS_CLUSTER_ID=develop-workflow
NATS_CLIENT_ID=workflow-client-1

# ============================================
# GitHub 配置
# ============================================
GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxx
GITHUB_WEBHOOK_SECRET=your-webhook-secret
GITHUB_OWNER=your-username

# ============================================
# AI 模型配置
# ============================================
ANTHROPIC_API_KEY=sk-ant-xxxxxxxxxxxx
OPENAI_API_KEY=sk-xxxxxxxxxxxx

# ============================================
# 应用配置
# ============================================
APP_ENV=development
LOG_LEVEL=debug
MAX_RETRIES=3
TIMEOUT=30s
```

---

## 快速启动

### 1. 使用 Docker Compose（推荐）

```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f temporal
```

### 2. 手动启动

```bash
# 启动 Temporal
temporal server start-dev

# 启动 Redis
redis-server

# 启动 NATS
nats-server
```

### 3. 验证服务

```bash
# 检查 Temporal
temporal operator namespace describe default

# 检查 Redis
redis-cli ping

# 检查 NATS
nats-server --help
```

---

## 健康检查

### Temporal

```bash
# 检查服务状态
temporal operator cluster health

# 检查命名空间
temporal operator namespace list
```

### Redis

```bash
# 连接测试
redis-cli ping
# 应返回: PONG

# 查看信息
redis-cli info server
```

### NATS

```bash
# 查看监控
curl http://localhost:8222/varz
```

---

## 常见问题

### 1. 端口冲突

如果端口被占用，修改配置：

```yaml
# 修改端口
temporal:
  address: localhost:7234  # 改为其他端口
```

### 2. 连接失败

检查服务是否启动：

```bash
# Linux/Mac
lsof -i :7233
lsof -i :6379
lsof -i :4222

# Windows
netstat -ano | findstr :7233
```

### 3. 权限问题

确保文件权限正确：

```bash
chmod 600 .env
chmod 700 config/
```

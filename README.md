# Develop Workflow

基于 Temporal 和 AI Agent 的自动化开发工作流系统。

## 项目结构

```
develop-workflow/
├── cmd/
│   ├── worker/          # Temporal Worker 入口
│   └── cli/             # CLI 命令行工具
├── internal/
│   ├── workflow/        # Temporal Workflow 定义
│   ├── activity/        # Temporal Activity 定义
│   ├── agent/           # AI Agent 实现
│   └── config/          # 配置管理
├── pkg/
│   └── models/          # 数据模型
├── configs/             # 配置文件
├── docker-compose.yml   # Docker 编排
├── Makefile            # 构建脚本
└── .env                # 环境变量
```

## 快速开始

### 1. 启动本地服务

```bash
# 使用 Docker Compose 启动所有依赖服务
make up

# 或者手动启动
docker-compose up -d
```

服务启动后：
- Temporal Server: `localhost:7233`
- Temporal UI: `http://localhost:8080`
- Redis: `localhost:6379`
- NATS: `localhost:4222`

### 2. 配置环境变量

复制 `.env.example` 为 `.env` 并填入配置：

```bash
cp .env.example .env
```

编辑 `.env` 文件，填入：
- `GITHUB_TOKEN`: GitHub Personal Access Token
- `ANTHROPIC_API_KEY`: Claude API Key

### 3. 构建项目

```bash
# 构建所有二进制文件
make build

# 或者分别构建
make build-worker  # 构建 Worker
make build-cli     # 构建 CLI
```

### 4. 运行 Worker

```bash
# 运行 Worker（监听任务队列）
make worker

# 或者直接运行
./bin/worker
```

### 5. 启动开发工作流

```bash
# 使用 CLI 启动工作流
./bin/cli <project-name> "<description>"

# 示例：创建一个 TODO 应用
./bin/cli my-todo-app "Create a simple TODO application with CRUD operations"
```

## 工作流程

1. **需求分析**: AI Agent 分析用户输入，生成 PRD 文档
2. **架构设计**: AI Agent 设计系统架构和技术选型
3. **任务规划**: 将需求拆分为可执行的任务
4. **代码实现**: AI Agent 根据技术规范生成代码
5. **测试验证**: 自动执行测试并生成报告
6. **文档生成**: 自动生成项目文档

## 开发指南

### 运行测试

```bash
# 运行所有测试
make test

# 运行特定包的测试
go test ./internal/config/...
```

### 代码格式化

```bash
# 格式化代码
make fmt

# 检查代码
make lint
```

### 查看日志

```bash
# 查看 Docker 服务日志
make logs
```

## 技术栈

- **工作流引擎**: Temporal
- **AI 模型**: Claude API (Anthropic)
- **状态存储**: Redis
- **消息系统**: NATS
- **开发语言**: Go

## 配置说明

详细的中间件配置请参考 [CONFIG.md](./CONFIG.md)

## MVP 范围

详细的 MVP 定义请参考 [MVP.md](./MVP.md)

## 系统设计

详细的系统设计请参考 [DESIGN.md](./DESIGN.md)

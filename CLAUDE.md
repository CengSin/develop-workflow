# 设计与规划

## README

[README.md](./README.md) 项目说明文档，包含：
- 项目结构
- 快速开始指南
- 工作流程说明
- 开发指南

## DESIGN.md

[DESIGN.md](./DESIGN.md) 定义了整体系统的设计，请参考该文件

## MVP

[MVP.md](./MVP.md) 定义了项目的最小可行产品，包含：
- MVP目标与成功标准
- 5个核心功能的用户故事与验收标准
- 不包含的功能（Out of Scope）
- 技术实现要点
- 开发计划（5周）
- 验收标准

## 中间件配置

[CONFIG.md](./CONFIG.md) 定义了所有中间件的本地开发配置，包含：
- Temporal（localhost:7233, UI: http://localhost:8080）
- Redis（localhost:6379）
- NATS（localhost:4222）
- GitHub API 配置
- AI 模型配置（Claude/OpenAI）
- 完整的环境变量模板
- Docker Compose 配置
- 快速启动指南
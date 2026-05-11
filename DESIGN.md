## Workflow 流程设计

### 第一阶段：需求对齐与资产初始化

1. **需求萃取（Requirement Extraction）：** Agent 接收模糊需求，通过对话引导用户明确核心功能、目标用户和技术栈。
2. **仓库与基准（Initialization）：**
* Agent 调用 GitHub API 创建 Private 仓库。
* 生成标准化的 `PRD.md`（需求说明书）和 `ARCHITECTURE.md`（架构设计文档）。
* **关键点：** 此时应包含一个 `context_window`（上下文窗口）管理策略，确保后续步骤能引用此基准。

### 第二阶段：人工网关与契约确认

3. **人工审核（Human-in-the-loop）：** 用户在 GitHub 上查看 PR（Pull Request）或 Issue。不合格则通过 Comment（评论）反馈，Agent 监听 Webhook 自动触发修订。
4. **MVP 拆解与技术选型：**
* Agent 将总需求拆分为一系列 `Milestones`（里程碑）。
* **人机协作：** 双方共同确认第一个 `Iteration_0` 的范围。
* 输出 `SPEC.md`（技术规范文档）和 `TEST_PLAN.md`（测试计划）。

### 第三阶段：循环开发与自动化验证

5. **Coding Agent 任务分发：**
* **实现（Implementation）：** Coding Agent 接收特定的里程碑任务，编写代码。
* **单元测试（Unit Testing）：** Coding Agent 必须同步输出测试用例。
* **CI/CD 集成：** 自动跑通测试。如果 CI 失败，Coding Agent 自动进行 `Self-Correction`（自愈/纠错）。

6. **交付物生成：** 输出 `USER_GUIDE.md`（用户手册）和 `API_REFERENCE.md`（API 参考）。

### 第四阶段：状态机回溯

7. **定义“完成”（Definition of Done）：**
* Agent 调用 `Validator` 模块，对比 `PRD.md` 的功能列表与当前代码的 `Coverage`（覆盖率）。
* **逻辑判断：** 如果 `Done_List < Total_List`，自动生成下一个里程碑的 Prompt，回到步骤 4。

---

## 优化建议与潜在风险提示

### 1. 明确环境与前提

* **依赖管理：** 在第一步创建仓库时，必须要求 Agent 探测或询问用户的运行环境（如 Go 版本、Java JDK 版本、数据库类型）。
* **权限控制：** 建议使用 GitHub App 而非个人 Token，以精细化管理仓库写权限。

### 2. 潜在风险与失败场景

* **上下文丢失（Context Drift）：** 随着项目迭代，`PRD` 可能会变得非常庞大。Coding Agent 在第五步可能会忘记第一步定义的全局约束。
* **对策：** 引入 **RAG (检索增强生成)** 或每轮迭代生成的 **Context Snapshots (上下文快照)**。


* **死循环风险：** 如果 Coding Agent 无法通过测试，可能会在步骤 5 反复重试消耗大量 Token。
* **对策：** 设置 `Max_Retries`（最大重试次数），超过阈值必须人工介入。



### 3. 验证方案

* 建议在步骤 4 和 5 之间增加一个 **"Dry Run" (演练)** 环节：Agent 先输出伪代码或接口定义（Interface Definition），由人确认无误后再进行大规模编码。

---

## 技术栈

### 核心组件

| 层面 | 技术选型 | 理由 |
|------|----------|------|
| **工作流引擎** | Temporal | 支持长时间运行、有状态的工作流，可观测性好，天然适合开发流程编排 |
| **Agent框架** | Claude SDK + OpenAI SDK | 双SDK支持，保持灵活性，可根据任务特性选择模型 |
| **开发语言** | Go | 两个SDK都有Go版本，性能好，并发处理能力强 |
| **状态存储** | Redis | Temporal外需要轻量级状态存储，用于Agent间共享上下文 |
| **事件通信** | NATS | 轻量级消息系统，支持Agent间解耦通信 |

### 技术选型详细说明

1. **Temporal** - 工作流编排核心
   - 支持长时间运行的工作流（开发流程可能持续数小时/天）
   - 内置重试、超时、状态持久化
   - 优秀的可观测性（UI、日志、指标）
   - 天然支持Human-in-the-loop（通过Signal/Query）

2. **Claude SDK / OpenAI SDK** - Agent能力核心
   - 用于需求分析、代码生成、测试生成等AI任务
   - 根据任务复杂度和成本选择模型

3. **NATS** - Agent间通信
   - 轻量、高性能、支持发布/订阅模式
   - 解耦Agent间的直接依赖
   - 支持消息持久化和重放

---

## Agent间通信协议设计

### 1. 统一消息格式

所有Agent间通信使用标准化的JSON消息格式：

```json
{
  "message_id": "uuid-v4",
  "timestamp": "2026-05-08T10:00:00Z",
  "sender": {
    "agent_id": "coding-agent-1",
    "agent_type": "implementation",
    "task_id": "task-uuid"
  },
  "receiver": {
    "agent_id": "review-agent",
    "agent_type": "validation"
  },
  "type": "task_result|status_update|error|context_request",
  "payload": {
    "content": {},
    "artifacts": ["file_path_1", "file_path_2"]
  },
  "context": {
    "snapshot_id": "ctx-snapshot-uuid",
    "workflow_id": "temporal-workflow-id",
    "iteration": 3
  }
}
```

### 2. Agent类型定义

| Agent类型 | 职责 | 输入 | 输出 |
|-----------|------|------|------|
| **Requirement Agent** | 需求分析与PRD生成 | 用户输入 | PRD.md, ARCHITECTURE.md |
| **Planning Agent** | 里程碑拆解与任务分配 | PRD.md | SPEC.md, TASKS.md |
| **Coding Agent** | 代码实现 | SPEC.md + 任务描述 | 代码文件 + 单元测试 |
| **Review Agent** | 代码审查 | 代码变更 | Review报告 |
| **Validation Agent** | 测试验证 | 代码 + 测试计划 | 测试结果报告 |

### 3. 通信模式

#### 3.1 任务分发（Task Distribution）
```
Planning Agent → NATS → Coding Agent
Topic: agent.task.assign.{agent_id}
```

#### 3.2 结果上报（Result Reporting）
```
Coding Agent → NATS → Review Agent
Topic: agent.result.submit.{workflow_id}
```

#### 3.3 状态同步（State Synchronization）
```
所有Agent → Redis → 共享状态
Key: workflow:{workflow_id}:state
```

#### 3.4 上下文请求（Context Request）
```
Agent → Redis → 加载上下文快照
Key: workflow:{workflow_id}:context:{snapshot_id}
```

### 4. 上下文管理策略

解决DESIGN.md中提到的"Context Drift"问题：

#### 4.1 上下文快照（Context Snapshot）
- 每个里程碑完成后，自动生成上下文快照
- 快照包含：PRD摘要、当前状态、已完成任务、关键决策
- 存储到Redis，通过snapshot_id引用

#### 4.2 增量上下文传递
- Agent间通信只传递差异部分（delta）
- 接收方基于已有上下文 + delta重建完整状态
- 减少网络传输和Token消耗

#### 4.3 RAG检索增强
- 对于大型项目，建立向量索引
- Agent可检索相关代码片段和历史决策
- 避免每次加载完整上下文

### 5. 错误恢复机制

#### 5.1 重试策略
```go
type RetryPolicy struct {
    MaxRetries:     3,
    BackoffType:    "exponential",
    InitialBackoff: 5 * time.Second,
    MaxBackoff:     5 * time.Minute,
}
```

#### 5.2 人工介入触发条件
- 重试次数超过阈值
- 检测到死循环（相同错误重复出现）
- 涉及安全或架构重大变更

#### 5.3 状态回滚
- 每个里程碑开始前保存检查点
- 失败时可回滚到上一个检查点
- 避免错误累积

---

## 规划评估

### 现有设计的优点

1. **流程完整性**
   - 覆盖从需求到交付的完整生命周期
   - 包含反馈循环（Human-in-the-loop）
   - 考虑了验证和回溯机制

2. **风险意识**
   - 识别了关键风险（上下文丢失、死循环）
   - 提出了相应的对策（RAG、Max_Retries）
   - 有验证机制确保质量

3. **可扩展性**
   - 模块化设计，各阶段相对独立
   - 支持不同技术栈的Agent
   - Temporal支持分布式部署

### 需要补充的内容

1. **状态机定义**
   - 建议明确定义所有状态（如：pending, in_progress, review, testing, done, failed）
   - 定义状态转换条件和触发事件
   - 绘制状态转换图

2. **Agent能力边界**
   - 明确每个Agent的职责范围
   - 定义Agent无法处理时的升级机制
   - 避免Agent越权操作

3. **性能指标**
   - 定义关键指标（任务完成时间、成功率、Token消耗）
   - 建立监控和告警机制
   - 用于优化和成本控制

4. **安全考虑**
   - Agent权限控制（最小权限原则）
   - 代码执行沙箱
   - 敏感信息处理

### 建议的下一步

1. **Phase 1: 基础框架搭建**
   - 实现Temporal工作流基础
   - 实现标准化消息协议
   - 搭建Agent间通信基础设施

2. **Phase 2: 核心Agent实现**
   - 实现Requirement Agent和Planning Agent
   - 实现基础的Coding Agent
   - 实现简单的Validation Agent

3. **Phase 3: 集成与优化**
   - 完整流程集成测试
   - 性能优化和成本控制
   - 监控和可观测性建设

---

## MVP

详细的MVP定义请参考 [MVP.md](./MVP.md)
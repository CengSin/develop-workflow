# MVP 定义

## MVP 目标

**核心价值主张：** 用户输入一个模糊需求，系统自动完成从需求分析到代码生成的完整流程，输出可运行的代码仓库。

**成功标准：**
- 用户无需手动干预，系统能自主完成一个简单项目的开发
- 输出的代码能通过基本测试
- 整个流程可在30分钟内完成（简单项目）

## MVP 功能范围

### 功能1：需求理解与PRD生成（Requirement Agent）

**用户故事：**
```
作为用户，我希望输入一个简单的需求描述（如"创建一个TODO应用"），
系统能自动分析并生成标准化的PRD文档。
```

**功能点：**
- [ ] 接收用户自然语言输入
- [ ] 调用LLM分析需求，提取核心功能点
- [ ] 生成标准化PRD.md文档
- [ ] 生成基础ARCHITECTURE.md（技术选型、目录结构）
- [ ] 用户可在GitHub上查看PRD并评论反馈

**输入：** 用户需求描述（文本）
**输出：** PRD.md, ARCHITECTURE.md

**验收标准：**
- PRD包含：功能列表、用户故事、技术约束
- ARCHITECTURE包含：技术栈、目录结构、关键接口定义
- 用户可在GitHub PR中查看并Approve

---

### 功能2：任务拆解与规划（Planning Agent）

**用户故事：**
```
作为系统，我需要将PRD拆分为可执行的开发任务，
每个任务对应一个具体的代码模块。
```

**功能点：**
- [ ] 解析PRD中的功能列表
- [ ] 按依赖关系排序任务
- [ ] 生成任务清单（TASKS.md）
- [ ] 为每个任务生成具体的技术规范（SPEC）
- [ ] 确定任务执行顺序（DAG）

**输入：** PRD.md, ARCHITECTURE.md
**输出：** TASKS.md, SPEC_*.md（每个任务一个）

**验收标准：**
- 任务拆分粒度合理（每个任务1-3小时工作量）
- 任务依赖关系正确
- 每个SPEC包含：输入/输出、接口定义、边界条件

---

### 功能3：代码实现（Coding Agent）

**用户故事：**
```
作为系统，我需要根据SPEC自动编写代码和单元测试。
```

**功能点：**
- [ ] 读取当前任务的SPEC
- [ ] 加载相关上下文（PRD摘要、已完成任务）
- [ ] 生成代码实现
- [ ] 生成对应的单元测试
- [ ] 执行本地测试验证
- [ ] 如果测试失败，自动修复（最多3次）

**输入：** SPEC_*.md, 上下文快照
**输出：** 代码文件 + 测试文件

**验收标准：**
- 代码符合SPEC定义的接口
- 单元测试覆盖核心逻辑
- 本地测试通过

---

### 功能4：测试验证（Validation Agent）

**用户故事：**
```
作为系统，我需要验证生成的代码是否满足需求，
并生成测试报告。
```

**功能点：**
- [ ] 执行单元测试套件
- [ ] 生成测试报告（通过率、覆盖率）
- [ ] 如果测试失败，触发Coding Agent重试
- [ ] 如果重试超过阈值，标记为失败并通知用户

**输入：** 代码 + 测试文件
**输出：** TEST_REPORT.md

**验收标准：**
- 测试通过率 > 80%
- 核心功能有测试覆盖
- 失败时有清晰的错误信息

---

### 功能5：工作流编排（Temporal Workflow）

**用户故事：**
```
作为系统，我需要协调所有Agent按正确顺序执行，
并处理异常情况。
```

**功能点：**
- [ ] 定义Temporal工作流（需求→规划→实现→验证）
- [ ] 实现Agent间的状态传递
- [ ] 实现Human-in-the-loop（PRD审核）
- [ ] 实现重试机制（测试失败自动重试）
- [ ] 实现超时控制（单任务最长30分钟）

**输入：** 用户需求
**输出：** 完整的代码仓库

**验收标准：**
- 工作流可正常执行
- 状态可持久化，支持恢复
- 异常情况有处理

---

## MVP 不包含的功能（Out of Scope）

为了控制MVP范围，以下功能暂时不实现：

1. **❌ 复杂项目支持**
   - 只支持单体应用，不支持微服务
   - 只支持单一语言（如Go或Python）

2. **❌ 高级上下文管理**
   - 不实现RAG
   - 不实现向量索引
   - 使用简单的快照机制

3. **❌ 完整的CI/CD**
   - 不集成GitHub Actions
   - 只做本地测试验证

4. **❌ 多轮迭代**
   - 只实现单次开发流程
   - 不支持"完成后再开发新功能"

5. **❌ 用户界面**
   - 只通过GitHub交互
   - 不开发Web UI

6. **❌ 成本优化**
   - 不实现Token消耗监控
   - 不实现模型选择优化

---

## MVP 技术实现要点

### 1. 工作流定义（Temporal）

```go
func DevelopWorkflow(ctx workflow.Context, req DevelopRequest) (*DevelopResult, error) {
    // Step 1: 需求分析
    var prd PRD
    err := workflow.GetSignalChannel(ctx, "prd_approved").Receive(ctx, &prd)

    // Step 2: 任务规划
    var tasks []Task
    err = workflow.ExecuteActivity(ctx, PlanningActivity, prd).Get(ctx, &tasks)

    // Step 3: 逐个任务实现
    for _, task := range tasks {
        var code []File
        err = workflow.ExecuteActivity(ctx, CodingActivity, task).Get(ctx, &code)

        // Step 4: 测试验证
        var report TestReport
        err = workflow.ExecuteActivity(ctx, ValidationActivity, code).Get(ctx, &report)

        if !report.Passed {
            // 重试逻辑
        }
    }

    return &DevelopResult{Success: true}, nil
}
```

### 2. Agent通信（简化版）

MVP阶段使用函数调用而非NATS：

```go
type AgentOrchestrator struct {
    requirementAgent *RequirementAgent
    planningAgent    *PlanningAgent
    codingAgent      *CodingAgent
    validationAgent  *ValidationAgent
}

func (o *AgentOrchestrator) Execute(ctx context.Context, req DevelopRequest) error {
    // 顺序调用各Agent
    prd, _ := o.requirementAgent.Analyze(ctx, req)
    tasks, _ := o.planningAgent.Plan(ctx, prd)

    for _, task := range tasks {
        code, _ := o.codingAgent.Implement(ctx, task)
        report, _ := o.validationAgent.Validate(ctx, code)

        if !report.Passed {
            // 重试
            code, _ = o.codingAgent.Retry(ctx, task, report.Errors)
        }
    }
    return nil
}
```

### 3. 上下文管理（简化版）

MVP使用文件系统存储上下文：

```
project/
├── .workflow/
│   ├── context/
│   │   ├── snapshot_1.json
│   │   ├── snapshot_2.json
│   │   └── ...
│   ├── prd.md
│   ├── tasks.md
│   └── state.json
├── src/
└── tests/
```

---

## MVP 开发计划

### Phase 1: 基础框架（1周）
- [ ] 搭建Go项目结构
- [ ] 集成Temporal SDK
- [ ] 实现基础工作流定义
- [ ] 实现Agent基础接口

### Phase 2: 核心Agent（2周）
- [ ] 实现Requirement Agent（调用Claude API）
- [ ] 实现Planning Agent
- [ ] 实现基础Coding Agent
- [ ] 实现基础Validation Agent

### Phase 3: 集成测试（1周）
- [ ] 端到端流程测试
- [ ] 修复集成问题
- [ ] 优化Prompt模板

### Phase 4: 示例项目（1周）
- [ ] 用MVP开发一个简单项目（如TODO应用）
- [ ] 验证流程完整性
- [ ] 收集反馈并优化

**预计总工期：5周**

---

## MVP 验收标准

### 功能验收
- [ ] 用户输入"创建一个TODO应用"，系统能输出完整的代码仓库
- [ ] 输出的代码能通过单元测试
- [ ] 用户可在GitHub上查看PRD并审核
- [ ] 整个流程无需用户手动干预

### 质量验收
- [ ] 生成的代码符合Go/Python最佳实践
- [ ] 单元测试覆盖核心逻辑（>70%）
- [ ] 工作流执行稳定，无死循环

### 性能验收
- [ ] 简单项目（<10个文件）30分钟内完成
- [ ] Token消耗在合理范围内（<$5）

package models

import "time"

// DevelopRequest 开发请求
type DevelopRequest struct {
	Description string `json:"description"`
	ProjectName string `json:"project_name"`
	Language    string `json:"language"`
}

// DevelopResult 开发结果
type DevelopResult struct {
	Success    bool   `json:"success"`
	Repository string `json:"repository"`
	Error      string `json:"error,omitempty"`
}

// PRD 产品需求文档
type PRD struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Features    []Feature  `json:"features"`
	Constraints []string   `json:"constraints"`
	CreatedAt   time.Time  `json:"created_at"`
}

// Feature 功能特性
type Feature struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
}

// Architecture 架构设计
type Architecture struct {
	TechStack    []string          `json:"tech_stack"`
	Directory    map[string]string `json:"directory"`
	Interfaces   []Interface       `json:"interfaces"`
}

// Interface 接口定义
type Interface struct {
	Name    string   `json:"name"`
	Methods []Method `json:"methods"`
}

// Method 方法定义
type Method struct {
	Name       string   `json:"name"`
	Params     []Param  `json:"params"`
	Returns    []Param  `json:"returns"`
}

// Param 参数定义
type Param struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Task 开发任务
type Task struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	DependsOn   []string `json:"depends_on"`
	Status      string   `json:"status"`
}

// TaskStatus 任务状态
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
	TaskStatusFailed     TaskStatus = "failed"
)

// SPEC 技术规范
type SPEC struct {
	TaskID      string      `json:"task_id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Input       []Param     `json:"input"`
	Output      []Param     `json:"output"`
	Interfaces  []Interface `json:"interfaces"`
	Boundary    []string    `json:"boundary"`
}

// CodeFile 代码文件
type CodeFile struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

// TestReport 测试报告
type TestReport struct {
	Passed      bool     `json:"passed"`
	Total       int      `json:"total"`
	PassedCount int      `json:"passed_count"`
	Failed      int      `json:"failed"`
	Coverage    float64  `json:"coverage"`
	Errors      []string `json:"errors,omitempty"`
}

// ContextSnapshot 上下文快照
type ContextSnapshot struct {
	ID             string    `json:"id"`
	WorkflowID     string    `json:"workflow_id"`
	PRDSummary     string    `json:"prd_summary"`
	CurrentTask    string    `json:"current_task"`
	CompletedTasks []string  `json:"completed_tasks"`
	Decisions      []string  `json:"decisions"`
	CreatedAt      time.Time `json:"created_at"`
}

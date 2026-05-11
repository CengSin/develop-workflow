package activity

import (
	"context"

	"go.temporal.io/sdk/activity"

	"github.com/cengsin/develop-workflow/internal/agent"
	"github.com/cengsin/develop-workflow/pkg/models"
)

var (
	requirementAgent *agent.RequirementAgent
	architectAgent   *agent.ArchitectAgent
	planningAgent    *agent.PlanningAgent
	codingAgent      *agent.CodingAgent
	validationAgent  *agent.ValidationAgent
)

func Init() {
	requirementAgent = agent.NewRequirementAgent()
	architectAgent = agent.NewArchitectAgent()
	planningAgent = agent.NewPlanningAgent()
	codingAgent = agent.NewCodingAgent()
	validationAgent = agent.NewValidationAgent()
}

// RequirementActivity 需求分析Activity
func RequirementActivity(ctx context.Context, req models.DevelopRequest) (*models.PRD, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("RequirementActivity started", "project", req.ProjectName)

	result, err := requirementAgent.Analyze(ctx, req)
	if err != nil {
		return nil, err
	}

	prd := result.(*models.PRD)
	logger.Info("RequirementActivity completed", "features", len(prd.Features))
	return prd, nil
}

// ArchitectureActivity 架构设计Activity
func ArchitectureActivity(ctx context.Context, prd models.PRD) (*models.Architecture, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("ArchitectureActivity started", "prd", prd.Title)

	result, err := architectAgent.Design(ctx, prd)
	if err != nil {
		return nil, err
	}

	arch := result.(*models.Architecture)
	logger.Info("ArchitectureActivity completed", "interfaces", len(arch.Interfaces))
	return arch, nil
}

// PlanningActivity 任务规划Activity
func PlanningActivity(ctx context.Context, prd models.PRD, arch models.Architecture) ([]models.Task, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("PlanningActivity started")

	result, err := planningAgent.Plan(ctx, prd, arch)
	if err != nil {
		return nil, err
	}

	tasks := result.([]models.Task)
	logger.Info("PlanningActivity completed", "tasks", len(tasks))
	return tasks, nil
}

// SpecActivity 技术规范Activity
func SpecActivity(ctx context.Context, task models.Task, prd models.PRD, arch models.Architecture) (*models.SPEC, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("SpecActivity started", "task", task.ID)

	result, err := planningAgent.GenerateSpec(ctx, task, prd, arch)
	if err != nil {
		return nil, err
	}

	spec := result.(*models.SPEC)
	logger.Info("SpecActivity completed", "spec", spec.Title)
	return spec, nil
}

// CodingActivity 代码实现Activity
func CodingActivity(ctx context.Context, spec models.SPEC) ([]models.CodeFile, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("CodingActivity started", "spec", spec.Title)

	result, err := codingAgent.Implement(ctx, spec)
	if err != nil {
		return nil, err
	}

	codeFiles := result.([]models.CodeFile)
	logger.Info("CodingActivity completed", "files", len(codeFiles))
	return codeFiles, nil
}

// ValidationActivity 测试验证Activity
func ValidationActivity(ctx context.Context, codeFiles []models.CodeFile) (*models.TestReport, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("ValidationActivity started")

	result, err := validationAgent.Validate(ctx, codeFiles)
	if err != nil {
		return nil, err
	}

	report := result.(*models.TestReport)
	logger.Info("ValidationActivity completed", "passed", report.Passed)
	return report, nil
}

// DocumentationActivity 文档生成Activity
func DocumentationActivity(ctx context.Context, prd models.PRD, arch models.Architecture, tasks []models.Task) error {
	logger := activity.GetLogger(ctx)
	logger.Info("DocumentationActivity started")

	// TODO: 实现文档生成逻辑
	logger.Info("DocumentationActivity completed")
	return nil
}

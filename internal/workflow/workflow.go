package workflow

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"github.com/cengsin/develop-workflow/internal/activity"
	"github.com/cengsin/develop-workflow/pkg/models"
)

const (
	MaxRetries    = 3
	RetryTimeout  = 5 * time.Minute
	ActivityTimeout = 30 * time.Minute
)

// DevelopWorkflow 开发工作流
func DevelopWorkflow(ctx workflow.Context, req models.DevelopRequest) (*models.DevelopResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting develop workflow", "project", req.ProjectName)

	// 设置重试策略
	retryPolicy := &temporal.RetryPolicy{
		MaximumAttempts:    MaxRetries,
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Minute,
	}

	// 设置Activity选项
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: ActivityTimeout,
		RetryPolicy:         retryPolicy,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// Step 1: 需求分析
	logger.Info("Step 1: Analyzing requirements")
	var prd models.PRD
	err := workflow.ExecuteActivity(ctx, activity.RequirementActivity, req).Get(ctx, &prd)
	if err != nil {
		return nil, fmt.Errorf("requirement activity failed: %w", err)
	}

	// Step 2: 架构设计
	logger.Info("Step 2: Designing architecture")
	var arch models.Architecture
	err = workflow.ExecuteActivity(ctx, activity.ArchitectureActivity, prd).Get(ctx, &arch)
	if err != nil {
		return nil, fmt.Errorf("architecture activity failed: %w", err)
	}

	// Step 3: 任务规划
	logger.Info("Step 3: Planning tasks")
	var tasks []models.Task
	err = workflow.ExecuteActivity(ctx, activity.PlanningActivity, prd, arch).Get(ctx, &tasks)
	if err != nil {
		return nil, fmt.Errorf("planning activity failed: %w", err)
	}

	// Step 4: 逐个任务实现
	logger.Info("Step 4: Implementing tasks", "count", len(tasks))
	var completedTasks []string
	for _, task := range tasks {
		logger.Info("Implementing task", "task_id", task.ID, "task_name", task.Name)

		// 生成技术规范
		var spec models.SPEC
		err = workflow.ExecuteActivity(ctx, activity.SpecActivity, task, prd, arch).Get(ctx, &spec)
		if err != nil {
			return nil, fmt.Errorf("spec activity failed for task %s: %w", task.ID, err)
		}

		// 实现代码
		var codeFiles []models.CodeFile
		err = workflow.ExecuteActivity(ctx, activity.CodingActivity, spec).Get(ctx, &codeFiles)
		if err != nil {
			return nil, fmt.Errorf("coding activity failed for task %s: %w", task.ID, err)
		}

		// 测试验证
		var testReport models.TestReport
		err = workflow.ExecuteActivity(ctx, activity.ValidationActivity, codeFiles).Get(ctx, &testReport)
		if err != nil {
			return nil, fmt.Errorf("validation activity failed for task %s: %w", task.ID, err)
		}

		// 如果测试失败，重试
		if !testReport.Passed {
			logger.Info("Test failed, retrying", "task_id", task.ID)
			err = workflow.ExecuteActivity(ctx, activity.CodingActivity, spec).Get(ctx, &codeFiles)
			if err != nil {
				return nil, fmt.Errorf("coding retry failed for task %s: %w", task.ID, err)
			}

			err = workflow.ExecuteActivity(ctx, activity.ValidationActivity, codeFiles).Get(ctx, &testReport)
			if err != nil {
				return nil, fmt.Errorf("validation retry failed for task %s: %w", task.ID, err)
			}

			if !testReport.Passed {
				return nil, fmt.Errorf("task %s failed after retry", task.ID)
			}
		}

		completedTasks = append(completedTasks, task.ID)
	}

	// Step 5: 生成交付物
	logger.Info("Step 5: Generating deliverables")
	err = workflow.ExecuteActivity(ctx, activity.DocumentationActivity, prd, arch, tasks).Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("documentation activity failed: %w", err)
	}

	logger.Info("Workflow completed successfully")
	return &models.DevelopResult{
		Success:    true,
		Repository: req.ProjectName,
	}, nil
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.temporal.io/sdk/client"

	"github.com/cengsin/develop-workflow/internal/config"
	"github.com/cengsin/develop-workflow/internal/workflow"
	"github.com/cengsin/develop-workflow/pkg/models"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 检查命令行参数
	if len(os.Args) < 3 {
		fmt.Println("Usage: cli <project-name> <description>")
		fmt.Println("Example: cli my-todo-app 'Create a simple TODO application'")
		os.Exit(1)
	}

	projectName := os.Args[1]
	description := os.Args[2]

	// 创建Temporal客户端
	temporalClient, err := client.Dial(client.Options{
		HostPort: cfg.Temporal.Address,
	})
	if err != nil {
		log.Fatalf("Unable to create Temporal client: %v", err)
	}
	defer temporalClient.Close()

	// 创建开发请求
	req := models.DevelopRequest{
		Description: description,
		ProjectName: projectName,
		Language:    "go",
	}

	// 启动工作流
	log.Printf("Starting develop workflow for project: %s", projectName)
	workflowID := fmt.Sprintf("develop-%s", projectName)

	workflowOptions := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: cfg.Temporal.TaskQueue,
	}

	we, err := temporalClient.ExecuteWorkflow(context.Background(), workflowOptions, workflow.DevelopWorkflow, req)
	if err != nil {
		log.Fatalf("Unable to execute workflow: %v", err)
	}

	log.Printf("Started workflow: WorkflowID=%s, RunID=%s", we.GetID(), we.GetRunID())

	// 等待工作流完成
	var result models.DevelopResult
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalf("Workflow failed: %v", err)
	}

	if result.Success {
		fmt.Printf("\n✅ Project '%s' developed successfully!\n", result.Repository)
		fmt.Printf("Repository: %s\n", result.Repository)
	} else {
		fmt.Printf("\n❌ Project development failed: %s\n", result.Error)
		os.Exit(1)
	}
}

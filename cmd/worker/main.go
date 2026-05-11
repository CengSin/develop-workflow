package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/cengsin/develop-workflow/internal/activity"
	"github.com/cengsin/develop-workflow/internal/config"
	"github.com/cengsin/develop-workflow/internal/workflow"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化Activity
	activity.Init()

	// 创建Temporal客户端
	temporalClient, err := client.Dial(client.Options{
		HostPort: cfg.Temporal.Address,
	})
	if err != nil {
		log.Fatalf("Unable to create Temporal client: %v", err)
	}
	defer temporalClient.Close()

	// 创建Worker
	w := worker.New(temporalClient, cfg.Temporal.TaskQueue, worker.Options{})

	// 注册Workflow
	w.RegisterWorkflow(workflow.DevelopWorkflow)

	// 注册Activities
	w.RegisterActivity(activity.RequirementActivity)
	w.RegisterActivity(activity.ArchitectureActivity)
	w.RegisterActivity(activity.PlanningActivity)
	w.RegisterActivity(activity.SpecActivity)
	w.RegisterActivity(activity.CodingActivity)
	w.RegisterActivity(activity.ValidationActivity)
	w.RegisterActivity(activity.DocumentationActivity)

	log.Printf("Starting Worker on task queue: %s", cfg.Temporal.TaskQueue)
	log.Printf("Temporal Server: %s", cfg.Temporal.Address)

	// 启动Worker
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalf("Unable to start Worker: %v", err)
	}
}

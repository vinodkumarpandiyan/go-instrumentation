// workflows.go
package main

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

func SampleWorkflow(ctx workflow.Context, name string) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Workflow started", "name", name)

	// Simulate some work
	workflow.Sleep(ctx, time.Second*2)

	return "Hello, " + name, nil
}

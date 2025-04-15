package main

import (
	"context"
	"time"

	"go.temporal.io/sdk/workflow"
)

func MyWorkflow(ctx workflow.Context, name string) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting workflow", "name", name)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var result string
	err := workflow.ExecuteActivity(ctx, MyActivity, name).Get(ctx, &result)
	if err != nil {
		return "", err
	}
	return result, nil
}

func MyActivity(ctx context.Context, name string) (string, error) {
	return "Hello " + name + " from Temporal activity!", nil
}

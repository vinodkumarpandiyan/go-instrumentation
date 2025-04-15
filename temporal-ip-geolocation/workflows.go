package iplocate

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func GetAddressFromIP(ctx workflow.Context, name string) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			MaximumInterval:    time.Minute,
			BackoffCoefficient: 2,
			MaximumAttempts:    5,
		},
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	var ipActivities *IPActivities
	var ip string

	err := workflow.ExecuteActivity(ctx, ipActivities.GetIP).Get(ctx, &ip)

	if err != nil {
		return "", fmt.Errorf("Failed to get IP: %s", err)
	}

	var location string
	err = workflow.ExecuteActivity(ctx, ipActivities.GetLocationInfo, ip).Get(ctx, &location)

	if err != nil {
		return "", fmt.Errorf("Failed to get location: %s", err)
	}

	return fmt.Sprintf("Hello, %s. Your IP is %s and your location is %s", name, ip, location), nil
}

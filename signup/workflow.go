package simpleworkflow

import (
	"context"
	"time"

	"go.temporal.io/sdk/workflow"
)

func SignupWorkflow(ctx workflow.Context, email string, name string) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var result string
	err := workflow.ExecuteActivity(ctx, SendWelcomeEmail, email, name).Get(ctx, &result)
	if err != nil {
		return "", err
	}

	return result, nil
}

func SendWelcomeEmail(ctx context.Context, email string, name string) (string, error) {
	return "Hello " + name + "!", nil
}

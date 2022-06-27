package simpleworkflow

import (
	"context"
	"fmt"
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

type SendWelcomeEmailInput struct {
	Name  string
	Email string
}

type SendWelcomeEmailResult struct {
	MessageID string
}

func SendWelcomeEmail(ctx context.Context, input SendWelcomeEmailInput) (SendWelcomeEmailResult, error) {
	// Here we'd send the email, for the demo we'll just pretend
	fmt.Printf("To: %s\nFrom: support@\nSubject: Welcome to our mailing list\n\nHi %s,\nWe'll be in touch soon!\n", input.Email, input.Name)

	return SendWelcomeEmailResult{
		MessageID: input.Email + "-1234",
	}, nil
}

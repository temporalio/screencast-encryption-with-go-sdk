package simpleworkflow

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

type SignupInput struct {
	Email string
	Name  string
}

func SignupWorkflow(ctx workflow.Context, input SignupInput) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	err := workflow.ExecuteActivity(ctx, SendWelcomeEmail, SendWelcomeEmailInput{Name: input.Name, Email: input.Email}).Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
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

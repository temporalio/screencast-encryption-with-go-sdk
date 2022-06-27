package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/temporalio/screencast-encryption-with-go-sdk/codec"
	signup "github.com/temporalio/screencast-encryption-with-go-sdk/signup"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/worker"
)

func NewClient(options client.Options) (client.Client, error) {
	keyID := os.Getenv("ENCRYPTION_KEY_ID")
	if keyID == "" {
		return nil, fmt.Errorf("ENCRYPTION_KEY_ID environment variable is required to create dataconverter for client")
	}

	options.DataConverter = codec.NewEncryptionDataConverter(
		converter.GetDefaultDataConverter(),
		codec.CodecOptions{KeyID: keyID},
	)

	return client.Dial(options)
}

func runWorker() {
	c, err := NewClient(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "default", worker.Options{})

	w.RegisterWorkflow(signup.SignupWorkflow)
	w.RegisterActivity(signup.SendWelcomeEmail)

	w.Run(nil)
}

func main() {
	go func() {
		runWorker()
	}()

	c, err := NewClient(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	wf, err := c.ExecuteWorkflow(
		context.Background(),
		client.StartWorkflowOptions{
			TaskQueue: "default",
		},
		signup.SignupWorkflow,
		signup.SignupInput{
			Name:  "John Smith",
			Email: "john@example.com",
		},
	)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	err = wf.Get(context.Background(), nil)
	if err != nil {
		log.Fatalln("Workflow failed", err)
	}

	log.Printf("Workflow completed\n")
}

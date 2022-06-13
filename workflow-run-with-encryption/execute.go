package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/temporalio/screencast-encryption-with-go-sdk/codec"
	"github.com/temporalio/screencast-encryption-with-go-sdk/simpleworkflow"
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

	w.RegisterWorkflow(simpleworkflow.SimpleWorkflow)
	w.RegisterActivity(simpleworkflow.SimpleActivity)

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
		simpleworkflow.SimpleWorkflow,
		"John Smith",
	)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	var result string
	err = wf.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Workflow failed", err)
	}

	log.Printf("Workflow completed. Result: %v\n", result)
}

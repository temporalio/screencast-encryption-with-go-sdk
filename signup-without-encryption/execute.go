package main

import (
	"context"
	"log"

	signup "github.com/temporalio/screencast-encryption-with-go-sdk/signup"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func runWorker() {
	c, err := client.Dial(client.Options{})
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

	c, err := client.Dial(client.Options{})
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
		"John Smith",
		"john@example.com",
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

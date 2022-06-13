package main

import (
	"context"
	"log"

	"github.com/temporalio/screencast-remote-codec-server-go/simpleworkflow"

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

	w.RegisterWorkflow(simpleworkflow.SimpleWorkflow)
	w.RegisterActivity(simpleworkflow.SimpleActivity)

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

// webserver.go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.temporal.io/sdk/client"
)

func startWorkflowHandler(w http.ResponseWriter, r *http.Request) {
	temporalClient, err := client.NewClient(client.Options{})
	if err != nil {
		http.Error(w, "Failed to create Temporal client", 500)
		return
	}
	defer temporalClient.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "sample_workflow_" + fmt.Sprint(time.Now().Unix()),
		TaskQueue: "sample-task-queue",
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Temporal"
	}

	we, err := temporalClient.ExecuteWorkflow(context.Background(), workflowOptions, SampleWorkflow, name)
	if err != nil {
		http.Error(w, "Failed to start workflow: "+err.Error(), 500)
		return
	}

	fmt.Fprintf(w, "Started workflow with ID: %s\n", we.GetID())
}

func main() {
	go startWorker() // optional: start the worker in the same app

	http.HandleFunc("/start", startWorkflowHandler)
	log.Println("Webserver listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

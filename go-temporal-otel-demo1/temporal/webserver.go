// webserver.go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.temporal.io/sdk/client"
)

func startWorkflowHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tr := otel.Tracer("webserver")

	// Start span for the handler
	ctx, span := tr.Start(ctx, "startWorkflowHandler")
	defer span.End()

	c, err := client.Dial(client.Options{})
	if err != nil {
		http.Error(w, "failed to connect to Temporal: "+err.Error(), 500)
		return
	}
	defer c.Close()

	workflowID := "workflow-" + fmt.Sprint(time.Now().Unix())
	opts := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: "sample-task-queue",
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Temporal"
	}

	we, err := c.ExecuteWorkflow(ctx, opts, SampleWorkflow, name)
	if err != nil {
		http.Error(w, "failed to start workflow: "+err.Error(), 500)
		return
	}

	fmt.Fprintf(w, "Started workflow: %s\n", we.GetID())
}

func main() {
	go startWorker() // optional: start the worker in the same app

	// http.HandleFunc("/start", startWorkflowHandler)
	// log.Println("Webserver listening on :8081")
	// log.Fatal(http.ListenAndServe(":8081", nil))

	cleanup := initTracerAuto()
	defer cleanup(context.Background())

	mux := http.NewServeMux()
	mux.HandleFunc("/start", startWorkflowHandler)

	log.Println("Service is running on port 8081")
	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		log.Fatalf("Service B failed: %v", err)
	}

	// r := gin.Default()
	// r.Use(otelgin.Middleware("otel-otlp-go-service"))

	// r.GET("/", func(c *gin.Context) {
	// 	c.String(200, "Hello, World!")
	// })

	// // Run the server
	// r.Run(":8080")
}

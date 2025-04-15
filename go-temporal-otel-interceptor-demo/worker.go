package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/opentelemetry"
	"go.temporal.io/sdk/interceptor"
	"go.temporal.io/sdk/worker"
)

func startWorker(c client.Client) {
	// Set up the OTEL interceptor
	otelInterceptor, err := opentelemetry.NewTracingInterceptor(opentelemetry.TracerOptions{})

	if err != nil {
		log.Fatalf("Could not create NewTracingInterceptor: %v", err)
	}

	w := worker.New(c, "sample-task-queue", worker.Options{
		Interceptors: []interceptor.WorkerInterceptor{otelInterceptor},
	})

	w.RegisterWorkflow(MyWorkflow)
	w.RegisterActivity(MyActivity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalf("unable to start worker: %v", err)
	}
}

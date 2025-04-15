package main

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/opentelemetry"
	"go.temporal.io/sdk/interceptor"
)

func main() {
	ctx := context.Background()

	// --- OpenTelemetry Tracer Setup with OTLP Exporter ---
	exp, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("localhost:4317"),   // Tempo OTLP endpoint
		otlptracegrpc.WithDialOption(grpc.WithBlock()), // optional: wait until connected
	)
	if err != nil {
		log.Fatalf("failed to create OTLP trace exporter: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("temporal-go-app"),
		)),
	)
	otel.SetTracerProvider(tp)
	defer func() { _ = tp.Shutdown(ctx) }()

	// --- Temporal Client with OTEL Interceptor ---
	clientInterceptor, err := opentelemetry.NewTracingInterceptor(opentelemetry.TracerOptions{})

	if err != nil {
		log.Fatalf("Could not create NewTracingInterceptor: %v", err)
	}

	c, err := client.Dial(client.Options{
		HostPort:     client.DefaultHostPort,
		Interceptors: []interceptor.ClientInterceptor{clientInterceptor},
	})
	if err != nil {
		log.Fatalf("Unable to create client: %v", err)
	}
	defer c.Close()

	// Start the worker in background
	go startWorker(c)

	// Start a workflow
	workflowOptions := client.StartWorkflowOptions{
		ID:        "hello-workflow",
		TaskQueue: "sample-task-queue",
	}

	we, err := c.ExecuteWorkflow(ctx, workflowOptions, MyWorkflow, "John")
	if err != nil {
		log.Fatalf("Unable to execute workflow: %v", err)
	}

	log.Printf("Started workflow. WorkflowID: %s RunID: %s\n", we.GetID(), we.GetRunID())

	var result string
	err = we.Get(ctx, &result)
	if err != nil {
		log.Fatalf("Failed to get workflow result: %v", err)
	}
	log.Printf("Workflow result: %s", result)

	// Give traces a bit of time to export
	time.Sleep(2 * time.Second)
}

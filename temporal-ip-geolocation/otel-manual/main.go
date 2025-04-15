package main

import (
	"context"
	"errors"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func initTracer() func() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	res, err := newResource(ctx)
	reportErr(err, "failed to create res")

	conn, err := grpc.DialContext(ctx, "localhost:4317", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	reportErr(err, "failed to create gRPC connection to collector")

	// Set up a trace exporter
	traceExporter, err := newExporter(ctx, conn)
	reportErr(err, "failed to create trace exporter")

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	batchSpanProcessor := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := newTraceProvider(res, batchSpanProcessor)
	otel.SetTracerProvider(tracerProvider)

	return func() {
		// Shutdown will flush any remaining spans and shut down the exporter.
		reportErr(tracerProvider.Shutdown(ctx), "failed to shutdown TracerProvider")
		cancel()
	}
}

func newTraceProvider(res *resource.Resource, bsp sdktrace.SpanProcessor) *sdktrace.TracerProvider {
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	return tracerProvider
}

func newExporter(ctx context.Context, conn *grpc.ClientConn) (*otlptrace.Exporter, error) {
	return otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
}

func newResource(ctx context.Context) (*resource.Resource, error) {
	return resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String("otel-otlp-go-service"),
			attribute.String("application", "otel-otlp-go-app"),
		),
	)
}

func reportErr(err error, message string) {
	if err != nil {
		log.Printf("%s: %v", message, err)
	}
}

func main() {

	log.Printf("Waiting for connection...")

	shutdown := initTracer()

	defer shutdown()

	tracer := otel.Tracer("demo1TracerName")
	ctx := context.Background()

	// work begins
	parentFunction(ctx, tracer)

	// bonus work
	exceptionFunction(ctx, tracer)

	log.Printf("Done!")
}

func parentFunction(ctx context.Context, tracer trace.Tracer) {
	ctx, parentSpan := tracer.Start(
		ctx,
		"parentSpanName",
		trace.WithAttributes(attribute.String("parentAttributeKey1", "parentAttributeValue1")))

	parentSpan.AddEvent("ParentSpan-Event")
	log.Printf("In parent span, before calling a child function.")

	defer parentSpan.End()

	childFunction(ctx, tracer)

	log.Printf("In parent span, after calling a child function. When this function ends, parentSpan will complete.")
}

func childFunction(ctx context.Context, tracer trace.Tracer) {
	ctx, childSpan := tracer.Start(
		ctx,
		"childSpanName",
		trace.WithAttributes(attribute.String("childAttributeKey1", "childAttributeValue1")))

	childSpan.AddEvent("ChildSpan-Event")
	defer childSpan.End()

	log.Printf("In child span, when this function returns, childSpan will complete.")
}

func exceptionFunction(ctx context.Context, tracer trace.Tracer) {
	ctx, exceptionSpan := tracer.Start(
		ctx,
		"exceptionSpanName",
		trace.WithAttributes(attribute.String("exceptionAttributeKey1", "exceptionAttributeValue1")))
	defer exceptionSpan.End()
	log.Printf("Call division function.")
	_, err := divide(10, 0)
	if err != nil {
		exceptionSpan.RecordError(err)
		exceptionSpan.SetStatus(codes.Error, err.Error())
	}
}

func divide(x int, y int) (int, error) {
	if y == 0 {
		return -1, errors.New("division by zero")
	}
	return x / y, nil
}

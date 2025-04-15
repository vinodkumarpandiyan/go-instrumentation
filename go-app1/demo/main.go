package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	header := r.Header

	carrier := propagation.HeaderCarrier{}
	carrier.Set("Traceparent", header.Get("traceparent"))
	fmt.Println(carrier)
	propgator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	parentCtx := propgator.Extract(ctx, carrier)
	ctx, span := otel.Tracer("propagator").Start(parentCtx, "app-b-child-span")
	defer span.End()

	spanCtx := trace.SpanContextFromContext(parentCtx)
	log.Printf("App A Trace ID: %s, Span ID: %s", spanCtx.TraceID(), spanCtx.SpanID())

	fmt.Fprintln(w, "Hello from App A 👋")

}

func main() {
	shutdown := initTracer()
	defer shutdown(context.Background())

	mux := http.NewServeMux()
	// mux.Handle("/hello", otelhttp.NewHandler(http.HandlerFunc(helloHandler), "HelloHandler"))
	mux.Handle("/hello", http.HandlerFunc(helloHandler))

	log.Println("App A listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}

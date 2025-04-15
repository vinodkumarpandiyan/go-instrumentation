package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func callAppAHandler(w http.ResponseWriter, r *http.Request) {
	tr := otel.Tracer("app-b")
	propgator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})

	ctx, span := tr.Start(context.Background(), "go-app2-handler")
	defer span.End()

	// Serialize the context into carrier
	// carrier := propagation.MapCarrier{}
	// propgator.Inject(ctx, carrier)
	// // This carrier is sent accros the process
	// fmt.Println(carrier)

	spanCtx := trace.SpanContextFromContext(ctx)
	log.Printf("App B Trace ID: %s, Span ID: %s", spanCtx.TraceID(), spanCtx.SpanID())

	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	req, _ := http.NewRequestWithContext(ctx, "GET", "http://localhost:8081/hello", nil)

	propgator.Inject(ctx, propagation.HeaderCarrier(req.Header))
	fmt.Println(req.Header)

	res, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to call App A: "+err.Error(), 500)
		return
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	fmt.Fprintf(w, "App B received: %s", string(body))
}

func main() {
	shutdown := initTracer()
	defer shutdown(context.Background())

	mux := http.NewServeMux()
	// mux.Handle("/hello", otelhttp.NewHandler(http.HandlerFunc(callAppAHandler), "Call-App-A"))
	mux.Handle("/hello", http.HandlerFunc(callAppAHandler))

	log.Println("App B listening on :8082")
	log.Fatal(http.ListenAndServe(":8082", mux))
}

## Start the temporal using below commands
temporal server start-dev --log-level=never --ui-port 8080 --db-filename=temporal.db

## For Manual instrumentation

### To start the go service
cd temporal

go run webserver.go worker.go workflows.go otel.go

<b> Note</b>: service exposed in port 8081

### To start the temporal worker
Access http://localhost:8081/start

### View traces in grafana explorer

http://localhost:3000

## For Auto Instrumentation

## To start go service

go run otel-auto/main.go

Note: service starts in 8080 port

### To Access Go Service

http://localhost:8080
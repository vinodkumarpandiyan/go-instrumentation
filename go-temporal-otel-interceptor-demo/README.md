# Go - Temporal - Otel - Interceptor 
<b>Note:</b> This example demonstrate auto instrumented traces through interceptors

### Start tempo using docker

### Start temporal
```bash
temporal server start-dev --log-level=never --ui-port 8080 --db-filename=temporal.db
```

### Run the go application
```bash 
go run workflows.go worker.go main.go
```

### Traces in Grafana
http://localhost:3000
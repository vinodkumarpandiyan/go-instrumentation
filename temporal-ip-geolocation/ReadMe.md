## to build the dependency 
go mod init

## to update the dependency
go mod tidy

# Otel manual instrumentation
Run the below command to generate traces ( manual instrumentation)
go run otel-manual/main.go

view traces in grafana
http://localhost:3000
 
# Otel auto instrumentation
Run the below command to generate traces ( auto instrumentation)
go run otel-auto/main.go

call the below url to call the api
http://localhost:8080

view traces in grafana
http://localhost:3000

# webservice demo1

go run webservice-demo1/main.go

access the below url in browser
http://localhost:8081/a
http://localhost:8081/b
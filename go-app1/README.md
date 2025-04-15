# Run the go-app1 , it runs in 8081 port
go run demo/main.go demo/otel.go

/hello api is called from go-app2 and trace is exported to tempo.

start go-app2 which runs on 8082 port

# access it in browser
http://localhost:8082/hello


# check the traces in grafana
http://localhost:3000
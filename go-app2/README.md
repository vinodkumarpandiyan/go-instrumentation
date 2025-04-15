# Run the go-app2 , it runs in 8082 port
go run demo/main.go demo/otel.go

/hello api is called from go-app2 and it calls /hello api from go-app1 and trace is exported to tempo.

start go-app1 which runs on 8081 port

# access it in browser
http://localhost:8082/hello


# check the traces in grafana
http://localhost:3000
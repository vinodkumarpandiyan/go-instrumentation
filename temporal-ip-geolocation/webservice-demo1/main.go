package main

import (
	"fmt"
	"log"
	"net/http"
)

// Service A Handler
func serviceAHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Service A!")
}

// Service B Handler
func serviceBHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Service B!")
}

// Run Service A on port 8081
func runServiceA() {
	mux := http.NewServeMux()
	mux.HandleFunc("/a", serviceAHandler)

	log.Println("Service A running on port 8081")
	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		log.Fatalf("Service A failed: %v", err)
	}
}

// Run Service B on port 8082
func runServiceB() {
	mux := http.NewServeMux()
	mux.HandleFunc("/b", serviceBHandler)

	log.Println("Service B running on port 8082")
	err := http.ListenAndServe(":8082", mux)
	if err != nil {
		log.Fatalf("Service B failed: %v", err)
	}
}

func main() {
	// Run services concurrently
	go runServiceA()
	go runServiceB()

	// Block forever
	select {}
}

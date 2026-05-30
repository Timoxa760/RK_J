package main

import (
	"fmt"
	"net/http"
	"os"

	root "backend_project/internal"
	reporting "backend_project/services/reporting/reporting-service/internal"
)

func main() {
	port := "8010"
	serviceName := "reporting-service"
	demoMode := getenv("DEMO_MODE", "true") == "true"

	r := root.NewRouter()
	reporting.New().Register(r)

	fmt.Printf("Service %s started on port %s (demo=%v)...\n", serviceName, port, demoMode)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		panic(err)
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

package main

import (
	"fmt"
	"net/http"
	"os"

	root "backend_project/internal"
	cred "backend_project/services/finance-core/credit-service/internal"
)

func main() {
	port := "8009"
	serviceName := "credit-service"
	demoMode := getenv("DEMO_MODE", "true") == "true"

	r := root.NewRouter()
	cred.New(demoMode).Register(r)

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

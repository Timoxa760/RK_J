package main

import (
	"fmt"
	"net/http"

	"backend_project/internal"
)

func main() {
	port := "8102"
	serviceName := "social-service"

	r := internal.NewRouter()

	fmt.Printf("Service %s started on port %s...\n", serviceName, port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		panic(err)
	}
}

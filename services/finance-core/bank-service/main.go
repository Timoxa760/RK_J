package main

import (
	"fmt"
	"net/http"

	root "backend_project/internal"
	bank "backend_project/services/finance-core/bank-service/internal"
)

func main() {
	port := "8011"
	serviceName := "bank-service"

	r := root.NewRouter()
	bank.New().Register(r)

	fmt.Printf("Service %s started on port %s...\n", serviceName, port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		panic(err)
	}
}

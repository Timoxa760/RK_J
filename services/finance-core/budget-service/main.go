package main

import (
	"fmt"
	"net/http"

	root "backend_project/internal"
	budget "backend_project/services/finance-core/budget-service/internal"
)

func main() {
	port := "8005"
	serviceName := "budget-service"

	r := root.NewRouter()
	budget.New().Register(r)

	fmt.Printf("Service %s started on port %s...\n", serviceName, port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		panic(err)
	}
}

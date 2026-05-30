package main

import (
	"fmt"
	"net/http"

	root "backend_project/internal"
	category "backend_project/services/finance-core/category-service/internal"
)

func main() {
	port := "8004"
	serviceName := "category-service"

	r := root.NewRouter()
	category.New().Register(r)

	fmt.Printf("Service %s started on port %s...\n", serviceName, port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		panic(err)
	}
}

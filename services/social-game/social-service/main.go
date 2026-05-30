package main

import (
	"fmt"
	"net/http"

	root "backend_project/internal"
	social "backend_project/services/social-game/social-service/internal"
)

func main() {
	port := "8102"
	serviceName := "social-service"

	r := root.NewRouter()
	social.New().Register(r)

	fmt.Printf("Service %s started on port %s...\n", serviceName, port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		panic(err)
	}
}

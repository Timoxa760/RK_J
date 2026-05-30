package main

import (
	"fmt"
	"net/http"
	"os"

	root "backend_project/internal"
	"backend_project/internal/creditstore"
	"backend_project/internal/llm"
	"backend_project/internal/profile"
	cred "backend_project/services/finance-core/credit-service/internal"
)

func main() {
	port := "8009"
	serviceName := "credit-service"

	creditStore, err := creditstore.NewFileStore(os.Getenv("CREDIT_DATA_DIR"))
	if err != nil {
		panic(err)
	}
	profileStore, err := profile.NewFileStore(os.Getenv("PROFILE_DATA_DIR"))
	if err != nil {
		panic(err)
	}
	llmClient := llm.NewFromEnv()

	r := root.NewRouter()
	cred.NewHandler(creditStore, profileStore, llmClient).Register(r)

	fmt.Printf("Service %s started on port %s...\n", serviceName, port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		panic(err)
	}
}

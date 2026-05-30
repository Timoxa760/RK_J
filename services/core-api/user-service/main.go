package main

import (
	"fmt"
	"net/http"
	"os"

	"backend_project/internal"
	"backend_project/internal/profile"
	"backend_project/services/core-api/user-service/auth"
	profilehandler "backend_project/services/core-api/user-service/profile"
	"backend_project/services/core-api/user-service/providers"
)

func main() {
	port := "8001"
	serviceName := "user-service"
	demoMode := os.Getenv("DEMO_MODE") == "true"

	r := internal.NewRouter()

	registerHandler := auth.NewRegisterHandler(demoMode)
	loginHandler := auth.NewLoginHandler(demoMode)
	connectHandler := providers.NewConnectHandler(demoMode)

	r.Post("/api/v1/auth/register", registerHandler.ServeHTTP)
	r.Post("/api/v1/auth/login", loginHandler.ServeHTTP)
	r.Post("/api/v1/providers/connect", connectHandler.ServeHTTP)

	profileStore, err := profile.NewFileStore(os.Getenv("PROFILE_DATA_DIR"))
	if err != nil {
		panic(err)
	}
	profilehandler.NewHandler(profileStore).Register(r)

	fmt.Printf("Service %s started on port %s...\n", serviceName, port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		panic(err)
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	root "backend_project/internal"
	"backend_project/internal/creditstore"
	"backend_project/internal/llm"
	"backend_project/internal/postgres"
	"backend_project/internal/profile"
	cred "backend_project/services/finance-core/credit-service/internal"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	port := "8009"
	serviceName := "credit-service"

	ctx := context.Background()
	pool := connectPostgres(ctx)

	creditStore, err := creditstore.NewFileStore(os.Getenv("CREDIT_DATA_DIR"))
	if err != nil {
		panic(err)
	}
	profileStore := buildProfileStore(pool)
	llmClient := llm.NewFromEnv()

	r := root.NewRouter()
	cred.NewHandler(creditStore, profileStore, llmClient).Register(r)

	fmt.Printf("Service %s started on port %s...\n", serviceName, port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		panic(err)
	}
}

func connectPostgres(ctx context.Context) *pgxpool.Pool {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil
	}
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Printf("pgxpool: %v (profile will use file store only)", err)
		return nil
	}
	if !postgres.Ping(ctx, pool) {
		log.Printf("postgres unavailable (profile will use file store only)")
		pool.Close()
		return nil
	}
	return pool
}

func buildProfileStore(pool *pgxpool.Pool) profile.Store {
	fileStore, err := profile.NewFileStore(os.Getenv("PROFILE_DATA_DIR"))
	if err != nil {
		panic(err)
	}
	demoMode := os.Getenv("DEMO_MODE") == "true"
	if pool != nil && !demoMode {
		return profile.NewDualStore(profile.NewPGStore(pool), fileStore)
	}
	return fileStore
}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"backend_project/internal"
	"backend_project/internal/otp"
	"backend_project/internal/postgres"
	"backend_project/internal/profile"
	"backend_project/internal/userstore"
	"backend_project/services/core-api/user-service/auth"
	profilehandler "backend_project/services/core-api/user-service/profile"
	"backend_project/services/core-api/user-service/providers"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {
	port := "8001"
	serviceName := "user-service"
	demoMode := os.Getenv("DEMO_MODE") == "true"

	ctx := context.Background()
	pool := connectPostgres(ctx)
	otpStore := connectOTP()

	var users userstore.Store
	if pool != nil {
		users = userstore.NewPostgres(pool)
	} else {
		log.Printf("WARN: no DATABASE_URL — using in-memory users (dev only)")
		users = userstore.NewMemory()
	}

	deps := auth.NewDeps(users, otpStore)

	r := internal.NewRouter()
	r.Post("/api/v1/auth/register", auth.NewRegisterHandler(deps).ServeHTTP)
	r.Post("/api/v1/auth/login", auth.NewLoginHandler(deps).ServeHTTP)
	r.Post("/api/v1/auth/password/forgot", auth.NewForgotPasswordHandler(deps).ServeHTTP)
	r.Post("/api/v1/auth/password/reset", auth.NewResetPasswordHandler(deps).ServeHTTP)
	r.Post("/api/v1/providers/connect", providers.NewConnectHandler(demoMode).ServeHTTP)

	profileStore := buildProfileStore(pool)
	profilehandler.NewHandler(profileStore).Register(r)

	fmt.Printf("Service %s started on port %s (auth=phone+password)...\n", serviceName, port)
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
		log.Fatalf("pgxpool: %v", err)
	}
	if !postgres.Ping(ctx, pool) {
		log.Fatalf("postgres unavailable")
	}
	return pool
}

func connectOTP() otp.Store {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		return otp.NewMemoryStore()
	}
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("redis url: %v", err)
	}
	client := redis.NewClient(opts)
	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Printf("redis unavailable, using memory OTP store: %v", err)
		return otp.NewMemoryStore()
	}
	return otp.NewRedisStore(client)
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

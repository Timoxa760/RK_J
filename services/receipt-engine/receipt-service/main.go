package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	root "backend_project/internal"
	"backend_project/internal/expensestore"
	"backend_project/services/receipt-engine/receipt-service/internal/dashboard"
	svckafka "backend_project/services/receipt-engine/receipt-service/internal/kafka"
	svc "backend_project/services/receipt-engine/receipt-service/internal"
)

func main() {
	port := "8002"
	serviceName := "receipt-service"

	r := root.NewRouter()

	brokers := strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ",")
	databaseURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/moneymind?sslmode=disable")

	bg := context.Background()
	pool, err := connectPostgres(bg, databaseURL, 30)
	if err != nil {
		log.Fatalf("pgxpool: %v", err)
	}
	defer pool.Close()

	repo := svc.NewReceiptRepo(pool)

	dedup := svc.NewDedup(repo, 7*24*time.Hour)
	defer dedup.Stop()

	producer := svckafka.NewProducer(brokers, "receipt.parsed")
	defer producer.Close()

	consumer := svckafka.NewConsumer(brokers, "receipt.raw", "receipt-service", func(ctx context.Context, receipt svc.RawReceipt) error {
		if err := svc.ValidateReceipt(&receipt); err != nil {
			log.Printf("validation failed: %v", err)
			return nil
		}

		dup, err := dedup.IsDuplicate(ctx, receipt.ID, receipt.Provider, receipt.UserID)
		if err != nil {
			return fmt.Errorf("dedup: %w", err)
		}
		if dup {
			log.Printf("duplicate receipt: %s", receipt.ID)
			return nil
		}

		if err := repo.Insert(ctx, &receipt); err != nil {
			return fmt.Errorf("save: %w", err)
		}

		if err := producer.SendParsed(ctx, receipt); err != nil {
			log.Printf("kafka producer: %v", err)
		}

		log.Printf("processed receipt: id=%s provider=%s store=%s total=%.2f",
			receipt.ID, receipt.Provider, receipt.Store, receipt.Total)
		return nil
	})
	defer consumer.Close()

	ctx, cancel := context.WithCancel(bg)
	defer cancel()

	go func() {
		if err := consumer.Start(ctx); err != nil && err != context.Canceled {
			log.Printf("consumer stopped: %v", err)
		}
	}()

	demoMode := getEnv("DEMO_MODE", "true") == "true"
	dash := dashboard.New(demoMode, getEnv("JWT_SECRET", ""))
	dash.SetPostgres(pool)
	if fileStore, err := expensestore.NewFileStore(expensestore.DefaultPath()); err != nil {
		log.Printf("expense file store: %v (dashboard will skip file fallback)", err)
	} else {
		dash.SetExpenseFile(fileStore)
	}
	if !demoMode {
		chHost := getEnv("CLICKHOUSE_HOST", "clickhouse")
		chUser := getEnv("CLICKHOUSE_USER", "clickhouse_user")
		chPass := getEnv("CLICKHOUSE_PASSWORD", "clickhouse_password")
		chDB := getEnv("CLICKHOUSE_DB", "default")
		if err := dash.ConnectClickHouse(ctx, chHost, chUser, chPass, chDB); err != nil {
			log.Printf("clickhouse: %v (fallback to mock data)", err)
		}
	}
	dash.Register(r)

	log.Printf("Service %s started on port %s (demo=%v)...", serviceName, port, demoMode)
	server := &http.Server{Addr: ":" + port, Handler: r}

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		cancel()
		server.Shutdown(context.Background())
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// connectPostgres подключается к PostgreSQL с повторными попытками при старте стека.
func connectPostgres(ctx context.Context, url string, attempts int) (*pgxpool.Pool, error) {
	var lastErr error
	for i := 0; i < attempts; i++ {
		pool, err := pgxpool.New(ctx, url)
		if err != nil {
			lastErr = err
			time.Sleep(time.Second)
			continue
		}
		if err := pool.Ping(ctx); err != nil {
			pool.Close()
			lastErr = err
			time.Sleep(time.Second)
			continue
		}
		return pool, nil
	}
	return nil, lastErr
}

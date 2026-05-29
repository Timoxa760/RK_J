package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	root "backend_project/internal"
	"backend_project/services/money-intelligence/ai-processor/internal"
	"backend_project/services/money-intelligence/ai-processor/internal/categorizer"
	"backend_project/services/money-intelligence/ai-processor/internal/clickhouse"
	"backend_project/services/money-intelligence/ai-processor/internal/manual"
	svckafka "backend_project/services/money-intelligence/ai-processor/internal/kafka"
)

func main() {
	port := "8100"
	demoMode := os.Getenv("DEMO_MODE") == "true"

	brokers := strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ",")
	chHost := getEnv("CLICKHOUSE_HOST", "localhost")
	chUser := getEnv("CLICKHOUSE_USER", "clickhouse_user")
	chPass := getEnv("CLICKHOUSE_PASSWORD", "clickhouse_password")
	chDB := getEnv("CLICKHOUSE_DB", "default")
	databaseURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/moneymind?sslmode=disable")

	cat := categorizer.NewDict()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var chWriter *clickhouse.Writer
	if !demoMode {
		var err error
		chWriter, err = clickhouse.NewWriter(ctx, chHost, chUser, chPass, chDB)
		if err != nil {
			log.Printf("clickhouse: %v (manual expenses will skip ClickHouse)", err)
		} else {
			defer chWriter.Close()
		}
	}

	pgPool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("pgxpool: %v", err)
	}
	defer pgPool.Close()

	consumer := svckafka.NewConsumer(brokers, "receipt.parsed", "ai-processor", func(ctx context.Context, receipt internal.RawReceipt) error {
		categorized := cat.Categorize(receipt.Items)
		cr := &internal.CategorizedReceipt{
			UserID:   receipt.UserID,
			Store:    receipt.Store,
			Provider: receipt.Provider,
			Date:     receipt.Date,
			Items:    categorized,
		}

		if chWriter != nil {
			if err := chWriter.InsertReceipt(ctx, cr); err != nil {
				return err
			}
		}

		log.Printf("processed receipt: id=%s provider=%s store=%s items=%d",
			receipt.ID, receipt.Provider, receipt.Store, len(categorized))
		return nil
	})
	defer consumer.Close()

	go func() {
		if err := consumer.Start(ctx); err != nil && err != context.Canceled {
			log.Printf("consumer stopped: %v", err)
		}
	}()

	var manualRepo *manual.Repo
	if chWriter != nil {
		manualRepo = manual.NewRepo(pgPool, chWriter.Conn())
	} else {
		manualRepo = manual.NewRepo(pgPool, nil)
	}
	manualHandler := manual.NewHandler(manualRepo)

	r := root.NewRouter()
	r.Post("/api/v1/expenses/manual", manualHandler.Create)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("ai-processor HTTP server started on port %s (demo=%v)", port, demoMode)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	log.Printf("ai-processor started (demo_mode=%v)", demoMode)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	srv.Shutdown(shutdownCtx)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

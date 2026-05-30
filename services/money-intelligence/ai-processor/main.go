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
	"backend_project/internal/onlysq"
	"backend_project/services/money-intelligence/ai-processor/internal"
	"backend_project/services/money-intelligence/ai-processor/internal/categorizer"
	"backend_project/services/money-intelligence/ai-processor/internal/clickhouse"
	"backend_project/services/money-intelligence/ai-processor/internal/expense"
	"backend_project/services/money-intelligence/ai-processor/internal/manual"
	svckafka "backend_project/services/money-intelligence/ai-processor/internal/kafka"
	"backend_project/services/money-intelligence/ai-processor/internal/voice"
	"backend_project/services/money-intelligence/ai-processor/internal/whisper"
)

func main() {
	port := "8100"
	demoMode := getEnv("DEMO_MODE", "true") == "true"

	brokers := strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ",")
	chHost := getEnv("CLICKHOUSE_HOST", "localhost")
	chUser := getEnv("CLICKHOUSE_USER", "clickhouse_user")
	chPass := getEnv("CLICKHOUSE_PASSWORD", "clickhouse_password")
	chDB := getEnv("CLICKHOUSE_DB", "default")
	databaseURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/moneymind?sslmode=disable")

	llmClient := onlysq.NewClient(onlysqBaseURL(), getEnv("ONLYSQ_API_KEY", ""), getEnv("ONLYSQ_MODEL", ""))
	expenseParser := expense.NewParser(llmClient)
	whisperClient := whisper.NewClient(getEnv("WHISPER_URL", ""), getEnv("WHISPER_API_KEY", ""), getEnv("WHISPER_MODEL", ""))

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

	var pgPool *pgxpool.Pool
	if !demoMode {
		var err error
		pgPool, err = pgxpool.New(ctx, databaseURL)
		if err != nil {
			log.Fatalf("pgxpool: %v", err)
		}
		defer pgPool.Close()
	}

	var consumer *svckafka.Consumer
	if !demoMode {
		consumer = svckafka.NewConsumer(brokers, "receipt.parsed", "ai-processor", func(ctx context.Context, receipt internal.RawReceipt) error {
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
	}

	r := root.NewRouter()
	var proc *manual.Processor
	if demoMode {
		var demoHandler *manual.DemoHandler
		demoHandler, proc = manual.NewDemoHandler(expenseParser)
		r.Post("/api/v1/expenses/manual", demoHandler.Create)
	} else {
		var manualRepo *manual.Repo
		if chWriter != nil {
			manualRepo = manual.NewRepo(pgPool, chWriter.Conn())
		} else {
			manualRepo = manual.NewRepo(pgPool, nil)
		}
		var handler *manual.Handler
		handler, proc = manual.NewHandler(manualRepo, expenseParser)
		r.Post("/api/v1/expenses/manual", handler.Create)
	}
	r.Post("/api/v1/expenses/voice", voice.NewHandler(whisperClient, proc).Create)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 120 * time.Second,
	}

	go func() {
		log.Printf("ai-processor HTTP server started on port %s (demo=%v onlysq=%v whisper=%v)",
			port, demoMode, llmClient.Enabled(), whisperClient.Enabled())
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

func onlysqBaseURL() string {
	if u := os.Getenv("ONLYSQ_BASE_URL"); u != "" {
		return u
	}
	return os.Getenv("ONLYSQ_URL")
}

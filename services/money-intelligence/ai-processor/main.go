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

	root 	"backend_project/internal"
	"backend_project/internal/creditstore"
	"backend_project/internal/expensestore"
	iadvisor "backend_project/internal/advisor"
	"backend_project/internal/llm"
	"backend_project/internal/postgres"
	"backend_project/internal/profile"
	"backend_project/services/money-intelligence/ai-processor/internal"
	"backend_project/services/money-intelligence/ai-processor/internal/advisor"
	"backend_project/services/money-intelligence/ai-processor/internal/categorizer"
	"backend_project/services/money-intelligence/ai-processor/internal/clickhouse"
	"backend_project/services/money-intelligence/ai-processor/internal/expense"
	"backend_project/services/money-intelligence/ai-processor/internal/manual"
	"backend_project/services/money-intelligence/ai-processor/internal/onboarding"
	"backend_project/services/money-intelligence/ai-processor/internal/receipt"
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

	llmClient := llm.NewFromEnv()
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
			log.Printf("pgxpool: %v (expenses will use file store)", err)
		} else {
			defer pgPool.Close()
		}
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
	fileStore, err := expensestore.NewFileStore(expensestore.DefaultPath())
	if err != nil {
		log.Fatalf("expense file store: %v", err)
	}

	var proc *manual.Processor
	if demoMode {
		var demoHandler *manual.DemoHandler
		demoHandler, proc = manual.NewDemoHandler(expenseParser)
		r.Post("/api/v1/expenses/manual", demoHandler.Create)
	} else {
		var pgStorage manual.Storage
		if pgPool != nil {
			if postgres.Ping(ctx, pgPool) {
				if err := postgres.EnsureManualExpenses(ctx, pgPool); err != nil {
					log.Printf("postgres schema: %v (expenses will use file store)", err)
				} else {
					var manualRepo *manual.Repo
					if chWriter != nil {
						manualRepo = manual.NewRepo(pgPool, chWriter.Conn())
					} else {
						manualRepo = manual.NewRepo(pgPool, nil)
					}
					pgStorage = manual.NewRepoStorage(manualRepo)
				}
			} else {
				log.Printf("postgres unavailable (expenses will use file store)")
			}
		}
		store := manual.NewFallbackStorage(pgStorage, fileStore)
		var handler *manual.Handler
		handler, proc = manual.NewHandlerWithStorage(store, expenseParser)
		r.Post("/api/v1/expenses/manual", handler.Create)
	}
	r.Post("/api/v1/expenses/voice", voice.NewHandler(whisperClient, proc).Create)

	receiptHandler := receipt.NewHandler(whisperClient, proc)
	r.Post("/api/v1/receipt/manual", receiptHandler.ManualCreate)
	r.Post("/api/v1/receipt/voice", receiptHandler.VoiceCreate)
	r.Post("/api/v1/receipt/from-text", receiptHandler.TextCreate)

	voiceHandler := voice.NewHandler(whisperClient, proc)
	r.Post("/api/v1/voice/transcribe", voiceHandler.Transcribe)

	onboardingHandler := onboarding.NewHandler(llmClient)
	r.Post("/api/v1/onboarding/parse", onboardingHandler.Parse)

	profileStore := buildProfileStore(pgPool, demoMode)
	creditStore, err := creditstore.NewFileStore(getEnv("CREDIT_DATA_DIR", ""))
	if err != nil {
		log.Fatalf("credit store: %v", err)
	}

	var spending iadvisor.SpendingProvider = iadvisor.NewPGSpendingProvider(pgPool, fileStore)
	var chatStore *iadvisor.ChatStore
	if pgPool != nil && postgres.Ping(ctx, pgPool) {
		chatStore = iadvisor.NewChatStore(pgPool)
		if err := chatStore.EnsureSchema(ctx); err != nil {
			log.Printf("advisor chat schema: %v", err)
		}
	}

	advisor.NewHandler(profileStore, creditStore, spending, chatStore, llmClient).Register(r)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 120 * time.Second,
	}

	go func() {
		log.Printf("ai-processor HTTP server started on port %s (demo=%v gemini=%v whisper=%v)",
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

func buildProfileStore(pool *pgxpool.Pool, demoMode bool) profile.Store {
	fileStore, err := profile.NewFileStore(getEnv("PROFILE_DATA_DIR", ""))
	if err != nil {
		log.Fatalf("profile store: %v", err)
	}
	if pool != nil && !demoMode {
		return profile.NewDualStore(profile.NewPGStore(pool), fileStore)
	}
	return fileStore
}

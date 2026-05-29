package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"backend_project/internal"
	scrap "backend_project/services/receipt-engine/scraper-service/internal"
	"backend_project/services/receipt-engine/scraper-service/internal/kafka"
	"backend_project/services/receipt-engine/scraper-service/internal/provider"
	"backend_project/services/receipt-engine/scraper-service/internal/providers/email"
	"backend_project/services/receipt-engine/scraper-service/internal/providers/fns"
	"backend_project/services/receipt-engine/scraper-service/internal/providers/fns_mco"
	"backend_project/services/receipt-engine/scraper-service/internal/providers/magnit"
	"backend_project/services/receipt-engine/scraper-service/internal/providers/mock"
	"backend_project/services/receipt-engine/scraper-service/internal/providers/x5club"
	"backend_project/services/receipt-engine/scraper-service/internal/scheduler"
)

func main() {
	port := "8003"
	serviceName := "scraper-service"
	demoMode := os.Getenv("DEMO_MODE") == "true"

	brokersEnv := os.Getenv("KAFKA_BROKERS")
	var brokers []string
	if brokersEnv == "" {
		brokers = []string{"localhost:9092"}
	} else {
		brokers = strings.Split(brokersEnv, ",")
	}
	producer := kafka.NewProducer(brokers, "receipt.raw")
	defer producer.Close()

	sched := scheduler.New()
	if demoMode {
		sched.Add(mock.New("mock", "Demo Store"), provider.TypeAPI)
	}
	sched.OnSync(func(ctx context.Context, name string, receipts []scrap.RawReceipt) error {
		return producer.SendBatch(ctx, receipts)
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go sched.Start(ctx)

	r := internal.NewRouter()

	oauthMgr := email.NewOAuthManager()
	imapCli := email.NewIMAPClient(oauthMgr)
	emailParser := email.NewParser()
	fnsHandler := fns.NewHandler(demoMode)
	mcoProvider := fns_mco.NewProvider(os.Getenv("RUCAPTCHA_KEY"), os.Getenv("MCO_TOKEN_DIR"))
	x5Client := x5club.NewClient(demoMode)
	magnitClient := magnit.NewClient(demoMode)
	magnitINN := magnit.NewINNDetector()
	magnitMapper := magnit.NewMapper(magnitINN)

	r.Post("/api/v1/fns/ticket", fnsHandler.ServeHTTP)

	r.Post("/api/v1/fns/mco/auth", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Phone string `json:"phone"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
			return
		}
		if body.Phone == "" {
			http.Error(w, `{"error":"phone required"}`, http.StatusBadRequest)
			return
		}

		if err := mcoProvider.StartAuth(body.Phone); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "SMS code sent to phone",
		})
	})

	r.Post("/api/v1/fns/mco/auth/verify", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Phone string `json:"phone"`
			Code  string `json:"code"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
			return
		}
		if body.Phone == "" || body.Code == "" {
			http.Error(w, `{"error":"phone and code required"}`, http.StatusBadRequest)
			return
		}

		if err := mcoProvider.VerifyAuth(body.Phone, body.Code); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "authorization successful",
		})
	})

	r.Post("/api/v1/fns/mco/sync", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Phone string `json:"phone"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
			return
		}
		if body.Phone == "" {
			http.Error(w, `{"error":"phone required"}`, http.StatusBadRequest)
			return
		}

		receipts, err := mcoProvider.SyncReceipts(r.Context(), body.Phone)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		for i := range receipts {
			if err := producer.Send(r.Context(), receipts[i]); err != nil {
				http.Error(w, fmt.Sprintf(`{"error":"kafka: %s"}`, err.Error()), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":  true,
			"count":    len(receipts),
			"provider": "fns_mco",
			"phone":    body.Phone,
		})
	})

	r.Post("/api/v1/fns/qr", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			QR     string `json:"qr"`
			UserID string `json:"user_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
			return
		}
		if body.QR == "" {
			http.Error(w, `{"error":"qr required"}`, http.StatusBadRequest)
			return
		}

		fn, fd, fp, err := fns.ParseQRString(body.QR)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
			return
		}

		receipt, err := fnsHandler.CheckTicket(r.Context(), fn, fd, fp)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		receipt.UserID = body.UserID

		if err := producer.Send(r.Context(), *receipt); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"kafka: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"receipt": receipt,
		})
	})

	r.Get("/api/v1/auth/oauth/{provider}", func(w http.ResponseWriter, r *http.Request) {
		provider := chi.URLParam(r, "provider")
		state := r.URL.Query().Get("state")
		url := oauthMgr.GetAuthURL(email.Provider(provider), state)
		if url == "" {
			http.Error(w, "unknown provider", http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, url, http.StatusFound)
	})

	callbackHandler := email.NewOAuthCallbackHandler(oauthMgr)
	r.Get("/api/v1/auth/oauth/callback", callbackHandler.ServeHTTP)

	r.Post("/api/v1/x5club/send-code", func(w http.ResponseWriter, r *http.Request) {
		phone := r.URL.Query().Get("phone")
		if phone == "" {
			http.Error(w, `{"error":"phone required"}`, http.StatusBadRequest)
			return
		}
		if err := x5Client.SendCode(phone); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"send-code: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success":true,"message":"code sent"}`))
	})

	r.Post("/api/v1/x5club/sync", func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		phone := r.URL.Query().Get("phone")
		code := r.URL.Query().Get("code")

		if userID == "" || phone == "" || code == "" {
			if !demoMode {
				http.Error(w, `{"error":"user_id, phone and code required"}`, http.StatusBadRequest)
				return
			}
		}

		if err := x5Client.Login(phone, code); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"login: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		mapper := x5club.NewMapper()
		pool := x5club.NewPool(x5Client, 3)
		items, err := pool.FetchAll(r.Context(), 1, 20)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"history: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		receipts := mapper.ToRawReceipts(items)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user_id":  userID,
			"provider": "x5club",
			"receipts": receipts,
			"count":    len(receipts),
		})
	})

	r.Post("/api/v1/magnit/send-code", func(w http.ResponseWriter, r *http.Request) {
		phone := r.URL.Query().Get("phone")
		if phone == "" {
			http.Error(w, `{"error":"phone required"}`, http.StatusBadRequest)
			return
		}
		if err := magnitClient.SendCode(phone); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"send-code: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success":true,"message":"code sent"}`))
	})

	r.Post("/api/v1/magnit/sync", func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		phone := r.URL.Query().Get("phone")
		code := r.URL.Query().Get("code")

		if userID == "" || phone == "" || code == "" {
			if !demoMode {
				http.Error(w, `{"error":"user_id, phone and code required"}`, http.StatusBadRequest)
				return
			}
		}

		if err := magnitClient.Login(phone, code); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"login: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		var allReceipts []interface{}
		page := 1
		for {
			items, pages, err := magnitClient.GetReceipts(page)
			if err != nil {
				http.Error(w, fmt.Sprintf(`{"error":"receipts: %s"}`, err.Error()), http.StatusInternalServerError)
				return
			}
			for _, r := range magnitMapper.ToRawReceiptsWithINN(items) {
				allReceipts = append(allReceipts, r)
			}
			if page >= pages {
				break
			}
			page++
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user_id":  userID,
			"provider": "magnit",
			"receipts": allReceipts,
			"count":    len(allReceipts),
		})
	})

	r.Post("/api/v1/email/receipts", func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		provider := r.URL.Query().Get("provider")

		htmlList, err := imapCli.FetchReceipts(r.Context(), email.Provider(provider), userID)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		var receipts []interface{}
		for _, html := range htmlList {
			receipt := emailParser.ParseReceiptHTML(html)
			if receipt != nil {
				receipts = append(receipts, receipt)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"receipts": receipts,
			"count":    len(receipts),
		})
	})

	srv := &http.Server{Addr: ":" + port, Handler: r}

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		cancel()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		srv.Shutdown(shutdownCtx)
	}()

	fmt.Printf("Service %s started on port %s...\n", serviceName, port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server: %v", err)
	}
}

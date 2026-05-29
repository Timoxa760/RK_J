package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"backend_project/internal"
	"backend_project/services/receipt-engine/scraper-service/internal/providers/email"
	"backend_project/services/receipt-engine/scraper-service/internal/providers/fns"
	"backend_project/services/receipt-engine/scraper-service/internal/providers/magnit"
	"backend_project/services/receipt-engine/scraper-service/internal/providers/x5club"
)

func main() {
	port := "8003"
	serviceName := "scraper-service"
	demoMode := os.Getenv("DEMO_MODE") == "true"

	r := internal.NewRouter()

	oauthMgr := email.NewOAuthManager()
	imapCli := email.NewIMAPClient(oauthMgr)
	emailParser := email.NewParser()
	fnsHandler := fns.NewHandler(demoMode)
	x5Client := x5club.NewClient(demoMode)
	magnitClient := magnit.NewClient(demoMode)
	magnitINN := magnit.NewINNDetector()
	magnitMapper := magnit.NewMapper(magnitINN)

	r.Post("/api/v1/fns/ticket", fnsHandler.ServeHTTP)

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

	r.Post("/api/v1/x5club/sync", func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		phone := r.URL.Query().Get("phone")
		password := r.URL.Query().Get("password")

		if userID == "" || phone == "" || password == "" {
			if !demoMode {
				http.Error(w, `{"error":"user_id, phone and password required"}`, http.StatusBadRequest)
				return
			}
		}

		if err := x5Client.Login(phone, password); err != nil {
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

	r.Post("/api/v1/magnit/sync", func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		phone := r.URL.Query().Get("phone")
		password := r.URL.Query().Get("password")

		if userID == "" || phone == "" || password == "" {
			if !demoMode {
				http.Error(w, `{"error":"user_id, phone and password required"}`, http.StatusBadRequest)
				return
			}
		}

		if err := magnitClient.Login(phone, password); err != nil {
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

	fmt.Printf("Service %s started on port %s...\n", serviceName, port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		panic(err)
	}
}

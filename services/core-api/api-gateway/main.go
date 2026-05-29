package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

type route struct {
	prefix  string
	target  string
	noAuth  bool
}

var routes = []route{
	{prefix: "/api/v1/auth/",       target: "http://user-service:8001", noAuth: true},
	{prefix: "/api/v1/providers/",  target: "http://user-service:8001"},
	{prefix: "/api/v1/dashboard/",  target: "http://receipt-service:8002"},
	{prefix: "/api/v1/receipts/",   target: "http://receipt-service:8002"},
	{prefix: "/api/v1/fns/",        target: "http://scraper-service:8003"},
	{prefix: "/api/v1/x5club/",     target: "http://scraper-service:8003"},
	{prefix: "/api/v1/magnit/",     target: "http://scraper-service:8003"},
	{prefix: "/api/v1/email/",      target: "http://scraper-service:8003"},
	{prefix: "/api/v1/credits/",    target: "http://credit-service:8009"},
	{prefix: "/api/v1/banks/",      target: "http://bank-service:8011"},
	{prefix: "/api/v1/categories/", target: "http://category-service:8004"},
	{prefix: "/api/v1/budgets/",    target: "http://budget-service:8005"},
	{prefix: "/api/v1/goals/",      target: "http://goal-service:8006"},
	{prefix: "/api/v1/expenses/",   target: "http://ai-processor:8100"},
	{prefix: "/api/v1/insights/",   target: "http://analytics-service:8101"},
	{prefix: "/api/v1/scenarios/",  target: "http://analytics-service:8101"},
	{prefix: "/api/v1/forecast/",   target: "http://analytics-service:8101"},
	{prefix: "/api/v1/challenges/", target: "http://social-service:8102"},
	{prefix: "/api/v1/digest/",     target: "http://reporting-service:8010"},
}

func main() {
	port := "8000"
	serviceName := "api-gateway"
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "change-me-to-a-random-secret"
	}

	r := chi.NewRouter()

	r.Use(recoverer)
	r.Use(corsMiddleware)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	for _, rt := range routes {
		rt := rt
		targetURL, err := url.Parse(rt.target)
		if err != nil {
			log.Fatalf("invalid target %s: %v", rt.target, err)
		}
		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		r.HandleFunc(rt.prefix+"*", func(w http.ResponseWriter, r *http.Request) {
			if !rt.noAuth {
				tokenStr := extractToken(r)
				if tokenStr == "" {
					http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
					return
				}
				token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
					if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
					}
					return []byte(jwtSecret), nil
				})
				if err != nil || !token.Valid {
					http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
					return
				}
			}
			r.URL.Path = strings.TrimPrefix(r.URL.Path, rt.prefix)
			r.URL.Path = "/api/v1" + r.URL.Path
			proxy.ServeHTTP(w, r)
		})
	}

	fmt.Printf("Service %s started on port %s...\n", serviceName, port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		panic(err)
	}
}

func extractToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return ""
}

func recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

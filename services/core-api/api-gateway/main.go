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
	"github.com/go-chi/cors"
	"github.com/golang-jwt/jwt/v5"
)

type route struct {
	prefix  string
	target  string
	noAuth  bool
}

var routes = []route{
	{prefix: "/api/v1/auth/",       target: svcURL("USER_SERVICE_URL", "http://user-service:8001"), noAuth: true},
	{prefix: "/api/v1/providers/",  target: svcURL("USER_SERVICE_URL", "http://user-service:8001")},
	{prefix: "/api/v1/dashboard/",  target: svcURL("RECEIPT_SERVICE_URL", "http://receipt-service:8002")},
	{prefix: "/api/v1/receipts/",   target: svcURL("RECEIPT_SERVICE_URL", "http://receipt-service:8002")},
	{prefix: "/api/v1/fns/",        target: svcURL("SCRAPER_SERVICE_URL", "http://scraper-service:8003")},
	{prefix: "/api/v1/x5club/",     target: svcURL("SCRAPER_SERVICE_URL", "http://scraper-service:8003")},
	{prefix: "/api/v1/magnit/",     target: svcURL("SCRAPER_SERVICE_URL", "http://scraper-service:8003")},
	{prefix: "/api/v1/email/",      target: svcURL("SCRAPER_SERVICE_URL", "http://scraper-service:8003")},
	{prefix: "/api/v1/credits/",    target: svcURL("CREDIT_SERVICE_URL", "http://credit-service:8009")},
	{prefix: "/api/v1/banks/",      target: svcURL("BANK_SERVICE_URL", "http://bank-service:8011")},
	{prefix: "/api/v1/categories/", target: svcURL("CATEGORY_SERVICE_URL", "http://category-service:8004")},
	{prefix: "/api/v1/budgets/",    target: svcURL("BUDGET_SERVICE_URL", "http://budget-service:8005")},
	{prefix: "/api/v1/goals/",      target: svcURL("GOAL_SERVICE_URL", "http://goal-service:8006")},
	{prefix: "/api/v1/expenses/",   target: svcURL("AI_PROCESSOR_URL", "http://ai-processor:8100")},
	{prefix: "/api/v1/receipt/",    target: svcURL("AI_PROCESSOR_URL", "http://ai-processor:8100")},
	{prefix: "/api/v1/voice/",      target: svcURL("AI_PROCESSOR_URL", "http://ai-processor:8100")},
	{prefix: "/api/v1/onboarding/", target: svcURL("AI_PROCESSOR_URL", "http://ai-processor:8100")},
	{prefix: "/api/v1/insights/",   target: svcURL("ANALYTICS_SERVICE_URL", "http://analytics-service:8101")},
	{prefix: "/api/v1/scenarios/",  target: svcURL("ANALYTICS_SERVICE_URL", "http://analytics-service:8101")},
	{prefix: "/api/v1/forecast/",   target: svcURL("ANALYTICS_SERVICE_URL", "http://analytics-service:8101")},
	{prefix: "/api/v1/analytics/",  target: svcURL("ANALYTICS_SERVICE_URL", "http://analytics-service:8101")},
	{prefix: "/api/v1/ai/",         target: svcURL("ANALYTICS_SERVICE_URL", "http://analytics-service:8101")},
	{prefix: "/api/v1/challenges/", target: svcURL("SOCIAL_SERVICE_URL", "http://social-service:8102")},
	{prefix: "/api/v1/digest/",     target: svcURL("REPORTING_SERVICE_URL", "http://reporting-service:8010")},
}

func main() {
	port := "8000"
	serviceName := "api-gateway"
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "test-secret"
	}

	r := chi.NewRouter()

	r.Use(recoverer)
	r.Use(cors.Handler(corsOptions()))

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
		proxy.ModifyResponse = stripUpstreamCORS
		registerProxy(r, rt, proxy, jwtSecret)
	}

	fmt.Printf("Service %s started on port %s...\n", serviceName, port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		panic(err)
	}
}

func registerProxy(r *chi.Mux, rt route, proxy *httputil.ReverseProxy, jwtSecret string) {
	handler := func(w http.ResponseWriter, req *http.Request) {
		if !rt.noAuth {
			tokenStr := extractToken(req)
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
		proxy.ServeHTTP(w, req)
	}
	base := strings.TrimSuffix(rt.prefix, "/")
	r.HandleFunc(base, handler)
	r.HandleFunc(base+"/*", handler)
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


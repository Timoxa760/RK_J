package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/cors"
)

// corsHeaderNames — заголовки CORS, которые задаёт только api-gateway.
// Upstream-сервисы используют internal.NewRouter с тем же middleware;
// без удаления браузер получает дубликаты (например Access-Control-Allow-Origin: *,*).
var corsHeaderNames = []string{
	"Access-Control-Allow-Origin",
	"Access-Control-Allow-Methods",
	"Access-Control-Allow-Headers",
	"Access-Control-Allow-Credentials",
	"Access-Control-Expose-Headers",
	"Access-Control-Max-Age",
	"Vary",
}

// stripUpstreamCORS удаляет CORS-заголовки из ответа upstream-сервиса.
func stripUpstreamCORS(resp *http.Response) error {
	for _, name := range corsHeaderNames {
		resp.Header.Del(name)
	}
	return nil
}

// corsOptions возвращает настройки CORS для единой точки входа.
func corsOptions() cors.Options {
	origins := parseAllowedOrigins(os.Getenv("CORS_ALLOWED_ORIGINS"))
	if len(origins) == 0 {
		origins = []string{"*"}
	}

	return cors.Options{
		AllowedOrigins: origins,
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
			"X-Requested-With",
			"Origin",
		},
		ExposedHeaders:   []string{"Content-Length", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}
}

// parseAllowedOrigins разбирает список origin через запятую из переменной окружения.
func parseAllowedOrigins(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		origin := strings.TrimSpace(part)
		if origin != "" {
			origins = append(origins, origin)
		}
	}
	return origins
}

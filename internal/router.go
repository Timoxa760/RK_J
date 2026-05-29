package internal

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

var logger = log.New(os.Stdout, "", 0)

// NewRouter собирает базовый chi-роутер с middleware и эндпоинтом health-check
func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	// 1. Предохранитель: если в коде упадет panic, сервер не вылетит, а вернет 500 ошибку
	r.Use(middleware.Recoverer)

	// 2. Настройка CORS для хакатона (разрешаем все origins для разработки по ТЗ)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		MaxAge:           300,
	}))

	// 3. Подключаем наш кастомный JSON-логгер
	r.Use(JSONLogger)

	// 4. Обязательный эндпоинт /health из Группы 0
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	return r
}

// JSONLogger форматирует логи запросов в структурированный JSON для аналитики
func JSONLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// Передаем запрос дальше по цепочке к нашему обработчику
		next.ServeHTTP(ww, r)

		// Собираем данные для лога
		logEntry := map[string]interface{}{
			"level":       "info",
			"timestamp":   time.Now().Format(time.RFC3339),
			"method":      r.Method,
			"path":        r.URL.Path,
			"status":      ww.Status(),
			"latency_ms":  time.Since(start).Milliseconds(),
		}

		// Выводим структурированный лог в консоль
		if jsonLog, err := json.Marshal(logEntry); err == nil {
			logger.Println(string(jsonLog))
		}
	})
}
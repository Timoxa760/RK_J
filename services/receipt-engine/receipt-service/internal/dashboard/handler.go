package dashboard

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"backend_project/internal/auth"
	"backend_project/internal/expensestore"
)

type Handler struct {
	chConn      driver.Conn
	pgPool      *pgxpool.Pool
	expenseFile *expensestore.FileStore
	demoMode    bool
	jwtSecret   string
}

func New(demoMode bool, jwtSecret string) *Handler {
	if jwtSecret == "" {
		jwtSecret = auth.JWTSecret()
	}
	return &Handler{demoMode: demoMode, jwtSecret: jwtSecret}
}

func (h *Handler) ConnectClickHouse(ctx context.Context, host, user, password, database string) error {
	addr := host + ":9000"
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Database: database,
			Username: user,
			Password: password,
		},
		Settings:     clickhouse.Settings{"max_execution_time": 60},
		Compression:  &clickhouse.Compression{Method: clickhouse.CompressionLZ4},
		DialTimeout:  5 * time.Second,
		MaxOpenConns: 5,
		MaxIdleConns: 5,
	})
	if err != nil {
		return err
	}
	if err := conn.Ping(ctx); err != nil {
		return err
	}
	h.chConn = conn
	return nil
}

func (h *Handler) Close() error {
	if h.chConn != nil {
		return h.chConn.Close()
	}
	return nil
}

func (h *Handler) Register(r chi.Router) {
	r.Get("/api/v1/dashboard/sankey", h.sankey)
	r.Get("/api/v1/dashboard/categories", h.categories)
	r.Get("/api/v1/dashboard/stores", h.stores)
	r.Get("/api/v1/dashboard/compare", h.compare)
	r.Get("/api/v1/dashboard/timemachine", h.timemachine)
	r.Get("/api/v1/receipts", h.listReceipts)
	r.Delete("/api/v1/receipts/{id}", h.deleteReceipt)
}

// --- Sankey ---

type sankeyNode struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type sankeyLink struct {
	Source string  `json:"source"`
	Target string  `json:"target"`
	Value  float64 `json:"value"`
}

type sankeyResponse struct {
	Nodes []sankeyNode `json:"nodes"`
	Links []sankeyLink `json:"links"`
}

func (h *Handler) sankey(w http.ResponseWriter, r *http.Request) {
	if resp := h.buildSankeyResponse(r); resp != nil {
		writeJSON(w, *resp)
		return
	}
	writeJSON(w, sankeyResponse{Nodes: []sankeyNode{}, Links: []sankeyLink{}})
}

func (h *Handler) buildSankeyResponse(r *http.Request) *sankeyResponse {
	userID, err := auth.UserIDFromRequest(r, h.jwtSecret)
	if err != nil {
		return nil
	}
	if h.pgPool != nil {
		if cats := h.queryCategoriesPG(r.Context(), userID); cats != nil && len(cats.Categories) > 0 {
			resp := buildSankeyFromCategories(cats.Categories)
			return &resp
		}
	}
	if h.chConn != nil {
		if cats := h.queryCategoriesCH(r.Context(), userID); cats != nil && len(cats.Categories) > 0 {
			resp := buildSankeyFromCategories(cats.Categories)
			return &resp
		}
	}
	return nil
}

// --- Categories ---

type catItem struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Total    float64 `json:"total"`
}

type catSubcategory struct {
	Name  string   `json:"name"`
	Total float64  `json:"total"`
	Items []catItem `json:"items"`
}

type catCategory struct {
	Name          string           `json:"name"`
	Total         float64          `json:"total"`
	Subcategories []catSubcategory `json:"subcategories"`
}

type categoriesResponse struct {
	Categories []catCategory `json:"categories"`
}

func (h *Handler) categories(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFromRequest(r, h.jwtSecret)
	if h.pgPool != nil && err == nil {
		if resp := h.queryCategoriesPG(r.Context(), userID); resp != nil && len(resp.Categories) > 0 {
			writeJSON(w, resp)
			return
		}
	}
	if h.chConn != nil && err == nil {
		if resp := h.queryCategoriesCH(r.Context(), userID); resp != nil && len(resp.Categories) > 0 {
			writeJSON(w, resp)
			return
		}
	}
	writeJSON(w, categoriesResponse{Categories: []catCategory{}})
}

func (h *Handler) queryCategoriesCH(ctx context.Context, userID string) *categoriesResponse {
	rows, err := h.chConn.Query(ctx, `
		SELECT category, sum(total) as total
		FROM spending_by_category
		WHERE user_id = $1 AND month = toStartOfMonth(now())
		GROUP BY category
		ORDER BY total DESC
	`, userID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var cats []catCategory
	var totalTotal float64
	for rows.Next() {
		var name string
		var total float64
		if err := rows.Scan(&name, &total); err != nil {
			return nil
		}
		cats = append(cats, catCategory{
			Name:  name,
			Total: total,
		})
		totalTotal += total
	}

	if len(cats) == 0 {
		return nil
	}

	// add sub-items per category from receipt_items
	for i, c := range cats {
		subRows, err := h.chConn.Query(ctx, `
			SELECT item_name, price, quantity
			FROM receipt_items
			WHERE user_id = $1 AND category = $2 AND purchased_at >= toStartOfMonth(now())
			ORDER BY purchased_at DESC
			LIMIT 10
		`, userID, c.Name)
		if err != nil {
			continue
		}

		var items []catItem
		for subRows.Next() {
			var name string
			var price float64
			var qty uint32
			if err := subRows.Scan(&name, &price, &qty); err != nil {
				continue
			}
			items = append(items, catItem{
				Name:     name,
				Price:    price,
				Quantity: int(qty),
				Total:    price * float64(qty),
			})
		}
		subRows.Close()

		if len(items) > 0 {
			cats[i].Subcategories = []catSubcategory{
				{Name: c.Name, Total: c.Total, Items: items},
			}
		}
	}

	return &categoriesResponse{Categories: cats}
}

func (h *Handler) queryStoresCH(ctx context.Context, userID string) *storesResponse {
	rows, err := h.chConn.Query(ctx, `
		SELECT store, purchases, avg_check, total, impulse_ratio
		FROM store_aggregates
		WHERE user_id = $1
		ORDER BY purchases DESC
	`, userID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var stores []storeItem
	for rows.Next() {
		var s storeItem
		if err := rows.Scan(&s.Name, &s.Purchases, &s.AvgCheck, &s.Total, &s.ImpulseRatio); err != nil {
			return nil
		}
		stores = append(stores, s)
	}

	if len(stores) == 0 {
		return nil
	}
	return &storesResponse{Stores: stores}
}

// --- Stores ---

type storeItem struct {
	Name         string  `json:"name"`
	AvgCheck     float64 `json:"avg_check"`
	Purchases    int     `json:"purchases"`
	Total        float64 `json:"total"`
	ImpulseRatio float64 `json:"impulse_ratio"`
}

type storesResponse struct {
	Stores []storeItem `json:"stores"`
}

func (h *Handler) stores(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFromRequest(r, h.jwtSecret)
	if h.pgPool != nil && err == nil {
		if resp := h.queryStoresPG(r.Context(), userID); resp != nil && len(resp.Stores) > 0 {
			writeJSON(w, resp)
			return
		}
	}
	if h.chConn != nil && err == nil {
		if resp := h.queryStoresCH(r.Context(), userID); resp != nil && len(resp.Stores) > 0 {
			writeJSON(w, resp)
			return
		}
	}
	writeJSON(w, storesResponse{Stores: []storeItem{}})
}

// --- Compare ---

type compareCategory struct {
	Name  string  `json:"name"`
	Total float64 `json:"total"`
}

type compareMonth struct {
	Label      string            `json:"label"`
	Categories []compareCategory `json:"categories"`
}

type compareChange struct {
	Category    string `json:"category"`
	Delta       int    `json:"delta"`
	DeltaPercent int   `json:"delta_percent"`
}

type compareInsights struct {
	BiggestChange compareChange `json:"biggest_change"`
}

type compareResponse struct {
	Months   []compareMonth  `json:"months"`
	Insights compareInsights `json:"insights"`
}

func (h *Handler) compare(w http.ResponseWriter, r *http.Request) {
	monthCount := parseMonthCount(r)
	if h.pgPool != nil {
		if userID, err := auth.UserIDFromRequest(r, h.jwtSecret); err == nil {
			if resp := h.queryComparePG(r.Context(), userID, monthCount); resp != nil {
				writeJSON(w, resp)
				return
			}
		}
	}
	writeJSON(w, compareResponse{Months: []compareMonth{}})
}

// --- Time Machine ---

type timemachineResponse struct {
	Months           []string  `json:"months"`
	RealSavings      []float64 `json:"real_savings"`
	OptimizedSavings []float64 `json:"optimized_savings"`
	DifferenceFinal  int       `json:"difference_final"`
}

func (h *Handler) timemachine(w http.ResponseWriter, r *http.Request) {
	if h.pgPool != nil {
		if userID, err := auth.UserIDFromRequest(r, h.jwtSecret); err == nil {
			if resp := h.queryTimemachinePG(r.Context(), userID); resp != nil {
				writeJSON(w, resp)
				return
			}
		}
	}
	writeJSON(w, timemachineResponse{
		Months:           []string{},
		RealSavings:      []float64{},
		OptimizedSavings: []float64{},
	})
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(v)
}

func parseMonthCount(r *http.Request) int {
	const defaultMonths = 2
	v := r.URL.Query().Get("months")
	if v == "" {
		return defaultMonths
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 1 || n > 12 {
		return defaultMonths
	}
	return n
}

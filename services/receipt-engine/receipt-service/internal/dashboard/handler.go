package dashboard

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	chConn   driver.Conn
	demoMode bool
}

func New(demoMode bool) *Handler {
	return &Handler{demoMode: demoMode}
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
	resp := sankeyResponse{
		Nodes: []sankeyNode{
			{Name: "Зарплата", Value: 180000},
			{Name: "Накопления", Value: 35000},
			{Name: "Продукты", Value: 52000},
			{Name: "Кафе и рестораны", Value: 28000},
			{Name: "Транспорт", Value: 15000},
			{Name: "Развлечения", Value: 12000},
			{Name: "Здоровье", Value: 7000},
			{Name: "Одежда", Value: 6000},
			{Name: "Доставка", Value: 14000},
			{Name: "Связь", Value: 3000},
			{Name: "Прочее", Value: 8000},
		},
		Links: []sankeyLink{
			{Source: "Зарплата", Target: "Накопления", Value: 35000},
			{Source: "Зарплата", Target: "Продукты", Value: 52000},
			{Source: "Зарплата", Target: "Кафе и рестораны", Value: 28000},
			{Source: "Зарплата", Target: "Транспорт", Value: 15000},
			{Source: "Зарплата", Target: "Развлечения", Value: 12000},
			{Source: "Зарплата", Target: "Здоровье", Value: 7000},
			{Source: "Зарплата", Target: "Одежда", Value: 6000},
			{Source: "Зарплата", Target: "Доставка", Value: 14000},
			{Source: "Зарплата", Target: "Связь", Value: 3000},
			{Source: "Зарплата", Target: "Прочее", Value: 8000},
		},
	}
	writeJSON(w, resp)
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
	if h.chConn != nil {
		if resp := h.queryCategories(r.Context()); resp != nil {
			writeJSON(w, resp)
			return
		}
	}

	resp := categoriesResponse{
		Categories: []catCategory{
			{
				Name: "Продукты", Total: 52000,
				Subcategories: []catSubcategory{
					{
						Name: "Молочные", Total: 8500,
						Items: []catItem{
							{Name: "Молоко 3.2%", Price: 78, Quantity: 12, Total: 936},
							{Name: "Творог 5%", Price: 120, Quantity: 6, Total: 720},
							{Name: "Сметана", Price: 85, Quantity: 4, Total: 340},
							{Name: "Кефир", Price: 65, Quantity: 8, Total: 520},
							{Name: "Йогурт", Price: 45, Quantity: 10, Total: 450},
						},
					},
					{
						Name: "Мясо", Total: 7200,
						Items: []catItem{
							{Name: "Курица", Price: 350, Quantity: 8, Total: 2800},
							{Name: "Говядина", Price: 600, Quantity: 4, Total: 2400},
							{Name: "Фарш", Price: 400, Quantity: 3, Total: 1200},
							{Name: "Сосиски", Price: 250, Quantity: 3, Total: 750},
						},
					},
					{
						Name: "Овощи и фрукты", Total: 6400,
						Items: []catItem{
							{Name: "Картофель", Price: 40, Quantity: 10, Total: 400},
							{Name: "Помидоры", Price: 180, Quantity: 6, Total: 1080},
							{Name: "Огурцы", Price: 120, Quantity: 5, Total: 600},
							{Name: "Бананы", Price: 90, Quantity: 8, Total: 720},
							{Name: "Яблоки", Price: 110, Quantity: 6, Total: 660},
						},
					},
					{
						Name: "Бакалея", Total: 5400,
						Items: []catItem{
							{Name: "Рис", Price: 120, Quantity: 5, Total: 600},
							{Name: "Макароны", Price: 80, Quantity: 8, Total: 640},
							{Name: "Мука", Price: 55, Quantity: 4, Total: 220},
							{Name: "Соль", Price: 25, Quantity: 3, Total: 75},
							{Name: "Масло растительное", Price: 150, Quantity: 4, Total: 600},
						},
					},
					{
						Name: "Хлеб и выпечка", Total: 3100,
						Items: []catItem{
							{Name: "Хлеб белый", Price: 45, Quantity: 15, Total: 675},
							{Name: "Хлеб ржаной", Price: 55, Quantity: 10, Total: 550},
							{Name: "Батон", Price: 40, Quantity: 8, Total: 320},
						},
					},
				},
			},
			{
				Name: "Кафе и рестораны", Total: 28000,
				Subcategories: []catSubcategory{
					{
						Name: "Кофе навынос", Total: 5600,
						Items: []catItem{
							{Name: "Латте", Price: 210, Quantity: 18, Total: 3780},
							{Name: "Капучино", Price: 190, Quantity: 14, Total: 2660},
							{Name: "Эспрессо", Price: 120, Quantity: 8, Total: 960},
						},
					},
					{
						Name: "Обеды", Total: 12000,
						Items: []catItem{
							{Name: "Бизнес-ланч", Price: 450, Quantity: 15, Total: 6750},
							{Name: "Суп", Price: 280, Quantity: 8, Total: 2240},
							{Name: "Салат", Price: 320, Quantity: 6, Total: 1920},
						},
					},
					{
						Name: "Ужины", Total: 10400,
						Items: []catItem{
							{Name: "Паста", Price: 520, Quantity: 6, Total: 3120},
							{Name: "Стейк", Price: 890, Quantity: 4, Total: 3560},
							{Name: "Пицца", Price: 620, Quantity: 6, Total: 3720},
						},
					},
				},
			},
			{
				Name: "Транспорт", Total: 15000,
				Subcategories: []catSubcategory{
					{
						Name: "Метро", Total: 4200,
						Items: []catItem{
							{Name: "Проездной", Price: 2650, Quantity: 1, Total: 2650},
							{Name: "Разовые поездки", Price: 65, Quantity: 24, Total: 1560},
						},
					},
					{
						Name: "Такси", Total: 6800,
						Items: []catItem{
							{Name: "Яндекс.Такси", Price: 350, Quantity: 14, Total: 4900},
							{Name: "Ситимобил", Price: 320, Quantity: 6, Total: 1920},
						},
					},
					{
						Name: "Каршеринг", Total: 4000,
						Items: []catItem{
							{Name: "Яндекс.Драйв", Price: 450, Quantity: 6, Total: 2700},
							{Name: "Делимобиль", Price: 420, Quantity: 3, Total: 1260},
						},
					},
				},
			},
			{
				Name: "Развлечения", Total: 12000,
				Subcategories: []catSubcategory{
					{
						Name: "Подписки", Total: 4500,
						Items: []catItem{
							{Name: "Яндекс.Плюс", Price: 299, Quantity: 3, Total: 897},
							{Name: "Kion", Price: 399, Quantity: 3, Total: 1197},
							{Name: "Spotify", Price: 249, Quantity: 3, Total: 747},
						},
					},
					{
						Name: "Кино", Total: 3500,
						Items: []catItem{
							{Name: "Билеты", Price: 450, Quantity: 6, Total: 2700},
							{Name: "Попкорн", Price: 320, Quantity: 4, Total: 1280},
						},
					},
					{
						Name: "Игры", Total: 4000,
						Items: []catItem{
							{Name: "Steam", Price: 1500, Quantity: 2, Total: 3000},
						},
					},
				},
			},
		},
	}
	writeJSON(w, resp)
}

func (h *Handler) queryCategories(ctx context.Context) *categoriesResponse {
	rows, err := h.chConn.Query(ctx, `
		SELECT category, sum(total) as total
		FROM spending_by_category
		WHERE month = toStartOfMonth(now())
		GROUP BY category
		ORDER BY total DESC
	`)
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
			WHERE category = $1 AND purchased_at >= toStartOfMonth(now())
			ORDER BY purchased_at DESC
			LIMIT 10
		`, c.Name)
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

func (h *Handler) queryStores(ctx context.Context) *storesResponse {
	rows, err := h.chConn.Query(ctx, `
		SELECT store, purchases, avg_check, total, impulse_ratio
		FROM store_aggregates
		ORDER BY purchases DESC
	`)
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
	if h.chConn != nil {
		if resp := h.queryStores(r.Context()); resp != nil {
			writeJSON(w, resp)
			return
		}
	}

	resp := storesResponse{
		Stores: []storeItem{
			{Name: "Пятёрочка", AvgCheck: 650, Purchases: 14, Total: 9100, ImpulseRatio: 0.25},
			{Name: "Магнит", AvgCheck: 720, Purchases: 10, Total: 7200, ImpulseRatio: 0.20},
			{Name: "ВкусВилл", AvgCheck: 980, Purchases: 7, Total: 6860, ImpulseRatio: 0.10},
			{Name: "Ozon", AvgCheck: 2100, Purchases: 4, Total: 8400, ImpulseRatio: 0.65},
			{Name: "Wildberries", AvgCheck: 1850, Purchases: 5, Total: 9250, ImpulseRatio: 0.70},
			{Name: "Лента", AvgCheck: 820, Purchases: 8, Total: 6560, ImpulseRatio: 0.15},
		},
	}
	writeJSON(w, resp)
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
	resp := compareResponse{
		Months: []compareMonth{
			{
				Label: "Март 2026",
				Categories: []compareCategory{
					{Name: "Продукты", Total: 52000},
					{Name: "Кафе и рестораны", Total: 28000},
					{Name: "Транспорт", Total: 15000},
					{Name: "Развлечения", Total: 12000},
					{Name: "Здоровье", Total: 7000},
					{Name: "Одежда", Total: 6000},
					{Name: "Доставка", Total: 14000},
					{Name: "Связь", Total: 3000},
					{Name: "Прочее", Total: 8000},
				},
			},
			{
				Label: "Февраль 2026",
				Categories: []compareCategory{
					{Name: "Продукты", Total: 48000},
					{Name: "Кафе и рестораны", Total: 25000},
					{Name: "Транспорт", Total: 14000},
					{Name: "Развлечения", Total: 11000},
					{Name: "Здоровье", Total: 6500},
					{Name: "Одежда", Total: 5500},
					{Name: "Доставка", Total: 12000},
					{Name: "Связь", Total: 3000},
					{Name: "Прочее", Total: 7500},
				},
			},
		},
		Insights: compareInsights{
			BiggestChange: compareChange{
				Category:     "Кафе и рестораны",
				Delta:        3000,
				DeltaPercent: 12,
			},
		},
	}
	writeJSON(w, resp)
}

// --- Time Machine ---

type timemachineResponse struct {
	Months           []string  `json:"months"`
	RealSavings      []float64 `json:"real_savings"`
	OptimizedSavings []float64 `json:"optimized_savings"`
	DifferenceFinal  int       `json:"difference_final"`
}

func (h *Handler) timemachine(w http.ResponseWriter, r *http.Request) {
	months := make([]string, 60)
	realSavings := make([]float64, 60)
	optSavings := make([]float64, 60)

	start := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	monthlyIncome := 180000.0
	monthlySpend := 145000.0
	monthlySave := monthlyIncome - monthlySpend
	deliveryAndCoffee := 14000.0 + 28000.0*0.25
	extraSave := deliveryAndCoffee * 0.3

	real := 500000.0
	opt := 500000.0

	for m := 0; m < 60; m++ {
		months[m] = start.AddDate(0, m, 0).Format("2006-01")
		growth := 1 + 0.003*float64(m)
		real += monthlySave*growth + float64(50000)
		opt += (monthlySave + extraSave)*growth + float64(50000)
		if m%12 == 11 {
			real += 50000
			opt += 50000
		}
		realSavings[m] = real
		optSavings[m] = opt
	}

	resp := timemachineResponse{
		Months:           months,
		RealSavings:      realSavings,
		OptimizedSavings: optSavings,
		DifferenceFinal:  int(optSavings[59] - realSavings[59]),
	}
	writeJSON(w, resp)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

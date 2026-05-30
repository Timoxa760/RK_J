package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type expense struct {
	UserID      string  `json:"user_id"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
	Source      string  `json:"source"`
}

var demoExpenses = []expense{
	{UserID: "+79991234567", Category: "Продукты", Description: "Молоко, хлеб, сыр"},
	{UserID: "+79991234567", Category: "Продукты", Description: "Овощи, фрукты, мясо"},
	{UserID: "+79991234567", Category: "Кафе и рестораны", Description: "Обед в столовой"},
	{UserID: "+79991234567", Category: "Транспорт", Description: "Бензин"},
	{UserID: "+79991234567", Category: "Доставка", Description: "Яндекс.Еда"},
	{UserID: "+79991234567", Category: "Подписки", Description: "Яндекс.Плюс"},
	{UserID: "+79991234567", Category: "ЖКХ", Description: "Квартплата"},
	{UserID: "+79991234567", Category: "Развлечения", Description: "Кино"},
	{UserID: "+79991234567", Category: "Здоровье", Description: "Аптека"},
	{UserID: "+79991234567", Category: "Прочие расходы", Description: "Подарок другу"},
	{UserID: "+79992345678", Category: "Продукты", Description: "Закупка в Пятёрочке"},
	{UserID: "+79992345678", Category: "Кафе и рестораны", Description: "Кофе навынос"},
	{UserID: "+79992345678", Category: "Транспорт", Description: "Метро"},
	{UserID: "+79992345678", Category: "Доставка", Description: "Delivery Club"},
	{UserID: "+79992345678", Category: "Подписки", Description: "Tele2"},
	{UserID: "+79992345678", Category: "Развлечения", Description: "Стим"},
	{UserID: "+79992345678", Category: "Одежда", Description: "Футболка"},
	{UserID: "+79992345678", Category: "Прочие расходы", Description: "Ремонт телефона"},
	{UserID: "+79993456789", Category: "Продукты", Description: "Корзина в Ашане"},
	{UserID: "+79993456789", Category: "Кафе и рестораны", Description: "Ужин в ресторане"},
	{UserID: "+79993456789", Category: "Транспорт", Description: "Такси"},
	{UserID: "+79993456789", Category: "Доставка", Description: "Яндекс.Еда"},
	{UserID: "+79993456789", Category: "Подписки", Description: "Домашний интернет"},
	{UserID: "+79993456789", Category: "ЖКХ", Description: "Электричество"},
	{UserID: "+79993456789", Category: "Одежда", Description: "Платье"},
	{UserID: "+79993456789", Category: "Здоровье", Description: "Спортзал"},
	{UserID: "+79993456789", Category: "Развлечения", Description: "Концерт"},
	{UserID: "+79993456789", Category: "Прочие расходы", Description: "Косметика"},
}

func main() {
	baseURL := os.Getenv("SEED_API_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8000"
	}
	api := baseURL + "/api/v1"
	phone := envStr("SEED_DEMO_PHONE", "+79991234567")
	code := envStr("SEED_DEMO_CODE", "0000")

	token := login(api, phone, code)
	if token != "" {
		seedProfile(api, token)
	} else {
		log.Println("warn: no JWT — profile seed skipped")
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for _, e := range demoExpenses {
		amount := 100.0 + float64(rng.Intn(5000)) + float64(rng.Intn(99))/100.0
		daysAgo := rng.Intn(60)
		date := time.Now().AddDate(0, 0, -daysAgo).Format("2006-01-02")

		e.Amount = amount
		e.Date = date
		e.Source = "seed"

		body, _ := json.Marshal(e)
		req, err := http.NewRequest(http.MethodPost, api+"/expenses/manual", bytes.NewReader(body))
		if err != nil {
			log.Printf("error: %v", err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("error: %s %s: %v", e.UserID, e.Description, err)
			continue
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			fmt.Printf("  [%s] %.0f₽ %s - %s\n", e.UserID[:8], amount, e.Category, e.Description)
		} else {
			fmt.Printf("  [%s] FAILED %d: %s\n", e.UserID[:8], resp.StatusCode, e.Description)
		}
	}

	fmt.Println("\nseed complete!")
}

func login(api, phone, code string) string {
	body, _ := json.Marshal(map[string]string{"phone": phone, "code": code})
	resp, err := http.Post(api+"/auth/login", "application/json", bytes.NewReader(body))
	if err != nil {
		log.Printf("login error: %v", err)
		return ""
	}
	defer resp.Body.Close()
	var out struct {
		AccessToken string `json:"access_token"`
	}
	if json.NewDecoder(resp.Body).Decode(&out) != nil || out.AccessToken == "" {
		return ""
	}
	fmt.Println("demo profile seed: JWT ok")
	return out.AccessToken
}

func seedProfile(api, token string) {
	payload := map[string]any{
		"active_income":          150_000,
		"passive_income":         30_000,
		"emergency_fund":         340_000,
		"goal_title":             "Отпуск",
		"goal_amount":            150_000,
		"goal_kind":              "save",
		"onboarding_completed":   true,
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPatch, api+"/users/me/profile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("profile seed error: %v", err)
		return
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		fmt.Println("demo profile seeded")
	} else {
		log.Printf("profile seed HTTP %d", resp.StatusCode)
	}
}

func envStr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

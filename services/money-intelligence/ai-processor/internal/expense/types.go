package expense

// Item — одна распознанная трата.
type Item struct {
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
}

// Result — результат парсинга фразы пользователя.
type Result struct {
	Expenses  []Item
	Advice    string
	ParsedBy  string // onlysq | regex
	Parsed    bool
}

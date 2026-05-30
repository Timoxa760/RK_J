package receipt

// ManualRequest — тело POST /receipt/manual (контракт Nuxt front).
type ManualRequest struct {
	Store    string  `json:"store"`
	Amount   float64 `json:"amount"`
	Category string  `json:"category"`
	Date     string  `json:"date"`
}

// ManualResponse — ответ POST /receipt/manual для front.
type ManualResponse struct {
	ReceiptID string  `json:"receipt_id"`
	Store     string  `json:"store"`
	Amount    float64 `json:"amount"`
	Category  string  `json:"category"`
	Date      string  `json:"date"`
	Status    string  `json:"status"`
}

// VoiceLineItem — позиция в голосовом чеке.
type VoiceLineItem struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity,omitempty"`
	Category string  `json:"category,omitempty"`
}

// VoiceResponse — ответ POST /receipt/voice для front.
type VoiceResponse struct {
	ReceiptID  string          `json:"receipt_id"`
	Store      string          `json:"store"`
	Items      []VoiceLineItem `json:"items"`
	Total      float64         `json:"total"`
	Category   string          `json:"category"`
	Confidence float64         `json:"confidence"`
}

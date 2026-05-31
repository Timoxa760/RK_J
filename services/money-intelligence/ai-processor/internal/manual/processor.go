package manual

import (
	"context"
	"time"

	"backend_project/internal/auth"
	"backend_project/services/money-intelligence/ai-processor/internal/expense"
)

// ExpenseDTO — трата в HTTP-ответе.
type ExpenseDTO struct {
	ID          string  `json:"id,omitempty"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description,omitempty"`
}

// CreateRequest — тело POST /expenses/manual.
type CreateRequest struct {
	UserID      string  `json:"user_id"`
	RawText     string  `json:"raw_text"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
	Source      string  `json:"source"`
}

// CreateResponse — ответ expenses API.
type CreateResponse struct {
	Success     bool         `json:"success"`
	ID          string       `json:"id,omitempty"`
	Amount      float64      `json:"amount"`
	Category    string       `json:"category"`
	Description string       `json:"description,omitempty"`
	Parsed      bool         `json:"parsed"`
	ParsedBy   string       `json:"parsed_by,omitempty"`
	Advice     string       `json:"advice,omitempty"`
	Transcript string       `json:"transcript,omitempty"`
	Source     string       `json:"source,omitempty"`
	Expenses   []ExpenseDTO `json:"expenses,omitempty"`
}

// Storage сохраняет распознанные траты.
type Storage interface {
	SaveExpenses(ctx context.Context, userID, source string, date time.Time, items []expense.Item) ([]ExpenseDTO, error)
}

// Processor оркестрирует парсинг и сохранение.
type Processor struct {
	parser *expense.Parser
	store  Storage
}

// NewProcessor создаёт обработчик manual/voice pipeline.
func NewProcessor(parser *expense.Parser, store Storage) *Processor {
	return &Processor{parser: parser, store: store}
}

// Create обрабатывает текстовый ввод.
func (p *Processor) Create(ctx context.Context, req CreateRequest) (CreateResponse, int, error) {
	req.UserID = auth.NormalizePhone(req.UserID)
	if req.UserID == "" {
		return CreateResponse{}, 400, &APIError{
			Code:    "user_id_required",
			Message: "Не удалось определить пользователя. Войдите снова.",
		}
	}
	if req.Source == "" {
		req.Source = "manual"
	}
	date := time.Now()
	if req.Date != "" {
		if t, err := time.Parse("2006-01-02", req.Date); err == nil {
			date = t
		}
	}

	parsed := p.parser.Parse(ctx, expense.ParseInput{
		RawText:     req.RawText,
		Amount:      req.Amount,
		Category:    req.Category,
		Description: req.Description,
	})
	if len(parsed.Expenses) == 0 {
		return CreateResponse{}, 400, &APIError{
			Code:    "amount_required",
			Message: "Не услышали сумму. Скажите, например: «колбаса 300 рублей».",
		}
	}

	saved, err := p.store.SaveExpenses(ctx, req.UserID, req.Source, date, parsed.Expenses)
	if err != nil {
		return CreateResponse{}, 500, &APIError{
			Code:    "save_failed",
			Message: "Не удалось сохранить покупку. Попробуйте ещё раз.",
		}
	}

	return buildResponse(saved, parsed, req.Source, ""), 200, nil
}

// CreateFromTranscript обрабатывает текст после Whisper.
func (p *Processor) CreateFromTranscript(ctx context.Context, userID, transcript, dateStr string) (CreateResponse, int, error) {
	return p.Create(ctx, CreateRequest{
		UserID:  userID,
		RawText: transcript,
		Source:  "voice",
		Date:    dateStr,
	})
}

func buildResponse(saved []ExpenseDTO, parsed expense.Result, source, transcript string) CreateResponse {
	resp := CreateResponse{
		Success:    true,
		Parsed:     parsed.Parsed,
		ParsedBy:   parsed.ParsedBy,
		Advice:     parsed.Advice,
		Transcript: transcript,
		Source:     source,
		Expenses:   saved,
	}
	if len(saved) > 0 {
		resp.ID = saved[0].ID
		resp.Amount = saved[0].Amount
		resp.Category = saved[0].Category
		resp.Description = saved[0].Description
	}
	return resp
}

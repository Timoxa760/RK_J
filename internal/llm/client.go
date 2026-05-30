package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const defaultBaseURL = "https://generativelanguage.googleapis.com/v1beta"
const defaultModel = "gemini-2.0-flash"

// Client — HTTP-клиент Google Gemini (generateContent).
type Client struct {
	baseURL    string
	apiKey     string
	model      string
	httpClient *http.Client
}

// NewClient создаёт клиент. Пустой apiKey означает, что LLM недоступен.
func NewClient(baseURL, apiKey, model string) *Client {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	if model == "" {
		model = defaultModel
	}
	return &Client{
		baseURL: baseURL,
		apiKey:  strings.TrimSpace(apiKey),
		model:   model,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// NewFromEnv читает GEMINI_API_KEY / GEMINI_MODEL из окружения.
func NewFromEnv() *Client {
	key := strings.TrimSpace(os.Getenv("GEMINI_API_KEY"))
	if key == "" {
		key = strings.TrimSpace(os.Getenv("GOOGLE_API_KEY"))
	}
	return NewClient(os.Getenv("GEMINI_BASE_URL"), key, os.Getenv("GEMINI_MODEL"))
}

// Enabled возвращает true, если задан API-ключ.
func (c *Client) Enabled() bool {
	return c != nil && c.apiKey != ""
}

// Complete выполняет generateContent и возвращает текст ответа.
func (c *Client) Complete(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	if !c.Enabled() {
		return "", fmt.Errorf("gemini: api key not configured")
	}

	body, err := json.Marshal(generateRequest{
		SystemInstruction: &contentBlock{
			Parts: []part{{Text: systemPrompt}},
		},
		Contents: []contentMessage{
			{
				Role:  "user",
				Parts: []part{{Text: userPrompt}},
			},
		},
		GenerationConfig: generationConfig{
			Temperature:     0.2,
			MaxOutputTokens: 4096,
		},
	})
	if err != nil {
		return "", fmt.Errorf("gemini: marshal: %w", err)
	}

	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", c.baseURL, c.model, c.apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("gemini: request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("gemini: http: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", fmt.Errorf("gemini: read: %w", err)
	}

	var parsed generateResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", fmt.Errorf("gemini: decode: %w", err)
	}
	if parsed.Error != nil && parsed.Error.Message != "" {
		return "", fmt.Errorf("gemini: api error: %s", parsed.Error.Message)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("gemini: status %d: %s", resp.StatusCode, string(raw))
	}

	text := extractText(parsed)
	if text == "" {
		return "", fmt.Errorf("gemini: empty response")
	}
	return text, nil
}

func extractText(resp generateResponse) string {
	if len(resp.Candidates) == 0 {
		return ""
	}
	var b strings.Builder
	for _, p := range resp.Candidates[0].Content.Parts {
		if t := strings.TrimSpace(p.Text); t != "" {
			if b.Len() > 0 {
				b.WriteByte('\n')
			}
			b.WriteString(t)
		}
	}
	return strings.TrimSpace(b.String())
}

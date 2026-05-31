package llm

import (
	"bufio"
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
const defaultModel = "gemini-2.5-flash"

// Antigravity Tools — OpenAI-совместимый маршрут (Gemini native для этого аккаунта блокируется по региону).
const defaultAntigravityBaseURL = "http://127.0.0.1:8045/v1"
const defaultAntigravityModel = "claude-sonnet-4-6"

// Client — LLM через Google Gemini API или Antigravity Tools (/v1/chat/completions).
type Client struct {
	baseURL    string
	apiKey     string
	model      string
	httpClient *http.Client
}

// NewClient создаёт клиент. Пустой apiKey означает, что LLM недоступен.
func NewClient(baseURL, apiKey, model string) *Client {
	baseURL = normalizeBaseURL(baseURL)
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	if model == "" {
		if isProxyBaseURL(baseURL) {
			model = defaultAntigravityModel
		} else {
			model = defaultModel
		}
	}
	return &Client{
		baseURL: baseURL,
		apiKey:  strings.TrimSpace(apiKey),
		model:   model,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// NewFromEnv читает GEMINI_* / GOOGLE_API_KEY из окружения.
func NewFromEnv() *Client {
	key := strings.TrimSpace(os.Getenv("GEMINI_API_KEY"))
	if key == "" {
		key = strings.TrimSpace(os.Getenv("GOOGLE_API_KEY"))
	}

	baseURL := strings.TrimSpace(os.Getenv("GEMINI_BASE_URL"))
	provider := strings.ToLower(strings.TrimSpace(os.Getenv("GEMINI_PROVIDER")))
	if baseURL == "" && provider == "antigravity" {
		baseURL = defaultAntigravityBaseURL
	}

	model := strings.TrimSpace(os.Getenv("GEMINI_MODEL"))
	if model == "" && (provider == "antigravity" || isProxyBaseURL(baseURL)) {
		model = defaultAntigravityModel
	}

	return NewClient(baseURL, key, model)
}

// Enabled возвращает true, если задан API-ключ.
func (c *Client) Enabled() bool {
	return c != nil && c.apiKey != ""
}

func normalizeBaseURL(baseURL string) string {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		return ""
	}
	// Antigravity: используем OpenAI-маршрут /v1, не native /v1beta (Gemini location block).
	if strings.HasSuffix(baseURL, "/v1beta") &&
		!strings.Contains(baseURL, "generativelanguage.googleapis.com") {
		baseURL = strings.TrimSuffix(baseURL, "/v1beta") + "/v1"
	}
	return baseURL
}

func isProxyBaseURL(baseURL string) bool {
	return baseURL != "" && !strings.Contains(baseURL, "generativelanguage.googleapis.com")
}

func (c *Client) isGoogleDirect() bool {
	return strings.Contains(c.baseURL, "generativelanguage.googleapis.com")
}

func (c *Client) applyAuth(req *http.Request) {
	if c.apiKey == "" {
		return
	}
	if c.isGoogleDirect() {
		return
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("x-api-key", c.apiKey)
}

func (c *Client) geminiEndpointURL(method string, stream bool) string {
	path := fmt.Sprintf("%s/models/%s:%s", c.baseURL, c.model, method)
	if stream {
		return fmt.Sprintf("%s?alt=sse&key=%s", path, c.apiKey)
	}
	return fmt.Sprintf("%s?key=%s", path, c.apiKey)
}

func (c *Client) newGeminiRequest(ctx context.Context, method string, stream bool, systemPrompt, userPrompt string) (*http.Request, error) {
	body, err := json.Marshal(generateRequest{
		SystemInstruction: &contentBlock{
			Parts: []part{{Text: systemPrompt}},
		},
		Contents: []contentMessage{
			{Role: "user", Parts: []part{{Text: userPrompt}}},
		},
		GenerationConfig: generationConfig{
			Temperature:     0.2,
			MaxOutputTokens: 4096,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("llm: marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.geminiEndpointURL(method, stream), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("llm: request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (c *Client) newChatRequest(ctx context.Context, stream bool, systemPrompt, userPrompt string) (*http.Request, error) {
	body, err := json.Marshal(chatCompletionRequest{
		Model: c.model,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.2,
		MaxTokens:   4096,
		Stream:      stream,
	})
	if err != nil {
		return nil, fmt.Errorf("llm: marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("llm: request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	c.applyAuth(req)
	return req, nil
}

// Complete выполняет запрос и возвращает текст ответа.
func (c *Client) Complete(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	if !c.Enabled() {
		return "", fmt.Errorf("llm: api key not configured")
	}

	var req *http.Request
	var err error
	if c.isGoogleDirect() {
		req, err = c.newGeminiRequest(ctx, "generateContent", false, systemPrompt, userPrompt)
	} else {
		req, err = c.newChatRequest(ctx, false, systemPrompt, userPrompt)
	}
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("llm: http: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", fmt.Errorf("llm: read: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("llm: status %d: %s", resp.StatusCode, string(raw))
	}

	if c.isGoogleDirect() {
		var parsed generateResponse
		if err := json.Unmarshal(raw, &parsed); err != nil {
			return "", fmt.Errorf("llm: decode: %w", err)
		}
		if parsed.Error != nil && parsed.Error.Message != "" {
			return "", fmt.Errorf("llm: api error: %s", parsed.Error.Message)
		}
		text := extractGeminiText(parsed)
		if text == "" {
			return "", fmt.Errorf("llm: empty response")
		}
		return text, nil
	}

	var parsed chatCompletionResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", fmt.Errorf("llm: decode: %w", err)
	}
	if parsed.Error != nil && parsed.Error.Message != "" {
		return "", fmt.Errorf("llm: api error: %s", parsed.Error.Message)
	}
	text := extractChatText(parsed)
	if text == "" {
		return "", fmt.Errorf("llm: empty response")
	}
	return text, nil
}

func extractGeminiText(resp generateResponse) string {
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

func extractChatText(resp chatCompletionResponse) string {
	if len(resp.Choices) == 0 {
		return ""
	}
	msg := resp.Choices[0].Message
	if t := strings.TrimSpace(msg.Content); t != "" {
		return t
	}
	return strings.TrimSpace(msg.ReasoningContent)
}

func extractChatDelta(chunk chatCompletionResponse) string {
	if len(chunk.Choices) == 0 {
		return ""
	}
	d := chunk.Choices[0].Delta
	if t := strings.TrimSpace(d.Content); t != "" {
		return t
	}
	return strings.TrimSpace(d.ReasoningContent)
}

// StreamComplete стримит фрагменты текста через onDelta и возвращает полный ответ.
func (c *Client) StreamComplete(ctx context.Context, systemPrompt, userPrompt string, onDelta func(string) error) (string, error) {
	if !c.Enabled() {
		return "", fmt.Errorf("llm: api key not configured")
	}

	var req *http.Request
	var err error
	if c.isGoogleDirect() {
		req, err = c.newGeminiRequest(ctx, "streamGenerateContent", true, systemPrompt, userPrompt)
	} else {
		req, err = c.newChatRequest(ctx, true, systemPrompt, userPrompt)
	}
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("llm: http: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		return "", fmt.Errorf("llm: status %d: %s", resp.StatusCode, string(raw))
	}

	var full strings.Builder
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || !strings.HasPrefix(line, "data: ") {
			continue
		}
		payload := strings.TrimPrefix(line, "data: ")
		if payload == "[DONE]" {
			break
		}

		var delta string
		if c.isGoogleDirect() {
			var chunk generateResponse
			if err := json.Unmarshal([]byte(payload), &chunk); err != nil {
				continue
			}
			if chunk.Error != nil && chunk.Error.Message != "" {
				return "", fmt.Errorf("llm: api error: %s", chunk.Error.Message)
			}
			delta = extractGeminiText(chunk)
		} else {
			var chunk chatCompletionResponse
			if err := json.Unmarshal([]byte(payload), &chunk); err != nil {
				continue
			}
			if chunk.Error != nil && chunk.Error.Message != "" {
				return "", fmt.Errorf("llm: api error: %s", chunk.Error.Message)
			}
			delta = extractChatDelta(chunk)
		}

		if delta == "" {
			continue
		}
		full.WriteString(delta)
		if onDelta != nil {
			if err := onDelta(delta); err != nil {
				return full.String(), err
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return full.String(), err
	}
	text := strings.TrimSpace(full.String())
	if text == "" {
		return "", fmt.Errorf("llm: empty stream response")
	}
	return text, nil
}

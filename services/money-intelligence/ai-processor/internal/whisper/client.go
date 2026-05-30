package whisper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

const maxAudioBytes = 10 << 20

type apiMode int

const (
	modeOpenAI apiMode = iota
	modeASR
)

// Client — HTTP-клиент транскрипции (OpenAI или whisper-asr-webservice /asr).
type Client struct {
	url        string
	apiKey     string
	model      string
	mode       apiMode
	httpClient *http.Client
}

// NewClient создаёт клиент транскрипции.
func NewClient(url, apiKey, model string) *Client {
	if model == "" {
		model = "whisper-1"
	}
	url = strings.TrimSpace(url)
	mode := modeASR
	if strings.Contains(url, "/v1/audio/transcriptions") {
		mode = modeOpenAI
	}
	return &Client{
		url:    url,
		apiKey: strings.TrimSpace(apiKey),
		model:  model,
		mode:   mode,
		httpClient: &http.Client{
			Timeout: 90 * time.Second,
		},
	}
}

// Enabled возвращает true, если задан URL сервиса.
func (c *Client) Enabled() bool {
	return c != nil && c.url != ""
}

func (c *Client) requestURL() string {
	if c.mode == modeOpenAI {
		return c.url
	}
	if strings.Contains(c.url, "/asr") {
		return c.url
	}
	base := strings.TrimRight(c.url, "/")
	return base + "/asr?output=json&language=ru&task=transcribe"
}

// Transcribe отправляет аудио и возвращает текст.
func (c *Client) Transcribe(ctx context.Context, filename string, audio []byte) (string, error) {
	if !c.Enabled() {
		return "", fmt.Errorf("whisper: url not configured")
	}
	if len(audio) == 0 {
		return "", fmt.Errorf("whisper: empty audio")
	}
	if len(audio) > maxAudioBytes {
		return "", fmt.Errorf("whisper: file too large")
	}

	fieldName := "file"
	if c.mode == modeASR {
		fieldName = "audio_file"
	}

	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	part, err := w.CreateFormFile(fieldName, filename)
	if err != nil {
		return "", fmt.Errorf("whisper: form file: %w", err)
	}
	if _, err := part.Write(audio); err != nil {
		return "", fmt.Errorf("whisper: write audio: %w", err)
	}
	if c.mode == modeOpenAI {
		_ = w.WriteField("model", c.model)
		_ = w.WriteField("language", "ru")
	}
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("whisper: close form: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.requestURL(), &body)
	if err != nil {
		return "", fmt.Errorf("whisper: request: %w", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("whisper: http: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", fmt.Errorf("whisper: read: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("whisper: status %d: %s", resp.StatusCode, string(raw))
	}

	text, err := parseTranscript(raw)
	if err != nil {
		return "", err
	}
	if text == "" {
		return "", fmt.Errorf("whisper: empty transcript")
	}
	return text, nil
}

func parseTranscript(raw []byte) (string, error) {
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" {
		return "", nil
	}
	if strings.HasPrefix(trimmed, "{") {
		var parsed struct {
			Text string `json:"text"`
		}
		if err := json.Unmarshal(raw, &parsed); err != nil {
			return "", fmt.Errorf("whisper: decode: %w", err)
		}
		return strings.TrimSpace(parsed.Text), nil
	}
	return trimmed, nil
}

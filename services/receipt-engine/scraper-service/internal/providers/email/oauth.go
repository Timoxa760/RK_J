package email

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/mailru"
	"golang.org/x/oauth2/yandex"
)

type Provider string

const (
	ProviderYandex Provider = "yandex"
	ProviderMail   Provider = "mailru"
)

type OAuthManager struct {
	mu       sync.RWMutex
	configs  map[Provider]*oauth2.Config
	tokens   map[string]*oauth2.Token
	httpCli  *http.Client
	demoMode bool
}

func NewOAuthManager() *OAuthManager {
	demoMode := os.Getenv("DEMO_MODE") == "true"

	m := &OAuthManager{
		configs:  make(map[Provider]*oauth2.Config),
		tokens:   make(map[string]*oauth2.Token),
		httpCli:  &http.Client{Timeout: 10 * time.Second},
		demoMode: demoMode,
	}

	if !demoMode {
		m.configs[ProviderYandex] = &oauth2.Config{
			ClientID:     os.Getenv("OAUTH_YANDEX_CLIENT_ID"),
			ClientSecret: os.Getenv("OAUTH_YANDEX_CLIENT_SECRET"),
			Endpoint:     yandex.Endpoint,
			RedirectURL:  os.Getenv("OAUTH_YANDEX_REDIRECT_URL"),
			Scopes:       []string{"mail:imap_read"},
		}

		m.configs[ProviderMail] = &oauth2.Config{
			ClientID:     os.Getenv("OAUTH_MAIL_CLIENT_ID"),
			ClientSecret: os.Getenv("OAUTH_MAIL_CLIENT_SECRET"),
			Endpoint:     mailru.Endpoint,
			RedirectURL:  os.Getenv("OAUTH_MAIL_REDIRECT_URL"),
			Scopes:       []string{"mail.imap"},
		}
	}

	return m
}

func (m *OAuthManager) GetAuthURL(provider Provider, state string) string {
	if m.demoMode {
		return fmt.Sprintf("/demo/oauth/%s?state=%s", provider, state)
	}

	cfg, ok := m.configs[provider]
	if !ok {
		return ""
	}
	return cfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (m *OAuthManager) ExchangeCode(ctx context.Context, provider Provider, code string) (*oauth2.Token, error) {
	if m.demoMode {
		tok := &oauth2.Token{
			AccessToken:  fmt.Sprintf("demo-access-token-%s-%s", provider, randomHex(8)),
			TokenType:    "Bearer",
			Expiry:       time.Now().Add(24 * time.Hour),
			RefreshToken: fmt.Sprintf("demo-refresh-%s", randomHex(8)),
		}
		return tok, nil
	}

	cfg, ok := m.configs[provider]
	if !ok {
		return nil, fmt.Errorf("oauth: unknown provider %s", provider)
	}

	tok, err := cfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("oauth: exchange code: %w", err)
	}
	return tok, nil
}

func (m *OAuthManager) GetToken(ctx context.Context, provider Provider, userID string) (*oauth2.Token, error) {
	key := tokenKey(provider, userID)

	m.mu.RLock()
	tok, ok := m.tokens[key]
	m.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("oauth: no token for %s/%s", provider, userID)
	}

	if !tok.Valid() && tok.RefreshToken != "" && !m.demoMode {
		cfg := m.configs[provider]
		src := cfg.TokenSource(ctx, tok)
		newTok, err := src.Token()
		if err != nil {
			return nil, fmt.Errorf("oauth: refresh token: %w", err)
		}
		m.mu.Lock()
		m.tokens[key] = newTok
		m.mu.Unlock()
		return newTok, nil
	}

	return tok, nil
}

func (m *OAuthManager) StoreToken(provider Provider, userID string, tok *oauth2.Token) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tokens[tokenKey(provider, userID)] = tok
}

func (m *OAuthManager) Client(ctx context.Context, provider Provider, userID string) (*http.Client, error) {
	if m.demoMode {
		return m.httpCli, nil
	}

	tok, err := m.GetToken(ctx, provider, userID)
	if err != nil {
		return nil, err
	}

	cfg, ok := m.configs[provider]
	if !ok {
		return nil, fmt.Errorf("oauth: unknown provider %s", provider)
	}

	return cfg.Client(ctx, tok), nil
}

func (m *OAuthManager) SaveTokens(w io.Writer) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return json.NewEncoder(w).Encode(m.tokens)
}

func (m *OAuthManager) LoadTokens(r io.Reader) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return json.NewDecoder(r).Decode(&m.tokens)
}

type OAuthCallbackHandler struct {
	manager *OAuthManager
}

func NewOAuthCallbackHandler(manager *OAuthManager) *OAuthCallbackHandler {
	return &OAuthCallbackHandler{manager: manager}
}

func (h *OAuthCallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	provider := Provider(r.URL.Query().Get("provider"))
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if provider == "" || code == "" {
		http.Error(w, "missing provider or code", http.StatusBadRequest)
		return
	}

	tok, err := h.manager.ExchangeCode(r.Context(), provider, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.manager.StoreToken(provider, state, tok)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":   "ok",
		"provider": string(provider),
	})
}

func tokenKey(provider Provider, userID string) string {
	return string(provider) + ":" + userID
}

func randomHex(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}

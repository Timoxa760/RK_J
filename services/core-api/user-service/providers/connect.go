package providers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

type ConnectRequest struct {
	UserID      string `json:"user_id"`
	Provider    string `json:"provider"`
	Credentials map[string]string `json:"credentials"`
}

// ConnectResponse — 200 OK по API_Contract.
type ConnectResponse struct {
	Message  string `json:"message"`
	Provider string `json:"provider"`
	Status   string `json:"status"`
}

type ConnectedProvider struct {
	UserID      string
	Provider    string
	Credentials string
}

var (
	mu        sync.Mutex
	providers = make(map[string]ConnectedProvider)
)

type ConnectHandler struct {
	demoMode bool
}

func NewConnectHandler(demoMode bool) *ConnectHandler {
	return &ConnectHandler{demoMode: demoMode}
}

func (h *ConnectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req ConnectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}

	provider := r.URL.Query().Get("provider")
	if provider == "" {
		provider = req.Provider
	}
	if provider == "" || req.Credentials == nil {
		http.Error(w, `{"error":"provider and credentials required"}`, http.StatusBadRequest)
		return
	}

	userID := req.UserID
	if userID == "" {
		userID = req.Credentials["phone"]
	}
	if userID == "" {
		http.Error(w, `{"error":"user_id or credentials.phone required"}`, http.StatusBadRequest)
		return
	}

	mu.Lock()
	conflictKey := userID + ":" + provider
	if _, exists := providers[conflictKey]; exists {
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": "provider already connected"})
		return
	}
	mu.Unlock()

	credsJSON, err := json.Marshal(req.Credentials)
	if err != nil {
		http.Error(w, `{"error":"marshal credentials"}`, http.StatusInternalServerError)
		return
	}

	encrypted, err := encrypt(credsJSON, getEncryptionKey())
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"encrypt: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	key := userID + ":" + provider
	mu.Lock()
	providers[key] = ConnectedProvider{
		UserID:      userID,
		Provider:    provider,
		Credentials: encrypted,
	}
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ConnectResponse{
		Message:  "Provider connected",
		Provider: provider,
		Status:   "active",
	})
}

func getEncryptionKey() []byte {
	if s := os.Getenv("ENCRYPTION_KEY"); s != "" {
		return deriveKey([]byte(s))
	}
	return deriveKey([]byte("demo-key-2026"))
}

func deriveKey(key []byte) []byte {
	if len(key) >= 32 {
		return key[:32]
	}
	padded := make([]byte, 32)
	copy(padded, key)
	return padded
}

func encrypt(plaintext []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("aes: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("gcm: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return hex.EncodeToString(ciphertext), nil
}

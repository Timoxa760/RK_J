package advisor

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"backend_project/internal/auth"
)

type ChatAction struct {
	Type          string `json:"type"`
	Label         string `json:"label"`
	Path          string `json:"path,omitempty"`
	Hash          string `json:"hash,omitempty"`
	Ask           string `json:"ask,omitempty"`
	ProfileField  string `json:"profileField,omitempty"`
}

type StoredChatMessage struct {
	ID        string       `json:"id"`
	Role      string       `json:"role"`
	Content   string       `json:"content"`
	Actions   []ChatAction `json:"actions,omitempty"`
	Source    string       `json:"source,omitempty"`
	CreatedAt int64        `json:"created_at"`
}

type ChatHistoryResponse struct {
	Messages []StoredChatMessage `json:"messages"`
}

type ChatStore struct {
	pool *pgxpool.Pool
}

func NewChatStore(pool *pgxpool.Pool) *ChatStore {
	if pool == nil {
		return nil
	}
	return &ChatStore{pool: pool}
}

func (s *ChatStore) Enabled() bool {
	return s != nil && s.pool != nil
}

func (s *ChatStore) EnsureSchema(ctx context.Context) error {
	if !s.Enabled() {
		return nil
	}
	_, err := s.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS advisor_chat_messages (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id VARCHAR(64) NOT NULL,
			role VARCHAR(16) NOT NULL CHECK (role IN ('user', 'assistant')),
			content TEXT NOT NULL,
			actions JSONB,
			source VARCHAR(16),
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_advisor_chat_user_created
			ON advisor_chat_messages (user_id, created_at DESC);
	`)
	if err != nil {
		return fmt.Errorf("ensure advisor_chat_messages: %w", err)
	}
	return nil
}

func (s *ChatStore) ListHistory(ctx context.Context, userID string, limit int) ([]StoredChatMessage, error) {
	if !s.Enabled() {
		return nil, nil
	}
	if limit <= 0 {
		limit = 50
	}
	userID = auth.NormalizePhone(userID)
	rows, err := s.pool.Query(ctx, `
		SELECT id, role, content, actions, source, created_at
		FROM (
			SELECT id, role, content, actions, source, created_at
			FROM advisor_chat_messages
			WHERE user_id = $1
			ORDER BY created_at DESC
			LIMIT $2
		) recent
		ORDER BY created_at ASC
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []StoredChatMessage
	for rows.Next() {
		var msg StoredChatMessage
		var actionsRaw []byte
		var createdAt time.Time
		if err := rows.Scan(&msg.ID, &msg.Role, &msg.Content, &actionsRaw, &msg.Source, &createdAt); err != nil {
			return nil, err
		}
		msg.CreatedAt = createdAt.UnixMilli()
		if len(actionsRaw) > 0 {
			_ = json.Unmarshal(actionsRaw, &msg.Actions)
		}
		out = append(out, msg)
	}
	return out, rows.Err()
}

func (s *ChatStore) ClearHistory(ctx context.Context, userID string) error {
	if !s.Enabled() {
		return nil
	}
	userID = auth.NormalizePhone(userID)
	_, err := s.pool.Exec(ctx, `DELETE FROM advisor_chat_messages WHERE user_id = $1`, userID)
	return err
}

func (s *ChatStore) Append(ctx context.Context, userID, role, content, source string, actions []ChatAction) (string, error) {
	if !s.Enabled() {
		return "", nil
	}
	userID = auth.NormalizePhone(userID)
	var actionsJSON []byte
	if len(actions) > 0 {
		var err error
		actionsJSON, err = json.Marshal(actions)
		if err != nil {
			return "", err
		}
	}
	var id string
	err := s.pool.QueryRow(ctx, `
		INSERT INTO advisor_chat_messages (user_id, role, content, actions, source)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, userID, role, content, actionsJSON, source).Scan(&id)
	return id, err
}

func (s *ChatStore) RecentTurns(ctx context.Context, userID string, limit int) ([]ChatMessage, error) {
	msgs, err := s.ListHistory(ctx, userID, limit)
	if err != nil {
		return nil, err
	}
	out := make([]ChatMessage, 0, len(msgs))
	for _, msg := range msgs {
		out = append(out, ChatMessage{Role: msg.Role, Content: msg.Content})
	}
	return out, nil
}

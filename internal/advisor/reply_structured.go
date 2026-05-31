package advisor

import (
	"encoding/json"
	"regexp"
	"strings"
)

// ReplyBlock — один блок ответа советника для UI.
type ReplyBlock struct {
	Type  string   `json:"type"`
	Text  string   `json:"text,omitempty"`
	Items []string `json:"items,omitempty"`
	Tone  string   `json:"tone,omitempty"`
}

// StructuredReply — JSON-ответ LLM.
type StructuredReply struct {
	Title  string       `json:"title,omitempty"`
	Blocks []ReplyBlock `json:"blocks"`
}

type storedReplyEnvelope struct {
	V      int          `json:"v"`
	Title  string       `json:"title,omitempty"`
	Blocks []ReplyBlock `json:"blocks"`
	Plain  string       `json:"plain"`
}

var reJSONFence = regexp.MustCompile("(?s)```(?:json)?\\s*(\\{.*\\})\\s*```")

var allowedBlockTypes = map[string]struct{}{
	"lead":      {},
	"heading":   {},
	"paragraph": {},
	"list":      {},
	"callout":   {},
}

// ParseStructuredReply разбирает JSON от LLM; при ошибке — один paragraph из текста.
func ParseStructuredReply(raw string) (StructuredReply, string) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return StructuredReply{}, ""
	}

	jsonText := extractJSONObject(raw)
	var parsed StructuredReply
	if err := json.Unmarshal([]byte(jsonText), &parsed); err == nil {
		parsed = normalizeStructuredReply(parsed)
		if len(parsed.Blocks) > 0 {
			plain := StructuredToPlainText(parsed)
			return parsed, plain
		}
	}

	fallback := FormatAdvisorReply(raw)
	if fallback == "" {
		fallback = raw
	}
	return StructuredReply{
		Blocks: []ReplyBlock{{Type: "paragraph", Text: fallback}},
	}, fallback
}

func extractJSONObject(raw string) string {
	raw = strings.TrimSpace(raw)
	if strings.HasPrefix(raw, "{") {
		return raw
	}
	if m := reJSONFence.FindStringSubmatch(raw); len(m) > 1 {
		return strings.TrimSpace(m[1])
	}
	return raw
}

func normalizeStructuredReply(in StructuredReply) StructuredReply {
	out := StructuredReply{Title: repairStructuredText(strings.TrimSpace(in.Title))}
	for _, b := range in.Blocks {
		typ := strings.ToLower(strings.TrimSpace(b.Type))
		if _, ok := allowedBlockTypes[typ]; !ok {
			continue
		}
		block := ReplyBlock{Type: typ, Tone: strings.TrimSpace(b.Tone)}
		block.Text = repairStructuredText(strings.TrimSpace(b.Text))
		for _, item := range b.Items {
			item = repairStructuredText(strings.TrimSpace(item))
			if item != "" {
				block.Items = append(block.Items, item)
			}
		}
		if typ == "list" && len(block.Items) == 0 {
			continue
		}
		if typ != "list" && block.Text == "" {
			continue
		}
		out.Blocks = append(out.Blocks, block)
	}
	return out
}

// StructuredToPlainText — fallback для старых клиентов и уведомлений.
func StructuredToPlainText(s StructuredReply) string {
	var parts []string
	if t := strings.TrimSpace(s.Title); t != "" {
		parts = append(parts, t)
	}
	for _, b := range s.Blocks {
		switch b.Type {
		case "list":
			for _, item := range b.Items {
				parts = append(parts, "• "+item)
			}
		default:
			if b.Text != "" {
				parts = append(parts, b.Text)
			}
		}
	}
	return strings.TrimSpace(strings.Join(parts, "\n\n"))
}

// EncodeStoredContent сохраняет структурированный ответ в БД.
func EncodeStoredContent(s StructuredReply, plain string) string {
	env := storedReplyEnvelope{
		V:      1,
		Title:  s.Title,
		Blocks: s.Blocks,
		Plain:  plain,
	}
	b, err := json.Marshal(env)
	if err != nil {
		return plain
	}
	return string(b)
}

// DecodeStoredContent читает content из истории (JSON v1 или plain text).
func DecodeStoredContent(content string) (StructuredReply, string) {
	content = strings.TrimSpace(content)
	if content == "" {
		return StructuredReply{}, ""
	}
	var env storedReplyEnvelope
	if err := json.Unmarshal([]byte(content), &env); err == nil && env.V == 1 && len(env.Blocks) > 0 {
		s := StructuredReply{Title: env.Title, Blocks: env.Blocks}
		plain := env.Plain
		if plain == "" {
			plain = StructuredToPlainText(s)
		}
		return s, plain
	}
	return StructuredReply{
		Blocks: []ReplyBlock{{Type: "paragraph", Text: content}},
	}, content
}

func structuredFromPlain(text string) StructuredReply {
	text = strings.TrimSpace(text)
	if text == "" {
		return StructuredReply{}
	}
	return StructuredReply{
		Blocks: []ReplyBlock{{Type: "paragraph", Text: text}},
	}
}

package advisor

import (
	"strings"
	"testing"
)

func TestParseStructuredReply_JSON(t *testing.T) {
	raw := `{"title":"Зачем расходы","blocks":[{"type":"lead","text":"Коротко."},{"type":"heading","text":"Что указать"},{"type":"list","items":["аренда","продукты"]}]}`
	got, plain := ParseStructuredReply(raw)
	if got.Title != "Зачем расходы" {
		t.Fatalf("title: %q", got.Title)
	}
	if len(got.Blocks) != 3 {
		t.Fatalf("blocks: %d", len(got.Blocks))
	}
	if !strings.Contains(plain, "• аренда") {
		t.Fatalf("plain: %q", plain)
	}
}

func TestParseStructuredReply_PreservesWordsInFallback(t *testing.T) {
	raw := "расходы обязательные указать транспорт"
	got, _ := ParseStructuredReply(raw)
	if len(got.Blocks) != 1 {
		t.Fatalf("expected 1 block")
	}
	if strings.Contains(got.Blocks[0].Text, "расход ы") {
		t.Fatalf("broke word: %q", got.Blocks[0].Text)
	}
}

func TestDecodeStoredContent_Envelope(t *testing.T) {
	stored := EncodeStoredContent(StructuredReply{
		Title:  "T",
		Blocks: []ReplyBlock{{Type: "paragraph", Text: "body"}},
	}, "plain")
	s, plain := DecodeStoredContent(stored)
	if s.Title != "T" || plain != "plain" {
		t.Fatalf("got title=%q plain=%q", s.Title, plain)
	}
}

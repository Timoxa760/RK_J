package expense

import (
	"context"
	"testing"
)

func TestParser_RegexSingle(t *testing.T) {
	p := NewParser(nil)
	res := p.Parse(context.Background(), ParseInput{RawText: "продукты 5000"})
	if len(res.Expenses) != 1 || res.Expenses[0].Amount != 5000 {
		t.Fatalf("unexpected %+v", res)
	}
	if res.ParsedBy != "regex" || !res.Parsed {
		t.Fatalf("parsed flags: %+v", res)
	}
}

func TestParser_ExplicitAmount(t *testing.T) {
	p := NewParser(nil)
	res := p.Parse(context.Background(), ParseInput{Amount: 1200, Category: "Транспорт"})
	if len(res.Expenses) != 1 || res.Expenses[0].Amount != 1200 {
		t.Fatalf("unexpected %+v", res)
	}
}

func TestParser_Empty(t *testing.T) {
	p := NewParser(nil)
	res := p.Parse(context.Background(), ParseInput{})
	if len(res.Expenses) != 0 {
		t.Fatalf("expected empty, got %+v", res)
	}
}

func TestExtractJSON_Fence(t *testing.T) {
	got := extractJSON("```json\n{\"expenses\":[]}\n```")
	if got != `{"expenses":[]}` {
		t.Fatalf("got %q", got)
	}
}

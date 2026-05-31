package expense

import (
	"context"
	"testing"
)

func TestParser_MultipleExpenses(t *testing.T) {
	p := NewParser(nil)
	res := p.Parse(context.Background(), ParseInput{RawText: "кофе 300 рублей и такси 500 рублей"})
	if len(res.Expenses) != 2 {
		t.Fatalf("unexpected %+v", res)
	}
	if res.Expenses[0].Category != "Кафе и рестораны" || res.Expenses[1].Category != "Транспорт" {
		t.Fatalf("categories: %+v", res.Expenses)
	}
}

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

func TestParser_Kolbasa(t *testing.T) {
	p := NewParser(nil)
	res := p.Parse(context.Background(), ParseInput{RawText: "Колбаса 300 руб"})
	if len(res.Expenses) != 1 || res.Expenses[0].Amount != 300 {
		t.Fatalf("unexpected %+v", res)
	}
	if res.Expenses[0].Category != "Продукты" {
		t.Fatalf("category=%q", res.Expenses[0].Category)
	}
	if res.Expenses[0].Description != "Колбаса" {
		t.Fatalf("description=%q", res.Expenses[0].Description)
	}
}

func TestParser_VoiceSlang(t *testing.T) {
	p := NewParser(nil)
	res := p.Parse(context.Background(), ParseInput{RawText: "вышел с пятерочки, оставил 10 тыщ"})
	if len(res.Expenses) != 1 || res.Expenses[0].Amount != 10000 {
		t.Fatalf("unexpected %+v", res)
	}
	if res.Expenses[0].Category != "Продукты" {
		t.Fatalf("category=%q", res.Expenses[0].Category)
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

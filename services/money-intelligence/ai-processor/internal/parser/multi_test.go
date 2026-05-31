package parser

import "testing"

func TestParseAll_MultipleCategories(t *testing.T) {
	got := ParseAll("кофе 300 рублей и такси 500 рублей")
	if len(got) != 2 {
		t.Fatalf("len=%d want 2: %+v", len(got), got)
	}
	if got[0].Amount != 300 || got[0].Category != "Кафе и рестораны" {
		t.Fatalf("first=%+v", got[0])
	}
	if got[1].Amount != 500 || got[1].Category != "Транспорт" {
		t.Fatalf("second=%+v", got[1])
	}
}

func TestParseAll_CommaSeparated(t *testing.T) {
	got := ParseAll("колбаса 300 руб, такси 500 руб")
	if len(got) != 2 {
		t.Fatalf("len=%d want 2: %+v", len(got), got)
	}
	if got[0].Amount != 300 || got[0].Category != "Продукты" {
		t.Fatalf("first=%+v", got[0])
	}
	if got[1].Amount != 500 || got[1].Category != "Транспорт" {
		t.Fatalf("second=%+v", got[1])
	}
}

func TestParseAll_SingleCombinedProducts(t *testing.T) {
	got := ParseAll("колбаса сыр 500 руб")
	if len(got) != 1 {
		t.Fatalf("len=%d want 1: %+v", len(got), got)
	}
	if got[0].Amount != 500 || got[0].Description != "Колбаса, Сыр" {
		t.Fatalf("got=%+v", got[0])
	}
}

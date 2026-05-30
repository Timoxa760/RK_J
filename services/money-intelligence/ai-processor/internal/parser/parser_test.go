package parser

import "testing"

func TestParse_TyshSlang(t *testing.T) {
	p := Parse("вышел с пятерочки, оставил 10 тыщ")
	if p == nil {
		t.Fatal("expected parse result")
	}
	if p.Amount != 10000 {
		t.Fatalf("amount=%v", p.Amount)
	}
	if p.Category != "Продукты" {
		t.Fatalf("category=%q", p.Category)
	}
	if p.Store != "Пятёрочка" {
		t.Fatalf("store=%q", p.Store)
	}
}

func TestParse_WhisperTypoStore(t *testing.T) {
	p := Parse("вышел спитёрочки, оставил 5 тысяч")
	if p == nil {
		t.Fatal("expected parse result")
	}
	if p.Amount != 5000 {
		t.Fatalf("amount=%v", p.Amount)
	}
	if p.Store != "Пятёрочка" {
		t.Fatalf("store=%q", p.Store)
	}
	if p.Description != "Пятёрочка" {
		t.Fatalf("desc=%q", p.Description)
	}
}

func TestParse_Kolbasa(t *testing.T) {
	p := Parse("Колбаса 300 руб")
	if p == nil {
		t.Fatal("expected parse result")
	}
	if p.Amount != 300 {
		t.Fatalf("amount=%v", p.Amount)
	}
	if p.Category != "Продукты" {
		t.Fatalf("category=%q", p.Category)
	}
	if p.Description != "Колбаса" {
		t.Fatalf("desc=%q", p.Description)
	}
}

func TestParse_Kosar(t *testing.T) {
	p := Parse("потратил 5 косарей на такси")
	if p == nil || p.Amount != 5000 || p.Category != "Транспорт" {
		t.Fatalf("unexpected %+v", p)
	}
}

func TestParse_PlainRubles(t *testing.T) {
	p := Parse("продукты 5000")
	if p == nil || p.Amount != 5000 {
		t.Fatalf("unexpected %+v", p)
	}
}

package email

import (
	"testing"
)

func TestParseReceiptHTML_HappyPath(t *testing.T) {
	html := `<html><body>
		<table>
			<tr><td>Магазин</td><td>Пятёрочка</td></tr>
			<tr><td>Дата</td><td>15.05.2026</td></tr>
			<tr><td>Сумма</td><td>1 032.50 ₽</td></tr>
			<tr><td>Молоко 3.2%</td><td>78.00 ₽</td><td>2</td></tr>
			<tr><td>Хлеб белый</td><td>45.00 ₽</td><td>1</td></tr>
			<tr><td>Сыр Российский</td><td>189.00 ₽</td><td>1</td></tr>
		</table>
	</body></html>`

	p := NewParser()
	r := p.ParseReceiptHTML(html)
	if r == nil {
		t.Fatal("ParseReceiptHTML returned nil")
	}

	if r.Store != "Пятёрочка" {
		t.Errorf("expected store Пятёрочка, got %s", r.Store)
	}

	if r.Total != 1032.50 {
		t.Errorf("expected total 1032.50, got %f", r.Total)
	}

	if len(r.Items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(r.Items))
	}

	expectedItems := []struct {
		name     string
		price    float64
		quantity int
	}{
		{"Молоко 3.2%", 78.00, 2},
		{"Хлеб белый", 45.00, 1},
		{"Сыр Российский", 189.00, 1},
	}

	for i, exp := range expectedItems {
		got := r.Items[i]
		if got.Name != exp.name {
			t.Errorf("item %d name: expected %s, got %s", i, exp.name, got.Name)
		}
		if got.Price != exp.price {
			t.Errorf("item %d price: expected %f, got %f", i, exp.price, got.Price)
		}
		if got.Quantity != exp.quantity {
			t.Errorf("item %d quantity: expected %d, got %d", i, exp.quantity, got.Quantity)
		}
	}
}

func TestParseReceiptHTML_WithoutTable(t *testing.T) {
	html := `<html><body>
		<div>Магазин: Пятёрочка</div>
		<div>Дата: 15.05.2026</div>
		<div>Сумма: 1032.50</div>
	</body></html>`

	p := NewParser()
	r := p.ParseReceiptHTML(html)
	if r == nil {
		t.Fatal("ParseReceiptHTML returned nil for div-only HTML")
	}

	if r.Store != "Пятёрочка" {
		t.Errorf("expected Пятёрочка, got %s", r.Store)
	}
}

func TestParseReceiptHTML_EmptyHTML(t *testing.T) {
	p := NewParser()
	r := p.ParseReceiptHTML("<html></html>")
	if r == nil {
		t.Fatal("ParseReceiptHTML returned nil")
	}
}

func TestParseReceiptHTML_MagnitFormat(t *testing.T) {
	html := `<html><body>
		<table class="receipt">
			<tr><th>Наименование</th><th>Цена</th><th>Кол-во</th></tr>
			<tr><td>Молоко</td><td>78.00</td><td>1</td></tr>
			<tr><td>Творог</td><td>120.00</td><td>2</td></tr>
			<tr><td>Хлеб</td><td>35.00</td><td>1</td></tr>
		</table>
		<p>Магазин: Магнит</p>
		<p>Дата: 2026-05-15</p>
		<p>Итого: 353.00 ₽</p>
	</body></html>`

	p := NewParser()
	r := p.ParseReceiptHTML(html)
	if r == nil {
		t.Fatal("ParseReceiptHTML returned nil")
	}

	if r.Store != "Магнит" {
		t.Errorf("expected Магнит, got %s", r.Store)
	}

	if len(r.Items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(r.Items))
	}
}

func TestParseReceiptHTML_NoItems(t *testing.T) {
	html := `<html><body>
		<table>
			<tr><td>Магазин</td><td>Тест</td></tr>
			<tr><td>Сумма</td><td>0.00 ₽</td></tr>
		</table>
	</body></html>`

	p := NewParser()
	r := p.ParseReceiptHTML(html)
	if r == nil {
		t.Fatal("ParseReceiptHTML returned nil")
	}
}

func TestParsePrice(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"78.00", 78.00},
		{"1 032.50", 1032.50},
		{"1 032,50", 1032.50},
		{"1032.50 ₽", 1032.50},
		{"1 200", 1200.00},
		{"0.00", 0.00},
		{"", 0},
		{"abc", 0},
	}

	for _, tt := range tests {
		got := parsePrice(tt.input)
		if got != tt.want {
			t.Errorf("parsePrice(%q) = %f, want %f", tt.input, got, tt.want)
		}
	}
}

func TestParseDate(t *testing.T) {
	tests := []struct {
		input string
		want  string // expected output in YYYY-MM-DD
	}{
		{"15.05.2026", "2026-05-15"},
		{"2026-05-15", "2026-05-15"},
		{"15/05/2026", "2026-05-15"},
	}

	for _, tt := range tests {
		got, err := parseDate(tt.input)
		if err != nil {
			t.Errorf("parseDate(%q) error: %v", tt.input, err)
			continue
		}
		gotStr := got.Format("2006-01-02")
		if gotStr != tt.want {
			t.Errorf("parseDate(%q) = %s, want %s", tt.input, gotStr, tt.want)
		}
	}
}

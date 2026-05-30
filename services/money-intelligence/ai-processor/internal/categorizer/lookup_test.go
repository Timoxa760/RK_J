package categorizer

import "testing"

func TestCategoryForText_Kolbasa(t *testing.T) {
	if got := CategoryForText("Колбаса 300 руб"); got != "Продукты" {
		t.Fatalf("category=%q", got)
	}
}

func TestProductLabelFromText_Kolbasa(t *testing.T) {
	if got := ProductLabelFromText("купил колбасу 300 руб"); got != "Колбас" {
		t.Fatalf("label=%q", got)
	}
}

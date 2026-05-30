package rublang

import "testing"

func TestExtractAll_Slang(t *testing.T) {
	cases := []struct {
		text string
		want []int
	}{
		{"оставил 10 тыщ", []int{10000}},
		{"10к на продукты", []int{10000}},
		{"180 тысяч в месяц", []int{180000}},
		{"5 косарей на аренду", []int{5000}},
		{"косарь на такси", []int{1000}},
		{"3 кеса", []int{3000}},
		{"сотка за кофе", []int{100}},
		{"2 сотки", []int{200}},
		{"150 000 зарплата и 20 000 с аренды", []int{150000, 20000}},
		{"вышел с пятерочки, оставил 10 тыщ", []int{10000}},
		{"5000 рублей", []int{5000}},
	}
	for _, tc := range cases {
		got := ExtractAll(tc.text)
		if len(got) != len(tc.want) {
			t.Fatalf("%q: got %v want %v", tc.text, got, tc.want)
		}
		for i := range tc.want {
			if got[i] != tc.want[i] {
				t.Fatalf("%q: got %v want %v", tc.text, got, tc.want)
			}
		}
	}
}

func TestExtractPrimary(t *testing.T) {
	v, ok := ExtractPrimary("продукты 500 рублей и кофе сотка")
	if !ok || v != 500 {
		t.Fatalf("got %v ok=%v", v, ok)
	}
	v, ok = ExtractPrimary("оставил 10 тыщ")
	if !ok || v != 10000 {
		t.Fatalf("got %v ok=%v", v, ok)
	}
}

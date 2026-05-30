package auth

import "testing"

func TestNormalizePhone(t *testing.T) {
	cases := map[string]string{
		"+79991234567":       "+79991234567",
		"+7 (123) 123-12-31": "+71231231231",
		"89991234567":        "+79991234567",
		"9991234567":         "+79991234567",
		"  +7 999 123 45 67 ": "+79991234567",
	}
	for in, want := range cases {
		if got := NormalizePhone(in); got != want {
			t.Errorf("NormalizePhone(%q) = %q, want %q", in, got, want)
		}
	}
}

package magnit

import "testing"

func TestINNDetector_ExactMatch(t *testing.T) {
	d := NewINNDetector()
	inn := d.Detect("Магнит")
	if inn != "2309085638" {
		t.Errorf("expected 2309085638 for Магнит, got %s", inn)
	}
}

func TestINNDetector_SubstringMatch(t *testing.T) {
	d := NewINNDetector()
	inn := d.Detect("Пятёрочка ул. Ленина")
	if inn != "7727563778" {
		t.Errorf("expected 7727563778 for Пятёрочка, got %s", inn)
	}
}

func TestINNDetector_UnknownStore(t *testing.T) {
	d := NewINNDetector()
	inn := d.Detect("Ашан ТРК")
	if inn != "2309085699" {
		t.Errorf("expected default INN 2309085699, got %s", inn)
	}
}

func TestINNDetector_CaseInsensitive(t *testing.T) {
	d := NewINNDetector()
	inn := d.Detect("МАГНИТ")
	if inn != "2309085638" {
		t.Errorf("expected 2309085638 for uppercase МАГНИТ, got %s", inn)
	}
}

func TestINNDetector_EmptyString(t *testing.T) {
	d := NewINNDetector()
	inn := d.Detect("")
	if inn != "2309085699" {
		t.Errorf("expected default INN for empty string, got %s", inn)
	}
}

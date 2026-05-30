package manual

import (
	"strings"
	"testing"
)

func TestParseFailureMessage_WithProduct(t *testing.T) {
	msg := parseFailureMessage("колбаса")
	if !strings.Contains(msg, "Колбаса") {
		t.Fatalf("unexpected message: %q", msg)
	}
}

func TestParseFailureMessage_Empty(t *testing.T) {
	msg := parseFailureMessage("")
	if !strings.Contains(msg, "300") {
		t.Fatalf("unexpected message: %q", msg)
	}
}

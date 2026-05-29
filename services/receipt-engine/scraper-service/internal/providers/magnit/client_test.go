package magnit

import (
	"testing"
)

func TestLogin_DemoMode(t *testing.T) {
	c := NewClient(true)
	if err := c.Login("+79990000001", "pass"); err != nil {
		t.Fatalf("Login: %v", err)
	}
}

func TestGetReceipts_DemoMode(t *testing.T) {
	c := NewClient(true)
	if err := c.Login("+79990000001", "pass"); err != nil {
		t.Fatal(err)
	}

	receipts, pages, err := c.GetReceipts(1)
	if err != nil {
		t.Fatalf("GetReceipts: %v", err)
	}
	if len(receipts) == 0 {
		t.Fatal("expected receipts")
	}
	if pages < 1 {
		t.Errorf("expected at least 1 page, got %d", pages)
	}
}

func TestGetReceipts_Pagination(t *testing.T) {
	c := NewClient(true)
	c.Login("+79990000001", "pass")

	r1, pages, _ := c.GetReceipts(1)
	r2, _, _ := c.GetReceipts(2)

	if len(r1) == 0 {
		t.Error("page 1 should have receipts")
	}
	if len(r2) != 0 {
		t.Errorf("page 2 should be empty (only 5 demo receipts), got %d", len(r2))
	}
	if pages != 1 {
		t.Errorf("expected 1 page, got %d", pages)
	}
}

func TestGetReceipts_RequireLogin(t *testing.T) {
	c := NewClient(false)
	if c.token != "" {
		t.Error("token should be empty before login")
	}
}

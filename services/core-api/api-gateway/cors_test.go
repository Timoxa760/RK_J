package main

import (
	"net/http"
	"testing"
)

func TestStripUpstreamCORS(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{
			"Access-Control-Allow-Origin":  []string{"*"},
			"Access-Control-Allow-Methods": []string{"GET, POST"},
			"Content-Type":                 []string{"application/json"},
		},
	}

	if err := stripUpstreamCORS(resp); err != nil {
		t.Fatalf("stripUpstreamCORS: %v", err)
	}

	if got := resp.Header.Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want empty", got)
	}
	if got := resp.Header.Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", got)
	}
}

func TestParseAllowedOrigins(t *testing.T) {
	tests := []struct {
		raw  string
		want []string
	}{
		{"", nil},
		{"http://127.0.0.1:3000", []string{"http://127.0.0.1:3000"}},
		{
			"http://localhost:3000, https://app.example.com",
			[]string{"http://localhost:3000", "https://app.example.com"},
		},
	}

	for _, tc := range tests {
		got := parseAllowedOrigins(tc.raw)
		if len(got) != len(tc.want) {
			t.Fatalf("parseAllowedOrigins(%q) = %v, want %v", tc.raw, got, tc.want)
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Fatalf("parseAllowedOrigins(%q)[%d] = %q, want %q", tc.raw, i, got[i], tc.want[i])
			}
		}
	}
}

package main

import (
	"testing"
)

func TestToLocalhost(t *testing.T) {
	got := toLocalhost("http://receipt-service:8002")
	want := "http://127.0.0.1:8002"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestSvcURL_LocalDefault(t *testing.T) {
	t.Setenv("RECEIPT_SERVICE_URL", "")
	t.Setenv("DOCKER", "")
	got := svcURL("RECEIPT_SERVICE_URL", "http://receipt-service:8002")
	if got != "http://127.0.0.1:8002" {
		t.Fatalf("got %q", got)
	}
}

func TestSvcURL_EnvOverride(t *testing.T) {
	t.Setenv("RECEIPT_SERVICE_URL", "http://custom:9999")
	got := svcURL("RECEIPT_SERVICE_URL", "http://receipt-service:8002")
	if got != "http://custom:9999" {
		t.Fatalf("got %q", got)
	}
}

func TestSvcURL_DockerDefault(t *testing.T) {
	t.Setenv("RECEIPT_SERVICE_URL", "")
	t.Setenv("DOCKER", "1")
	got := svcURL("RECEIPT_SERVICE_URL", "http://receipt-service:8002")
	if got != "http://receipt-service:8002" {
		t.Fatalf("got %q", got)
	}
}

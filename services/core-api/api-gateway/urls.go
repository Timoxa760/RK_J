package main

import (
	"net/url"
	"os"
)

// svcURL возвращает URL upstream-сервиса: env → docker (в контейнере) → 127.0.0.1 (локальные бинарники).
func svcURL(envKey, dockerDefault string) string {
	if u := os.Getenv(envKey); u != "" {
		return u
	}
	if runningInDocker() {
		return dockerDefault
	}
	return toLocalhost(dockerDefault)
}

func runningInDocker() bool {
	if os.Getenv("DOCKER") == "1" || os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		return true
	}
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	return false
}

func toLocalhost(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Port() == "" {
		return raw
	}
	u.Host = "127.0.0.1:" + u.Port()
	u.Scheme = "http"
	return u.String()
}

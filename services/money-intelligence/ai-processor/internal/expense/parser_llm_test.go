package expense

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend_project/internal/onlysq"
)

func TestParser_LLM(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"{\"expenses\":[{\"amount\":1332,\"category\":\"Продукты\",\"description\":\"продуктовый\"}],\"advice\":\"Чек чуть выше обычного.\"}"}}]}`))
	}))
	defer srv.Close()

	p := NewParser(onlysq.NewClient(srv.URL+"/v1", "key", "gpt-4o-mini"))
	res := p.Parse(context.Background(), ParseInput{RawText: "продуктовый 1332"})
	if res.ParsedBy != "onlysq" || len(res.Expenses) != 1 {
		t.Fatalf("unexpected %+v", res)
	}
	if res.Expenses[0].Amount != 1332 || res.Advice == "" {
		t.Fatalf("unexpected %+v", res)
	}
}

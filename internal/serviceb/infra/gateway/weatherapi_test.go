package gateway

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWeatherAPIFindCelsius(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if query.Get("key") != "secret" || query.Get("q") != "São Paulo" {
			t.Errorf("unexpected query %s", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"current":{"temp_c":28.5}}`))
	}))
	defer server.Close()

	celsius, err := NewWeatherAPI(server.URL, "secret", server.Client()).FindCelsius(context.Background(), "São Paulo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if celsius != 28.5 {
		t.Fatalf("expected 28.5, got %v", celsius)
	}
}

func TestWeatherAPIFindCelsiusUnexpectedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	if _, err := NewWeatherAPI(server.URL, "invalid", server.Client()).FindCelsius(context.Background(), "Recife"); err == nil {
		t.Fatal("expected error, got nil")
	}
}

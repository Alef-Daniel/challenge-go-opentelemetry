package gateway

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Alef-Daniel/challenge-go-opentelemetry/internal/serviceb/domain"
)

func TestViaCEPFindCity(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ws/01001000/json/" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"cep":"01001-000","localidade":"São Paulo","uf":"SP"}`))
	}))
	defer server.Close()

	city, err := NewViaCEP(server.URL, server.Client()).FindCity(context.Background(), "01001000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if city != "São Paulo" {
		t.Fatalf("expected São Paulo, got %s", city)
	}
}

func TestViaCEPFindCityNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"erro":"true"}`))
	}))
	defer server.Close()

	_, err := NewViaCEP(server.URL, server.Client()).FindCity(context.Background(), "99999999")
	if !errors.Is(err, domain.ErrZipcodeNotFound) {
		t.Fatalf("expected %v, got %v", domain.ErrZipcodeNotFound, err)
	}
}

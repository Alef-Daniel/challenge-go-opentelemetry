package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Alef-Daniel/challenge-go-opentelemetry/internal/servicea/domain"
	"github.com/Alef-Daniel/challenge-go-opentelemetry/internal/servicea/usecase"
)

type stubGateway struct {
	weather domain.Weather
	err     error
}

func (s stubGateway) FindByZipcode(context.Context, domain.Zipcode) (domain.Weather, error) {
	return s.weather, s.err
}

func TestWeatherHandler(t *testing.T) {
	weather := domain.Weather{City: "São Paulo", Celsius: 28.5, Fahrenheit: 83.3, Kelvin: 301.5}

	tests := []struct {
		name       string
		body       string
		gateway    stubGateway
		wantStatus int
		wantBody   string
	}{
		{
			name:       "valid zipcode",
			body:       `{"cep":"01001000"}`,
			gateway:    stubGateway{weather: weather},
			wantStatus: http.StatusOK,
			wantBody:   `{"city":"São Paulo","temp_C":28.5,"temp_F":83.3,"temp_K":301.5}`,
		},
		{name: "zipcode as number", body: `{"cep":1001000}`, wantStatus: http.StatusUnprocessableEntity, wantBody: "invalid zipcode"},
		{name: "zipcode too short", body: `{"cep":"0100100"}`, wantStatus: http.StatusUnprocessableEntity, wantBody: "invalid zipcode"},
		{name: "missing zipcode", body: `{}`, wantStatus: http.StatusUnprocessableEntity, wantBody: "invalid zipcode"},
		{name: "malformed body", body: `{`, wantStatus: http.StatusUnprocessableEntity, wantBody: "invalid zipcode"},
		{
			name:       "zipcode not found",
			body:       `{"cep":"99999999"}`,
			gateway:    stubGateway{err: domain.ErrZipcodeNotFound},
			wantStatus: http.StatusNotFound,
			wantBody:   "can not find zipcode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := NewRouter(NewWeatherHandler(usecase.NewGetWeather(tt.gateway)))
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}
			if got := strings.TrimSpace(rec.Body.String()); got != tt.wantBody {
				t.Fatalf("expected body %q, got %q", tt.wantBody, got)
			}
		})
	}
}

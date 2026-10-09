package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Alef-Daniel/challenge-go-opentelemetry/internal/servicea/domain"
)

type serviceBResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type ServiceB struct {
	baseURL string
	client  *http.Client
}

func NewServiceB(baseURL string, client *http.Client) *ServiceB {
	return &ServiceB{baseURL: baseURL, client: client}
}

func (g *ServiceB) FindByZipcode(ctx context.Context, zipcode domain.Zipcode) (domain.Weather, error) {
	endpoint := fmt.Sprintf("%s/weather/%s", g.baseURL, zipcode)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return domain.Weather{}, err
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return domain.Weather{}, fmt.Errorf("service b request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusUnprocessableEntity:
		return domain.Weather{}, domain.ErrInvalidZipcode
	case http.StatusNotFound:
		return domain.Weather{}, domain.ErrZipcodeNotFound
	default:
		return domain.Weather{}, fmt.Errorf("service b unexpected status: %d", resp.StatusCode)
	}

	var body serviceBResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return domain.Weather{}, fmt.Errorf("service b decode: %w", err)
	}

	return domain.Weather{
		City:       body.City,
		Celsius:    body.TempC,
		Fahrenheit: body.TempF,
		Kelvin:     body.TempK,
	}, nil
}

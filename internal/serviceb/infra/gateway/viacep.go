package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Alef-Daniel/challenge-go-opentelemetry/internal/serviceb/domain"
)

type viaCEPResponse struct {
	City string `json:"localidade"`
}

type ViaCEP struct {
	baseURL string
	client  *http.Client
}

func NewViaCEP(baseURL string, client *http.Client) *ViaCEP {
	return &ViaCEP{baseURL: baseURL, client: client}
}

func (g *ViaCEP) FindCity(ctx context.Context, zipcode domain.Zipcode) (string, error) {
	endpoint := fmt.Sprintf("%s/ws/%s/json/", g.baseURL, zipcode)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("viacep request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusBadRequest:
		return "", domain.ErrInvalidZipcode
	default:
		return "", fmt.Errorf("viacep unexpected status: %d", resp.StatusCode)
	}

	var body viaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("viacep decode: %w", err)
	}

	if body.City == "" {
		return "", domain.ErrZipcodeNotFound
	}

	return body.City, nil
}

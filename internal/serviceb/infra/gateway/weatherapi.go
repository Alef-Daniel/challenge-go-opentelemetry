package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type weatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

type WeatherAPI struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

func NewWeatherAPI(baseURL, apiKey string, client *http.Client) *WeatherAPI {
	return &WeatherAPI{baseURL: baseURL, apiKey: apiKey, client: client}
}

func (g *WeatherAPI) FindCelsius(ctx context.Context, city string) (float64, error) {
	query := url.Values{}
	query.Set("key", g.apiKey)
	query.Set("q", city)
	query.Set("aqi", "no")

	endpoint := fmt.Sprintf("%s/v1/current.json?%s", g.baseURL, query.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return 0, err
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("weatherapi request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("weatherapi unexpected status: %d", resp.StatusCode)
	}

	var body weatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return 0, fmt.Errorf("weatherapi decode: %w", err)
	}

	return body.Current.TempC, nil
}

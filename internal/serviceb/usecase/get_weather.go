package usecase

import (
	"context"

	"github.com/Alef-Daniel/challenge-go-opentelemetry/internal/serviceb/domain"
)

type LocationFinder interface {
	FindCity(ctx context.Context, zipcode domain.Zipcode) (string, error)
}

type TemperatureFinder interface {
	FindCelsius(ctx context.Context, city string) (float64, error)
}

type GetWeather struct {
	locations    LocationFinder
	temperatures TemperatureFinder
}

func NewGetWeather(locations LocationFinder, temperatures TemperatureFinder) *GetWeather {
	return &GetWeather{locations: locations, temperatures: temperatures}
}

func (uc *GetWeather) Execute(ctx context.Context, rawZipcode string) (domain.Weather, error) {
	zipcode, err := domain.NewZipcode(rawZipcode)
	if err != nil {
		return domain.Weather{}, err
	}

	city, err := uc.locations.FindCity(ctx, zipcode)
	if err != nil {
		return domain.Weather{}, err
	}

	celsius, err := uc.temperatures.FindCelsius(ctx, city)
	if err != nil {
		return domain.Weather{}, err
	}

	return domain.Weather{
		City:        city,
		Temperature: domain.NewTemperatureFromCelsius(celsius),
	}, nil
}

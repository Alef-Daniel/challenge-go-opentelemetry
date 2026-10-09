package usecase

import (
	"context"

	"github.com/Alef-Daniel/challenge-go-opentelemetry/internal/servicea/domain"
)

type WeatherGateway interface {
	FindByZipcode(ctx context.Context, zipcode domain.Zipcode) (domain.Weather, error)
}

type GetWeather struct {
	gateway WeatherGateway
}

func NewGetWeather(gateway WeatherGateway) *GetWeather {
	return &GetWeather{gateway: gateway}
}

func (uc *GetWeather) Execute(ctx context.Context, rawZipcode string) (domain.Weather, error) {
	zipcode, err := domain.NewZipcode(rawZipcode)
	if err != nil {
		return domain.Weather{}, err
	}
	return uc.gateway.FindByZipcode(ctx, zipcode)
}

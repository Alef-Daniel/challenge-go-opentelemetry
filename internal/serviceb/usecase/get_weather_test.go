package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/Alef-Daniel/challenge-go-opentelemetry/internal/serviceb/domain"
)

type stubLocationFinder struct {
	city string
	err  error
}

func (s stubLocationFinder) FindCity(context.Context, domain.Zipcode) (string, error) {
	return s.city, s.err
}

type stubTemperatureFinder struct {
	celsius float64
	err     error
	called  *bool
}

func (s stubTemperatureFinder) FindCelsius(context.Context, string) (float64, error) {
	if s.called != nil {
		*s.called = true
	}
	return s.celsius, s.err
}

func TestGetWeatherExecute(t *testing.T) {
	uc := NewGetWeather(stubLocationFinder{city: "São Paulo"}, stubTemperatureFinder{celsius: 28.5})

	weather, err := uc.Execute(context.Background(), "01001000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := domain.Weather{City: "São Paulo", Temperature: domain.NewTemperatureFromCelsius(28.5)}
	if weather != want {
		t.Fatalf("expected %+v, got %+v", want, weather)
	}
}

func TestGetWeatherExecuteInvalidZipcode(t *testing.T) {
	uc := NewGetWeather(stubLocationFinder{}, stubTemperatureFinder{})

	if _, err := uc.Execute(context.Background(), "123"); !errors.Is(err, domain.ErrInvalidZipcode) {
		t.Fatalf("expected %v, got %v", domain.ErrInvalidZipcode, err)
	}
}

func TestGetWeatherExecuteZipcodeNotFound(t *testing.T) {
	called := false
	uc := NewGetWeather(
		stubLocationFinder{err: domain.ErrZipcodeNotFound},
		stubTemperatureFinder{called: &called},
	)

	if _, err := uc.Execute(context.Background(), "99999999"); !errors.Is(err, domain.ErrZipcodeNotFound) {
		t.Fatalf("expected %v, got %v", domain.ErrZipcodeNotFound, err)
	}
	if called {
		t.Fatal("temperature finder must not be called when zipcode is not found")
	}
}

func TestGetWeatherExecuteTemperatureError(t *testing.T) {
	temperatureErr := errors.New("weather unavailable")
	uc := NewGetWeather(stubLocationFinder{city: "Recife"}, stubTemperatureFinder{err: temperatureErr})

	if _, err := uc.Execute(context.Background(), "50000000"); !errors.Is(err, temperatureErr) {
		t.Fatalf("expected %v, got %v", temperatureErr, err)
	}
}

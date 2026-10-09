package web

import (
	"context"
	"errors"
	"net/http"

	"github.com/Alef-Daniel/challenge-go-opentelemetry/internal/platform/httpresponse"
	"github.com/Alef-Daniel/challenge-go-opentelemetry/internal/serviceb/domain"
)

type WeatherFinder interface {
	Execute(ctx context.Context, rawZipcode string) (domain.Weather, error)
}

type weatherResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type WeatherHandler struct {
	finder WeatherFinder
}

func NewWeatherHandler(finder WeatherFinder) *WeatherHandler {
	return &WeatherHandler{finder: finder}
}

func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	weather, err := h.finder.Execute(r.Context(), r.PathValue("cep"))
	if err != nil {
		writeError(w, err)
		return
	}

	httpresponse.JSON(w, http.StatusOK, weatherResponse{
		City:  weather.City,
		TempC: weather.Temperature.Celsius,
		TempF: weather.Temperature.Fahrenheit,
		TempK: weather.Temperature.Kelvin,
	})
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidZipcode):
		httpresponse.Message(w, http.StatusUnprocessableEntity, domain.ErrInvalidZipcode.Error())
	case errors.Is(err, domain.ErrZipcodeNotFound):
		httpresponse.Message(w, http.StatusNotFound, domain.ErrZipcodeNotFound.Error())
	default:
		httpresponse.Message(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
}

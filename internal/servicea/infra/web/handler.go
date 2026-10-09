package web

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Alef-Daniel/challenge-go-opentelemetry/internal/platform/httpresponse"
	"github.com/Alef-Daniel/challenge-go-opentelemetry/internal/servicea/domain"
)

type WeatherFinder interface {
	Execute(ctx context.Context, rawZipcode string) (domain.Weather, error)
}

type weatherRequest struct {
	CEP *string `json:"cep"`
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
	var request weatherRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || request.CEP == nil {
		writeError(w, domain.ErrInvalidZipcode)
		return
	}

	weather, err := h.finder.Execute(r.Context(), *request.CEP)
	if err != nil {
		writeError(w, err)
		return
	}

	httpresponse.JSON(w, http.StatusOK, weatherResponse{
		City:  weather.City,
		TempC: weather.Celsius,
		TempF: weather.Fahrenheit,
		TempK: weather.Kelvin,
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

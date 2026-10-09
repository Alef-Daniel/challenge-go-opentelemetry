package web

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const getWeatherRoute = "GET /weather/{cep}"

func NewRouter(handler *WeatherHandler) http.Handler {
	mux := http.NewServeMux()
	mux.Handle(getWeatherRoute, otelhttp.NewHandler(http.HandlerFunc(handler.GetWeather), getWeatherRoute))
	return mux
}

package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Alef-Daniel/challenge-go-opentelemetry/internal/platform/env"
	"github.com/Alef-Daniel/challenge-go-opentelemetry/internal/platform/httpserver"
	"github.com/Alef-Daniel/challenge-go-opentelemetry/internal/platform/telemetry"
	"github.com/Alef-Daniel/challenge-go-opentelemetry/internal/serviceb/infra/gateway"
	"github.com/Alef-Daniel/challenge-go-opentelemetry/internal/serviceb/infra/tracing"
	"github.com/Alef-Daniel/challenge-go-opentelemetry/internal/serviceb/infra/web"
	"github.com/Alef-Daniel/challenge-go-opentelemetry/internal/serviceb/usecase"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

const (
	httpClientTimeout     = 10 * time.Second
	telemetryFlushTimeout = 5 * time.Second
)

var errMissingWeatherAPIKey = errors.New("WEATHER_API_KEY is required")

type config struct {
	port              string
	serviceName       string
	collectorEndpoint string
	viaCEPURL         string
	weatherAPIURL     string
	weatherAPIKey     string
}

func loadConfig() (config, error) {
	var cfg = config{
		port:              env.Get("PORT", "8081"),
		serviceName:       env.Get("OTEL_SERVICE_NAME", "service-b"),
		collectorEndpoint: env.Get("OTEL_COLLECTOR_ENDPOINT", "localhost:4317"),
		viaCEPURL:         env.Get("VIACEP_URL", "https://viacep.com.br"),
		weatherAPIURL:     env.Get("WEATHER_API_URL", "https://api.weatherapi.com"),
		weatherAPIKey:     env.Get("WEATHER_API_KEY", ""),
	}
	if cfg.weatherAPIKey == "" {
		return config{}, errMissingWeatherAPIKey
	}
	return cfg, nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	shutdown, err := telemetry.Setup(ctx, telemetry.Config{
		ServiceName:       cfg.serviceName,
		CollectorEndpoint: cfg.collectorEndpoint,
	})
	if err != nil {
		return err
	}
	defer func() {
		flushCtx, cancel := context.WithTimeout(context.Background(), telemetryFlushTimeout)
		defer cancel()
		_ = shutdown(flushCtx)
	}()

	client := &http.Client{
		Timeout:   httpClientTimeout,
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	tracer := otel.Tracer(cfg.serviceName)

	locations := tracing.NewLocationFinder(gateway.NewViaCEP(cfg.viaCEPURL, client), tracer)
	temperatures := tracing.NewTemperatureFinder(gateway.NewWeatherAPI(cfg.weatherAPIURL, cfg.weatherAPIKey, client), tracer)

	getWeather := usecase.NewGetWeather(locations, temperatures)
	router := web.NewRouter(web.NewWeatherHandler(getWeather))

	log.Printf("%s listening on :%s", cfg.serviceName, cfg.port)
	return httpserver.Run(ctx, ":"+cfg.port, router)
}

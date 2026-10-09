package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Alef-Daniel/challenge-go-opentelemetry/internal/platform/env"
	"github.com/Alef-Daniel/challenge-go-opentelemetry/internal/platform/httpserver"
	"github.com/Alef-Daniel/challenge-go-opentelemetry/internal/platform/telemetry"
	"github.com/Alef-Daniel/challenge-go-opentelemetry/internal/servicea/infra/gateway"
	"github.com/Alef-Daniel/challenge-go-opentelemetry/internal/servicea/infra/web"
	"github.com/Alef-Daniel/challenge-go-opentelemetry/internal/servicea/usecase"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	httpClientTimeout     = 15 * time.Second
	telemetryFlushTimeout = 5 * time.Second
)

type config struct {
	port              string
	serviceName       string
	collectorEndpoint string
	serviceBURL       string
}

func loadConfig() config {
	return config{
		port:              env.Get("PORT", "8080"),
		serviceName:       env.Get("OTEL_SERVICE_NAME", "service-a"),
		collectorEndpoint: env.Get("OTEL_COLLECTOR_ENDPOINT", "localhost:4317"),
		serviceBURL:       env.Get("SERVICE_B_URL", "http://localhost:8081"),
	}
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := loadConfig()

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

	getWeather := usecase.NewGetWeather(gateway.NewServiceB(cfg.serviceBURL, client))
	router := web.NewRouter(web.NewWeatherHandler(getWeather))

	log.Printf("%s listening on :%s", cfg.serviceName, cfg.port)
	return httpserver.Run(ctx, ":"+cfg.port, router)
}

package tracing

import (
	"context"

	"github.com/Alef-Daniel/challenge-go-opentelemetry/internal/serviceb/usecase"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const findTemperatureSpan = "busca-temperatura-weatherapi"

type TemperatureFinder struct {
	next   usecase.TemperatureFinder
	tracer trace.Tracer
}

func NewTemperatureFinder(next usecase.TemperatureFinder, tracer trace.Tracer) *TemperatureFinder {
	return &TemperatureFinder{next: next, tracer: tracer}
}

func (f *TemperatureFinder) FindCelsius(ctx context.Context, city string) (float64, error) {
	ctx, span := f.tracer.Start(ctx, findTemperatureSpan, trace.WithAttributes(
		attribute.String("city", city),
	))
	defer span.End()

	celsius, err := f.next.FindCelsius(ctx, city)
	if err != nil {
		recordError(span, err)
		return 0, err
	}

	span.SetAttributes(attribute.Float64("temperature.celsius", celsius))
	return celsius, nil
}

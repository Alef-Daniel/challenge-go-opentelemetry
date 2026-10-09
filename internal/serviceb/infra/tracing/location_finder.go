package tracing

import (
	"context"

	"github.com/Alef-Daniel/challenge-go-opentelemetry/internal/serviceb/domain"
	"github.com/Alef-Daniel/challenge-go-opentelemetry/internal/serviceb/usecase"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const findLocationSpan = "busca-cep-viacep"

type LocationFinder struct {
	next   usecase.LocationFinder
	tracer trace.Tracer
}

func NewLocationFinder(next usecase.LocationFinder, tracer trace.Tracer) *LocationFinder {
	return &LocationFinder{next: next, tracer: tracer}
}

func (f *LocationFinder) FindCity(ctx context.Context, zipcode domain.Zipcode) (string, error) {
	ctx, span := f.tracer.Start(ctx, findLocationSpan, trace.WithAttributes(
		attribute.String("zipcode", zipcode.String()),
	))
	defer span.End()

	city, err := f.next.FindCity(ctx, zipcode)
	if err != nil {
		recordError(span, err)
		return "", err
	}

	span.SetAttributes(attribute.String("city", city))
	return city, nil
}

package tracer

import (
	"context"
	"net/url"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// InitTracer initializes the OpenTelemetry tracer provider with OTLP HTTP exporter
func InitTracer(serviceName, collectorEndpoint string) (*sdktrace.TracerProvider, error) {
	opts := []otlptracehttp.Option{}
	endpoint := collectorEndpoint
	if endpoint != "" {
		if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
			u, err := url.Parse(endpoint)
			if err == nil {
				if u.Host != "" {
					opts = append(opts, otlptracehttp.WithEndpoint(u.Host))
				}
				if u.Path != "" {
					opts = append(opts, otlptracehttp.WithURLPath(strings.TrimSpace(u.Path)))
				}
				if u.Scheme == "http" {
					opts = append(opts, otlptracehttp.WithInsecure())
				}
			}
		} else {
			opts = append(opts, otlptracehttp.WithEndpoint(endpoint))
		}
	}
	exporter, err := otlptracehttp.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp, nil
}

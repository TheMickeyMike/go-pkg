package jaeger

import (
	"fmt"
	"log"

	"contrib.go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/trace"
)

// RegisterExporter add jaeger as trace exporter
func RegisterExporter(log *log.Logger, conf Config) (func(), error) {
	// Validate config first
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	// Initialize an exporter
	exporter, err := jaeger.NewExporter(jaeger.Options{
		CollectorEndpoint: conf.CollectorEndpoint,
		AgentEndpoint:     conf.AgentEndpoint,
		OnError: func(err error) {
			log.Printf("Error occured in Jaeger exporter: %s", err)
		},
		Process: jaeger.Process{
			ServiceName: conf.ServiceName,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("platform: failed to create jaeger exporter: %w", err)
	}

	// Register it
	trace.RegisterExporter(exporter)

	// No error
	return func() {
		exporter.Flush()
	}, nil
}

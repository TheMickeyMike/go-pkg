package jaeger

import (
	"fmt"

	"contrib.go.opencensus.io/exporter/jaeger"
	"github.com/TheMickeyMike/go-pkg/log"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
)

// RegisterExporter add jaeger as trace exporter
func RegisterExporter(conf Config) (func(), error) {
	// Validate config first
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	// Initialize an exporter
	exporter, err := jaeger.NewExporter(jaeger.Options{
		CollectorEndpoint: conf.CollectorEndpoint,
		AgentEndpoint:     conf.AgentEndpoint,
		OnError: func(err error) {
			log.Error("Error occured in Jaeger exporter", zap.Error(err))
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

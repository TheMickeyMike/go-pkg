package prometheus

import (
	"fmt"
	"log"
	"net/http"

	"contrib.go.opencensus.io/exporter/prometheus"
	"github.com/TheMickeyMike/go-pkg/log"
	"go.opencensus.io/stats/view"
	"go.uber.org/zap"
)

// RegisterExporter adds prometheus exporter
func RegisterExporter(conf Config, r *http.ServeMux) (func() error, error) {
	// Start prometheus
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	exporter, err := prometheus.NewExporter(prometheus.Options{
		Namespace: conf.Namespace,
		OnError: func(err error) {
			log.Error("Error occured in Prometheus exporter", zap.Error(err))
		},
	})

	if err != nil {
		return nil, fmt.Errorf("platform: unable to register prometheus exporter: %w", err)
	}

	// Add exporter
	view.RegisterExporter(exporter)

	// Add metrics handler
	r.Handle("/metrics", exporter)

	// No error
	return nil, nil
}

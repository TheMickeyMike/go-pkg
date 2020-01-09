package platform

import (
	"net/http"

	"os"

	"go.uber.org/zap"
	"github.com/TheMickeyMike/go-pkg/log"
	"github.com/TheMickeyMike/go-pkg/platform/jaeger"
	"github.com/TheMickeyMike/go-pkg/platform/prometheus"
)

type Server struct {
	RuntimeError    chan<- error
	Interrupt       <-chan os.Signal
	Instrumentation InstrumentationConfig
}

type InstrumentationConfig struct {
	Addr       string `json:"addres"`
	Prometheus struct {
		Enabled bool              `json:"enabled"`
		Config  prometheus.Config `json:"config"`
	} `json:"prometheus"`
	Jaeger struct {
		Enabled bool          `json:"enabled" `
		Config  jaeger.Config `json:"config"`
	} `json:"jaeger"`
}

func InstrumentationServer(cfg InstrumentationConfig) {
	instrumentationRouter := http.NewServeMux()

	if cfg.Prometheus.Enabled {
		if _, err := prometheus.RegisterExporter(cfg.Prometheus.Config, instrumentationRouter); err != nil {
			log.Fatal("Unable to register prometheus instrumentation", zap.Error(err))
		}
	}

	if cfg.Jaeger.Enabled {
		cancelFunc, err := jaeger.RegisterExporter(cfg.Jaeger.Config)
		if err != nil {
			log.Fatal("Unable to register jaeger instrumentation", zap.Error(err))
		}
		defer cancelFunc()
	}

	instrumentationServer := &http.Server{
		Addr:    cfg.Addr,
		Handler: instrumentationRouter,
	}
	go func() {
		_ = instrumentationServer.ListenAndServe()
	}()
}

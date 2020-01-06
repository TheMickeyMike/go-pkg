package platform

import (
	"net/http"

	"log"
	"os"

	"github.com/TheMickeyMike/go-pkg/platform/jaeger"
	"github.com/TheMickeyMike/go-pkg/platform/prometheus"
	"go.zenithar.org/pkg/log"
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

func InstrumentationServer(log *log.Logger, cfg InstrumentationConfig) {
	instrumentationRouter := http.NewServeMux()

	if cfg.Prometheus.Enabled {
		if _, err := prometheus.RegisterExporter(log, cfg.Prometheus.Config, instrumentationRouter); err != nil {
			log.Fatal("Unable to register prometheus instrumentation", err)
		}
	}

	if cfg.Jaeger.Enabled {
		cancelFunc, err := jaeger.RegisterExporter(log, cfg.Jaeger.Config)
		if err != nil {
			log.Fatal("Unable to register jaeger instrumentation", err)
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

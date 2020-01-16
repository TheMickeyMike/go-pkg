package platform

import (
	"context"
	"net/http"
	"time"

	"github.com/TheMickeyMike/go-pkg/log"
	"github.com/TheMickeyMike/go-pkg/platform/jaeger"
	"github.com/TheMickeyMike/go-pkg/platform/prometheus"
	"github.com/TheMickeyMike/go-pkg/platform/diagnostic"
	"go.uber.org/zap"
)

type instrumentation struct {
	server       *http.Server
	stopExporter func()
}

type InstrumentationConfig struct {
	Addr       string `json:"addres"`
	Diagnostic struct {
		Enabled bool              `json:"enabled"`
		Config  diagnostic.Config `json:"config"`
	} `json:"diagnostic"`
	Prometheus struct {
		Enabled bool              `json:"enabled"`
		Config  prometheus.Config `json:"config"`
	} `json:"prometheus"`
	Jaeger struct {
		Enabled bool          `json:"enabled" `
		Config  jaeger.Config `json:"config"`
	} `json:"jaeger"`
}

func NewInstrumentation(cfg InstrumentationConfig) *instrumentation {
	var (
		instrumentation instrumentation
		err             error
	)

	router := http.NewServeMux()

	if cfg.Diagnostic.Enabled {
		err := diagnostic.Register(cfg.Diagnostic.Config, router)
		if err != nil {
			log.Fatal("Unable to register diagnostic instrumentation", zap.Error(err))
		}
	}

	if cfg.Prometheus.Enabled {
		if _, err = prometheus.RegisterExporter(cfg.Prometheus.Config, router); err != nil {
			log.Fatal("Unable to register prometheus instrumentation", zap.Error(err))
		}
	}

	if cfg.Jaeger.Enabled {
		instrumentation.stopExporter, err = jaeger.RegisterExporter(cfg.Jaeger.Config)
		if err != nil {
			log.Fatal("Unable to register jaeger instrumentation", zap.Error(err))
		}
	}

	instrumentation.server = &http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	return &instrumentation

}

func (instrumentation *instrumentation) Run(ctx context.Context) <-chan error {
	err := make(chan error, 1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				defer instrumentation.stopExporter()
				shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				err := instrumentation.server.Shutdown(shutdownCtx)
				if err != nil {
					log.Fatal("Error raised while shutting down the server", zap.Error(err))
				}
				return
			case err <- instrumentation.server.ListenAndServe():
				return
			}
		}
	}()
	return err
}

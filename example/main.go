package main

import (
	"context"
	"fmt"

	"github.com/TheMickeyMike/go-pkg/config"
	"github.com/TheMickeyMike/go-pkg/platform"
)

const conf = `
{
	"addres": ":8080",
	"diagnostic": {
		"enabled": true,
		"config": {
			"gops": {
				"enabled": true
			},
			"pprof": {
				"enabled": true
			},
			"zpages": {
				"enabled": true
			}
		}
	},
    "prometheus": {
        "enabled": true,
        "config": {
            "namespace": "example_app"
        }
    },
    "jaeger": {
        "enabled": true,
        "config": {
            "collector_endpoint": "http://localhost:14268/api/traces",
            "agent_endpoint": "localhost:6831",
            "service_name": "maciej-test"
        }
    }
}
`

func main() {
	var instrumentationConfig platform.InstrumentationConfig
	if err := config.LoadConfig(&instrumentationConfig, conf); err != nil {
		fmt.Println(err)
	}
	fmt.Println(instrumentationConfig)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	instrumenattionServer := platform.NewInstrumentation(instrumentationConfig)
	fmt.Println("Start")
	instrumenattionError := instrumenattionServer.Run(ctx)
	err := <-instrumenattionError
	if err != nil {
		fmt.Println(err)
	}
}

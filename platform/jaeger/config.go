package jaeger

import "fmt"

// Config holds information necessary for sending trace to jaeger.
type Config struct {
	// CollectorEndpoint is the Jaeger HTTP Thrift endpoint.
	// For example, http://localhost:14268/api/traces?format=jaeger.thrift.
	CollectorEndpoint string `json:"collector_endpoint"`

	// AgentEndpoint instructs exporter to send spans to Jaeger agent at this address.
	// For example, localhost:6831.
	AgentEndpoint string `json:"agent_endpoint"`

	// ServiceName is the name of the process.
	ServiceName string `json:"service_name"`
}

// Validate checks that the configuration is valid.
func (c Config) Validate() error {
	if c.CollectorEndpoint == "" && c.AgentEndpoint == "" {
		return fmt.Errorf("jaeger: either collector endpoint or agent endpoint must be configured")
	}
	if c.ServiceName == "" {
		return fmt.Errorf("jaeger: service name must not be blank")
	}

	return nil
}

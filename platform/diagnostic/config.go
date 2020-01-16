package diagnostic

// Config holds information for diagnostic handlers.
type Config struct {
	GOPS struct {
		Enabled   bool   `json:"enabled"`
		RemoteURL string `json:"remoteDebugURL"`
	}`json:"gops"`
	PProf struct {
		Enabled bool `json:"enabled"`
	}`json:"pprof"`
	ZPages struct {
		Enabled bool `json:"enabled"`
	}`json:"zpages"`
}

// Validate checks that the configuration is valid.
func (c Config) Validate() error {
	return nil
}

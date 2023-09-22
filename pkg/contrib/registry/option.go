package registry

type GetRegistryOption func() (*Option, error)

type Option struct {
	Enabled   bool     `json:"enabled" yaml:"enabled"`
	Endpoints []string `json:"endpoints" yaml:"endpoints"`
}

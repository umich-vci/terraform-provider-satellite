package satellite

import (
	"github.com/umich-vci/gosatellite"
)

// Config holds the provider configuration
type Config struct {
	Username      string
	Password      string
	SatelliteHost string
	SSLVerify     bool
}

// Client returns a new client for accessing Red Hat Satellite
func (c *Config) Client() (*gosatellite.Client, error) {
	config := new(gosatellite.Config)
	config.Username = c.Username
	config.Password = c.Password
	config.SatelliteHost = c.SatelliteHost
	config.SSLVerify = c.SSLVerify

	return gosatellite.NewClient(config)
}

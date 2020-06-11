package teamcity

import (
	"net/http"

	api "github.com/cvbarros/go-teamcity/teamcity"
)

// Config Used to configure an api client for TeamCity
type Config struct {
	Address  string
	Token    string
	Username string
	Password string
}

// Client Returns a new TeamCity api client configured with this instance parameters
func (c *Config) Client() (*api.Client, error) {
	// `http.DefaultClient` doesn't configure a proxy by default - this does
	httpClient := &http.Client{
		Transport: http.DefaultTransport,
	}

	if c.Token != "" {
		return api.NewClientWithAddress(api.TokenAuth(c.Token), c.Address, httpClient)
	}

	return api.NewClientWithAddress(api.BasicAuth(c.Username, c.Password), c.Address, httpClient)
}

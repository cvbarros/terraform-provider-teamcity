package teamcity

import (
	api "github.com/cvbarros/go-teamcity-sdk/client"
	runtime "github.com/go-openapi/runtime/client"
)

// Config Used to configure an api client for TeamCity
type Config struct {
	Address  string
	Username string
	Password string
}

// Client Returns a new TeamCity api client configured with this instance parameters
func (c *Config) Client() *api.TeamCityREST {
	clientTransport := runtime.New(c.Address, "/", []string{"http"})
	//clientTransport.SetDebug(true)
	clientTransport.DefaultAuthentication = runtime.BasicAuth(c.Username, c.Password)

	client := api.New(clientTransport, nil)
	return client
}

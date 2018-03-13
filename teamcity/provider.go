package teamcity

import (
	api "github.com/cvbarros/go-teamcity-sdk/client"
	runtime "github.com/go-openapi/runtime/client"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"teamcity_project": resourceProject(),
		},
		Schema: map[string]*schema.Schema{
			"server_url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TEAMCITY_URL", nil),
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TEAMCITY_USER", nil),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TEAMCITY_PASSWORD", nil),
			},
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	clientTransport := runtime.New(d.Get("server_url").(string), "/", []string{"http"})
	clientTransport.DefaultAuthentication = runtime.BasicAuth(
		d.Get("username").(string),
		d.Get("password").(string))

	client := api.New(clientTransport, nil)

	return client, nil
}

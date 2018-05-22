package teamcity

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

//Provider is the plugin entry point
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"teamcity_project":                         resourceProject(),
			"teamcity_vcs_root_git":                    resourceVcsRootGit(),
			"teamcity_build_config":                    resourceBuildConfiguration(),
			"teamcity_snapshot_dependency":             resourceSnapshotDependency(),
			"teamcity_trigger":                         resourceTrigger(),
			"teamcity_agent_requirement":               resourceAgentRequirement(),
			"teamcity_feature_commit_status_publisher": resourceFeatureCommitStatusPublisher(),
		},
		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
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
	config := Config{
		Address:  d.Get("address").(string),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
	}
	client := config.Client()
	return client, nil
}

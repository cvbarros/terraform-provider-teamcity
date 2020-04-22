package teamcity

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

//Provider is the plugin entry point
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"teamcity_artifact_dependency":                resourceArtifactDependency(),
			"teamcity_agent_requirement":                  resourceAgentRequirement(),
			"teamcity_build_config":                       resourceBuildConfig(),
			"teamcity_build_trigger_build_finish":         resourceBuildTriggerBuildFinish(),
			"teamcity_build_trigger_schedule":             resourceBuildTriggerSchedule(),
			"teamcity_build_trigger_vcs":                  resourceBuildTriggerVcs(),
			"teamcity_feature_commit_status_publisher":    resourceFeatureCommitStatusPublisher(),
			"teamcity_group":                              resourceGroup(),
			"teamcity_project":                            resourceProject(),
			"teamcity_project_feature_versioned_settings": resourceProjectFeatureVersionedSettings(),
			"teamcity_snapshot_dependency":                resourceSnapshotDependency(),
			"teamcity_vcs_root_git":                       resourceVcsRootGit(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"teamcity_project": dataSourceProject(),
		},
		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TEAMCITY_ADDR", nil),
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
	return config.Client()
}

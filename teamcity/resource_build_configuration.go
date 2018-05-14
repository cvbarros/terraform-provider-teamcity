package teamcity

import (
	api "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBuildConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceBuildConfigurationCreate,
		Read:   resourceBuildConfigurationRead,
		Update: resourceBuildConfigurationUpdate,
		Delete: resourceBuildConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceBuildConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	bt := &api.BuildType{}
	var projectId string

	if v, ok := d.GetOk("project_id"); ok {
		projectId = v.(string)
		bt.ProjectID = projectId
	}

	if v, ok := d.GetOk("name"); ok {
		bt.Name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		bt.Description = v.(string)
	}

	created, err := client.BuildTypes.Create(projectId, bt)
	if err != nil {
		return err
	}

	d.MarkNewResource()
	d.SetId(created.ID)
	return nil
}

func resourceBuildConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBuildConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBuildConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	return client.BuildTypes.Delete(d.Id())
	return nil
}

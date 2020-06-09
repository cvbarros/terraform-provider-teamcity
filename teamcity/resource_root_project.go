package teamcity

import (
	api "github.com/cvbarros/go-teamcity/teamcity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRootProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceRootProjectCreate,
		Delete: resourceRootProjectDelete,
		Read:   resourceProjectRead,
		Importer: &schema.ResourceImporter{
			State: resourceProjectImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"env_params": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"config_params": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"sys_params": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceRootProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	d.SetId("_Root")
	return resourceProjectUpdate(d, client)
}

func resourceRootProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	dt, err := getProject(client, d.Id())

	// Empty all the parameters, so that it will be destroyed
	dt.Parameters.Items = nil
	_, err = client.Projects.Update(dt)

	if err != nil {
		return err
	}
	return nil
}

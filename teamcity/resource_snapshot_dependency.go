package teamcity

import (
	"fmt"

	api "github.com/cvbarros/go-teamcity-sdk/teamcity"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSnapshotDependency() *schema.Resource {
	return &schema.Resource{
		Create: resourceSnapshotDependencyCreate,
		Read:   resourceSnapshotDependencyRead,
		Delete: resourceSnapshotDependencyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"build_config_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_build_config_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSnapshotDependencyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	var buildConfigID string

	if v, ok := d.GetOk("build_config_id"); ok {
		buildConfigID = v.(string)
	}
	// validates the Build Configuration exists
	if _, err := client.BuildTypes.GetByID(buildConfigID); err != nil {
		return fmt.Errorf("invalid build_config_id '%s' - Build configuration does not exist", buildConfigID)
	}

	depService := client.DependencyService(buildConfigID)
	dep := api.NewSnapshotDependency(d.Get("source_build_config_id").(string))

	out, err := depService.AddSnapshotDependency(dep)

	if err != nil {
		return err
	}

	d.SetId(out.ID)

	return resourceSnapshotDependencyRead(d, meta)
}

func resourceSnapshotDependencyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client).DependencyService(d.Get("build_config_id").(string))

	dt, err := getSnapshotDependency(client, d.Id())
	if err != nil {
		return err
	}

	if err := d.Set("build_config_id", dt.BuildTypeID); err != nil {
		return err
	}

	return d.Set("source_build_config_id", dt.SourceBuildType.ID)
}

func resourceSnapshotDependencyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	dep := client.DependencyService(d.Get("build_config_id").(string))

	return dep.DeleteSnapshot(d.Id())
}

func getSnapshotDependency(c *api.DependencyService, id string) (*api.SnapshotDependency, error) {

	dt, err := c.GetSnapshotByID(id)
	if err != nil {
		return nil, err
	}

	return dt, nil
}

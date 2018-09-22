package teamcity

import (
	"fmt"

	api "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBuildTriggerVcs() *schema.Resource {
	return &schema.Resource{
		Create: resourceBuildTriggerVcsCreate,
		Read:   resourceBuildTriggerVcsRead,
		Update: resourceBuildTriggerVcsUpdate,
		Delete: resourceBuildTriggerVcsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"build_config_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rules": {
				Type:     schema.TypeString,
				Required: true,
			},
			"branch_filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceBuildTriggerVcsCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	var buildConfigID string
	var err error

	if v, ok := d.GetOk("build_config_id"); ok {
		buildConfigID = v.(string)
	}
	// validates the Build Configuration exists
	if _, err := client.BuildTypes.GetByID(buildConfigID); err != nil {
		return fmt.Errorf("invalid build_config_id '%s' - Build configuration does not exist", buildConfigID)
	}

	ts := client.TriggerService(buildConfigID)
	var dt *api.TriggerVcs
	if v, ok := d.GetOk("rules"); ok {
		dt, err = api.NewTriggerVcs(v.(string), "")
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Error getting required property 'rules' for vcs trigger")
	}

	if v, ok := d.GetOk("branch_filter"); ok {
		dt.BranchFilter = v.(string)
	}

	out, err := ts.AddTrigger(dt)

	if err != nil {
		return err
	}

	d.SetId(out.ID())

	return resourceBuildTriggerVcsRead(d, meta)
}

func resourceBuildTriggerVcsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client).TriggerService(d.Get("build_config_id").(string))

	ret, err := getTrigger(client, d.Id())
	if err != nil {
		return err
	}
	dt, ok := ret.(*api.TriggerVcs)
	if !ok {
		return fmt.Errorf("invalid trigger type when reading build_trigger_vcs resource")
	}

	if err := d.Set("build_config_id", dt.BuildTypeID()); err != nil {
		return err
	}

	if err := d.Set("rules", dt.Rules); err != nil {
		return err
	}

	if dt.BranchFilter != "" {
		if err := d.Set("branch_filter", dt.BranchFilter); err != nil {
			return err
		}
	}

	return nil
}

func resourceBuildTriggerVcsUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBuildTriggerVcsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	ts := client.TriggerService(d.Get("build_config_id").(string))

	return ts.Delete(d.Id())
}

func getTrigger(c *api.TriggerService, id string) (api.Trigger, error) {

	dt, err := c.GetByID(id)
	if err != nil {
		return nil, err
	}

	return dt, nil
}

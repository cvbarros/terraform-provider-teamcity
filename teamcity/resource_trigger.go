package teamcity

import (
	"fmt"

	api "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTrigger() *schema.Resource {
	return &schema.Resource{
		Create: resourceTriggerCreate,
		Read:   resourceTriggerRead,
		Update: resourceTriggerUpdate,
		Delete: resourceTriggerDelete,
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

func resourceTriggerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	var buildConfigID string

	if v, ok := d.GetOk("build_config_id"); ok {
		buildConfigID = v.(string)
	}
	// validates the Build Configuration exists
	if _, err := client.BuildTypes.GetByID(buildConfigID); err != nil {
		return fmt.Errorf("invalid build_config_id '%s' - Build configuration does not exist", buildConfigID)
	}

	ts := client.TriggerService(buildConfigID)
	var dt *api.Trigger
	if v, ok := d.GetOk("rules"); ok {
		dt = api.NewVcsTrigger(v.(string), "")
	} else {
		return fmt.Errorf("Error getting required property 'rules' for trigger")
	}

	if v, ok := d.GetOk("branch_filter"); ok {
		dt.SetBranchFilter(v.(string))
	}

	out, err := ts.AddTrigger(dt)

	if err != nil {
		return err
	}

	d.SetId(out.ID)

	return resourceTriggerRead(d, meta)
}

func resourceTriggerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client).TriggerService(d.Get("build_config_id").(string))

	dt, err := getTrigger(client, d.Id())
	if err != nil {
		return err
	}

	if err := d.Set("build_config_id", dt.BuildTypeID); err != nil {
		return err
	}

	if err := d.Set("rules", dt.Rules()); err != nil {
		return err
	}

	if v, ok := dt.BranchFilterOk(); ok {
		if err := d.Set("branch_filter", v); err != nil {
			return err
		}
	}

	return nil
}

func resourceTriggerUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceTriggerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	ts := client.TriggerService(d.Get("build_config_id").(string))

	return ts.Delete(d.Id())
}

func getTrigger(c *api.TriggerService, id string) (*api.Trigger, error) {

	dt, err := c.GetByID(id)
	if err != nil {
		return nil, err
	}

	return dt, nil
}

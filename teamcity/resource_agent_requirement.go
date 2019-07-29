package teamcity

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/validation"

	api "github.com/cvbarros/go-teamcity-sdk/teamcity"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAgentRequirement() *schema.Resource {
	return &schema.Resource{
		Create: resourceAgentRequirementCreate,
		Read:   resourceAgentRequirementRead,
		Delete: resourceAgentRequirementDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"build_config_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"condition": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(api.ConditionStrings, false),
				Required:     true,
				ForceNew:     true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAgentRequirementCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	var buildConfigID string

	if v, ok := d.GetOk("build_config_id"); ok {
		buildConfigID = v.(string)
	}
	// validates the Build Configuration exists
	if _, err := client.BuildTypes.GetByID(buildConfigID); err != nil {
		return fmt.Errorf("invalid build_config_id '%s' - Build configuration does not exist", buildConfigID)
	}

	srv := client.AgentRequirementService(buildConfigID)

	var condition, name, value string

	if v, ok := d.GetOk("condition"); ok {
		condition = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	if v, ok := d.GetOk("value"); ok {
		value = v.(string)
	}

	dt, err := api.NewAgentRequirement(condition, name, value)
	if err != nil {
		return err
	}
	out, err := srv.Create(dt)

	if err != nil {
		return err
	}

	d.SetId(out.ID)

	return resourceAgentRequirementRead(d, meta)
}

func resourceAgentRequirementRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client).AgentRequirementService(d.Get("build_config_id").(string))

	dt, err := getAgentRequirement(client, d.Id())
	if err != nil {
		return err
	}

	if err := d.Set("build_config_id", dt.BuildTypeID); err != nil {
		return err
	}

	if err := d.Set("condition", dt.Condition); err != nil {
		return err
	}

	if err := d.Set("name", dt.Name()); err != nil {
		return err
	}

	if v := dt.Value(); v != "" {
		if err := d.Set("value", v); err != nil {
			return err
		}
	}

	return nil
}

func resourceAgentRequirementDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	svr := client.AgentRequirementService(d.Get("build_config_id").(string))

	return svr.Delete(d.Id())
}

func getAgentRequirement(c *api.AgentRequirementService, id string) (*api.AgentRequirement, error) {
	dt, err := c.GetByID(id)
	if err != nil {
		return nil, err
	}

	return dt, nil
}

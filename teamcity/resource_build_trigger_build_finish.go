package teamcity

import (
	"fmt"
	"log"

	api "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBuildTriggerBuildFinish() *schema.Resource {
	return &schema.Resource{
		Create: resourceBuildTriggerBuildFinishCreate,
		Read:   resourceBuildTriggerBuildFinishRead,
		Update: resourceBuildTriggerBuildFinishUpdate,
		Delete: resourceBuildTriggerBuildFinishDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"build_config_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"after_successful_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"branch_filter": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

func resourceBuildTriggerBuildFinishCreate(d *schema.ResourceData, meta interface{}) error {
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
	opt := api.NewTriggerBuildFinishOptions(false, nil)
	dt, err := api.NewTriggerBuildFinish(buildConfigID, opt)
	if err != nil {
		return err
	}

	if v, ok := d.GetOk("after_successful_only"); ok {
		dt.Options.AfterSuccessfulBuildOnly = v.(bool)
	}

	log.Printf("[INFO] BranchFilter: %s, State: %s", dt.Options.BranchFilter, d.Get("branch_filter"))
	if v, ok := d.GetOk("branch_filter"); ok {
		dt.Options.BranchFilter = expandStringSlice(v.([]interface{}))
		log.Printf("[INFO] BranchFilter: %s, State: %s", dt.Options.BranchFilter, v)
	}

	out, err := ts.AddTrigger(dt)

	if err != nil {
		return err
	}

	d.SetId(out.ID())

	return resourceBuildTriggerBuildFinishRead(d, meta)
}

func resourceBuildTriggerBuildFinishRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client).TriggerService(d.Get("build_config_id").(string))

	ret, err := getTrigger(client, d.Id())
	if err != nil {
		return err
	}
	dt, ok := ret.(*api.TriggerBuildFinish)
	if !ok {
		return fmt.Errorf("invalid trigger type when reading build_trigger_build_finish resource")
	}

	if err := d.Set("build_config_id", dt.BuildTypeID()); err != nil {
		return err
	}
	log.Printf("[INFO] READ: BranchFilter: %s, State: %s", dt.Options.BranchFilter, d.Get("branch_filter"))
	if err := d.Set("branch_filter", flattenStringSlice(dt.Options.BranchFilter)); err != nil {
		return err
	}

	if dt.Options.AfterSuccessfulBuildOnly {
		if err := d.Set("after_sucessful_only", dt.Options.AfterSuccessfulBuildOnly); err != nil {
			return err
		}
	}

	return nil
}

func resourceBuildTriggerBuildFinishUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBuildTriggerBuildFinishDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	ts := client.TriggerService(d.Get("build_config_id").(string))

	return ts.Delete(d.Id())
}

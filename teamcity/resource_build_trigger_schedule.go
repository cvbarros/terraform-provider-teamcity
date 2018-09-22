package teamcity

import (
	"fmt"

	api "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceBuildTriggerSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceBuildTriggerScheduleCreate,
		Read:   resourceBuildTriggerScheduleRead,
		Update: resourceBuildTriggerScheduleUpdate,
		Delete: resourceBuildTriggerScheduleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"build_config_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"schedule": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"daily", "weekly"}, false),
			},
			"hour": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"minute": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "SERVER",
			},
			"weekday": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Sunday",
					"Monday",
					"Tuesday",
					"Wednesday",
					"Thursday",
					"Friday",
					"Saturday"}, false),
			},
			"rules": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

func resourceBuildTriggerScheduleCreate(d *schema.ResourceData, meta interface{}) error {
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

	hour := d.Get("hour").(int)
	minute := d.Get("minute").(int)
	timezone := d.Get("timezone").(string)
	rules := expandStringSlice(d.Get("rules").([]interface{}))
	schedule := d.Get("schedule").(string)
	weekday, _ := parseWeekday(d.Get("weekday").(string))
	dt, err := api.NewTriggerSchedule(schedule, buildConfigID, weekday, uint(hour), uint(minute), timezone, rules, api.NewTriggerScheduleOptions())

	if err != nil {
		return err
	}

	out, err := ts.AddTrigger(dt)

	if err != nil {
		return err
	}

	d.SetId(out.ID())

	return resourceBuildTriggerScheduleRead(d, meta)
}

func resourceBuildTriggerScheduleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client).TriggerService(d.Get("build_config_id").(string))

	ret, err := getTrigger(client, d.Id())
	if err != nil {
		return err
	}
	dt, ok := ret.(*api.TriggerSchedule)
	if !ok {
		return fmt.Errorf("invalid trigger type when reading build_trigger_schedule resource")
	}

	if err := d.Set("build_config_id", dt.BuildTypeID()); err != nil {
		return err
	}
	if err := d.Set("schedule", dt.SchedulingPolicy); err != nil {
		return err
	}
	if err := d.Set("hour", dt.Hour); err != nil {
		return err
	}
	if err := d.Set("minute", dt.Minute); err != nil {
		return err
	}
	if err := d.Set("timezone", dt.Timezone); err != nil {
		return err
	}
	if err := d.Set("rules", flattenStringSlice(dt.Rules)); err != nil {
		return err
	}
	if dt.SchedulingPolicy == api.TriggerSchedulingWeekly {
		if err := d.Set("weekday", dt.Weekday.String()); err != nil {
			return err
		}
	}

	return nil
}

func resourceBuildTriggerScheduleUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBuildTriggerScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	ts := client.TriggerService(d.Get("build_config_id").(string))

	return ts.Delete(d.Id())
}

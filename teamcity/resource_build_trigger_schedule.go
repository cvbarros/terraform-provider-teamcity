package teamcity

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/robfig/cron"
	api "github.com/yext/go-teamcity/teamcity"
)

func resourceBuildTriggerSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceBuildTriggerScheduleCreate,
		Read:   resourceBuildTriggerScheduleRead,
		Delete: resourceBuildTriggerScheduleDelete,
		Importer: &schema.ResourceImporter{
			State: subresourceImporter(resourceBuildTriggerScheduleRead),
		},

		Schema: map[string]*schema.Schema{
			"build_config_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"schedule": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"daily", "weekly", "cron"}, false),
			},
			"hour": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Optional: true,
			},
			"minute": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Optional: true,
				Default:  0,
			},
			"timezone": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  "SERVER",
			},
			"weekday": {
				Type:     schema.TypeString,
				ForceNew: true,
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
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"cron_schedule": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"seconds": {
							Type:     schema.TypeString,
							Required: true,
						},
						"minutes": {
							Type:     schema.TypeString,
							Required: true,
						},
						"hours": {
							Type:     schema.TypeString,
							Required: true,
						},
						"day_of_month": {
							Type:     schema.TypeString,
							Required: true,
						},
						"month": {
							Type:     schema.TypeString,
							Required: true,
						},
						"day_of_week": {
							Type:     schema.TypeString,
							Required: true,
						},
						"year": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"enforce_clean_checkout": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
			},
			"enforce_clean_checkout_dependencies": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
			},
			"queue_optimization": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
				Default:  true,
			},
			"on_all_compatible_agents": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
			},
			"with_pending_changes_only": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
				Default:  true,
			},
			"promote_watched_build": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
				Default:  true,
			},
			"only_if_watched_changes": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
			},
			"watched_build_config_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"revision": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  "lastFinished",
				ValidateFunc: validation.StringInSlice([]string{
					string(api.LatestFinishedBuild),
					string(api.LatestSuccessfulBuild),
					string(api.LastBuildFinishedWithTag),
				}, false),
			},
			"watched_branch": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  "<default>",
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

	opt, err := expandTriggerScheduleOptions(d)
	if err != nil {
		return err
	}

	var cronSchedule *api.CronSchedule
	v, ok := d.GetOk("cron_schedule")
	if ok {
		cronSchedule, err = expandCronSchedule(v.([]interface{}))
		if err != nil {
			return err
		}
	}

	dt, err := api.NewTriggerSchedule(schedule, buildConfigID, weekday, uint(hour), uint(minute), timezone, rules, cronSchedule, opt)
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
	if dt.SchedulingPolicy == api.TriggerSchedulingCron {
		if dt.CronExpression == nil {
			return fmt.Errorf("cron expression was not specified")
		}
		err := d.Set("cron_schedule", []map[string]interface{}{flattenCronSchedule(dt.CronExpression)})
		if err != nil {
			return err
		}
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
	flatOpt := flattenTriggerScheduleOptions(dt.Options)
	for k, v := range flatOpt {
		if err := d.Set(k, v); err != nil {
			return err
		}
	}
	return nil
}

func resourceBuildTriggerScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	ts := client.TriggerService(d.Get("build_config_id").(string))

	return ts.Delete(d.Id())
}

func expandTriggerScheduleOptions(d *schema.ResourceData) (*api.TriggerScheduleOptions, error) {
	opt := api.NewTriggerScheduleOptions()

	if v, ok := d.GetOkExists("queue_optimization"); ok {
		opt.QueueOptimization = v.(bool)
	}
	if v, ok := d.GetOkExists("on_all_compatible_agents"); ok {
		opt.BuildOnAllCompatibleAgents = v.(bool)
	}
	if v, ok := d.GetOkExists("with_pending_changes_only"); ok {
		opt.BuildWithPendingChangesOnly = v.(bool)
	}
	if v, ok := d.GetOkExists("promote_watched_build"); ok {
		opt.PromoteWatchedBuild = v.(bool)
	}
	if v, ok := d.GetOkExists("enforce_clean_checkout"); ok {
		opt.EnforceCleanCheckout = v.(bool)
	}
	if v, ok := d.GetOkExists("enforce_clean_checkout_dependencies"); ok {
		opt.EnforceCleanCheckoutForDependencies = v.(bool)
	}
	if v, ok := d.GetOkExists("only_if_watched_changes"); ok {
		opt.TriggerIfWatchedBuildChanges = v.(bool)
	}
	if v, ok := d.GetOkExists("watched_build_config_id"); ok {
		opt.RevisionRuleSourceBuildID = v.(string)
	}
	if v, ok := d.GetOkExists("revision"); ok {
		opt.RevisionRule = api.ArtifactDependencyRevision(v.(string))
	}
	if v, ok := d.GetOkExists("watched_branch"); ok {
		opt.RevisionRuleBuildBranch = v.(string)
	}

	return opt, nil
}

func expandCronSchedule(v []interface{}) (*api.CronSchedule, error) {
	raw := v[0].(map[string]interface{})
	var cronSchedule api.CronSchedule

	if v, ok := raw["seconds"]; ok {
		cronSchedule.Seconds = v.(string)
	}
	if v, ok := raw["minutes"]; ok {
		cronSchedule.Minutes = v.(string)
	}
	if v, ok := raw["hours"]; ok {
		cronSchedule.Hours = v.(string)
	}
	if v, ok := raw["day_of_month"]; ok {
		cronSchedule.DayOfMonth = v.(string)
	}
	if v, ok := raw["month"]; ok {
		cronSchedule.Month = v.(string)
	}
	if v, ok := raw["day_of_week"]; ok {
		cronSchedule.DayOfWeek = v.(string)
	}
	if v, ok := raw["year"]; ok {
		cronSchedule.Year = v.(string)
	}
	// Validate the cron expression
	cronParser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	_, err := cronParser.Parse(
		fmt.Sprintf(
			"%v %v %v %v %v %v",
			cronSchedule.Seconds,
			cronSchedule.Minutes,
			cronSchedule.Hours,
			cronSchedule.DayOfMonth,
			cronSchedule.Month,
			cronSchedule.DayOfWeek,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("parsing cron expression: %v", err)
	}

	return &cronSchedule, nil
}

func flattenTriggerScheduleOptions(dt *api.TriggerScheduleOptions) map[string]interface{} {
	out := make(map[string]interface{})
	out["queue_optimization"] = dt.QueueOptimization
	out["promote_watched_build"] = dt.PromoteWatchedBuild
	out["with_pending_changes_only"] = dt.BuildWithPendingChangesOnly
	out["revision"] = dt.RevisionRule
	out["watched_branch"] = dt.RevisionRuleBuildBranch

	if dt.BuildOnAllCompatibleAgents {
		out["on_all_compatible_agents"] = dt.BuildOnAllCompatibleAgents
	}
	if dt.EnforceCleanCheckout {
		out["enforce_clean_checkout"] = dt.EnforceCleanCheckout
	}
	if dt.EnforceCleanCheckoutForDependencies {
		out["enforce_clean_checkout_dependencies"] = dt.EnforceCleanCheckoutForDependencies
	}
	if dt.TriggerIfWatchedBuildChanges {
		out["only_if_watched_changes"] = dt.TriggerIfWatchedBuildChanges
	}
	if dt.RevisionRuleSourceBuildID != "" {
		out["watched_build_config_id"] = dt.RevisionRuleSourceBuildID
	}

	return out
}

func flattenCronSchedule(dt *api.CronSchedule) map[string]interface{} {
	m := make(map[string]interface{})

	m["seconds"] = dt.Seconds
	m["minutes"] = dt.Minutes
	m["hours"] = dt.Hours
	m["day_of_month"] = dt.DayOfMonth
	m["day_of_week"] = dt.DayOfWeek
	m["month"] = dt.Month
	m["year"] = dt.Year

	return m
}

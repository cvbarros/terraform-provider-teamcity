package teamcity

import (
	"fmt"

	api "github.com/cvbarros/go-teamcity-sdk/teamcity"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceBuildTriggerSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceBuildTriggerScheduleCreate,
		Read:   resourceBuildTriggerScheduleRead,
		Delete: resourceBuildTriggerScheduleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				ValidateFunc: validation.StringInSlice([]string{"daily", "weekly"}, false),
			},
			"hour": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Required: true,
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

	dt, err := api.NewTriggerSchedule(schedule, buildConfigID, weekday, uint(hour), uint(minute), timezone, rules, opt)

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

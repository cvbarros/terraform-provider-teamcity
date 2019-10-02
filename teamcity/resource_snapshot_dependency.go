package teamcity

import (
	"fmt"
	"strconv"
	"strings"

	api "github.com/cvbarros/go-teamcity-sdk/teamcity"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

var (
	runBuildOptions  = []string{"RUN_ADD_PROBLEM", "RUN", "MAKE_FAILED_TO_START", "CANCEL"}
	snapshotIDFormat = "%s %s" // "<build_config_id> <snapshot dependency ID>"
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
			"run_build_if_dependency_failed": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      api.DefaultSnapshotDependencyOptions.OnFailedDependency,
				ValidateFunc: validation.StringInSlice(runBuildOptions, false),
			},
			"run_build_if_dependency_failed_to_start": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      api.DefaultSnapshotDependencyOptions.OnFailedToStartOrCanceledDependency,
				ValidateFunc: validation.StringInSlice(runBuildOptions, false),
			},
			"run_build_on_the_same_agent": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  api.DefaultSnapshotDependencyOptions.RunSameAgent,
				ForceNew: true,
			},
			"take_started_build_with_same_revisions": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  api.DefaultSnapshotDependencyOptions.DoNotRunNewBuildIfThereIsASuitable,
				ForceNew: true,
			},
			"take_successful_builds_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  api.DefaultSnapshotDependencyOptions.TakeSuccessfulBuildsOnly,
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

	opts := getSnapshotDependencyOptions(d)
	dep := api.NewSnapshotDependencyWithOptions(d.Get("source_build_config_id").(string), opts)

	depService := client.DependencyService(buildConfigID)
	out, err := depService.AddSnapshotDependency(dep)

	if err != nil {
		return err
	}

	// Since snapshote dependencies are sub-resources of build_configs, the
	// buildConfigID is included as part of the snapshot dependency's ID.
	d.SetId(fmt.Sprintf(snapshotIDFormat, buildConfigID, out.ID))

	return resourceSnapshotDependencyRead(d, meta)
}

func getSnapshotDependencyOptions(d *schema.ResourceData) *api.SnapshotDependencyOptions {
	// Initialize as a copy of defaults
	opts := *api.DefaultSnapshotDependencyOptions

	if v, ok := d.GetOk("run_build_if_dependency_failed"); ok {
		opts.OnFailedDependency = v.(string)
	}
	if v, ok := d.GetOk("run_build_if_dependency_failed_to_start"); ok {
		opts.OnFailedToStartOrCanceledDependency = v.(string)
	}
	if v, ok := d.GetOkExists("run_build_on_the_same_agent"); ok {
		opts.RunSameAgent = v.(bool)
	}
	if v, ok := d.GetOkExists("take_started_build_with_same_revisions"); ok {
		opts.DoNotRunNewBuildIfThereIsASuitable = v.(bool)
	}
	if v, ok := d.GetOkExists("take_successful_builds_only"); ok {
		opts.TakeSuccessfulBuildsOnly = v.(bool)
	}

	return &opts
}

func setSnapshotDependencyProperties(d *schema.ResourceData, p *api.Properties) error {
	if err := setSnapshotDependencyStringProperty(d, p, "run-build-if-dependency-failed"); err != nil {
		return err
	}
	if err := setSnapshotDependencyStringProperty(d, p, "run-build-if-dependency-failed-to-start"); err != nil {
		return err
	}
	if err := setSnapshotDependencyBoolProperty(d, p, "run-build-on-the-same-agent"); err != nil {
		return err
	}
	if err := setSnapshotDependencyBoolProperty(d, p, "take-started-build-with-same-revisions"); err != nil {
		return err
	}
	if err := setSnapshotDependencyBoolProperty(d, p, "take-successful-builds-only"); err != nil {
		return err
	}

	return nil
}

func setSnapshotDependencyStringProperty(d *schema.ResourceData, p *api.Properties, propName string) error {
	if v, ok := p.GetOk(propName); ok {
		schemaName := strings.ReplaceAll(propName, "-", "_")
		if err := d.Set(schemaName, v); err != nil {
			return err
		}
	}
	return nil
}

func setSnapshotDependencyBoolProperty(d *schema.ResourceData, p *api.Properties, propName string) error {
	if v, ok := p.GetOk(propName); ok {
		schemaName := strings.ReplaceAll(propName, "-", "_")
		b, err := strconv.ParseBool(v)
		if err != nil {
			return err
		}
		if err := d.Set(schemaName, b); err != nil {
			return err
		}
	}
	return nil
}

func resourceSnapshotDependencyRead(d *schema.ResourceData, meta interface{}) error {
	_, dt, err := getSnapshotDependency(meta.(*api.Client), d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			// This dependency was deleted out-of-band
			d.SetId("")
			return nil
		}
		return err
	}

	if err := d.Set("build_config_id", dt.BuildTypeID); err != nil {
		return err
	}

	if err := d.Set("source_build_config_id", dt.SourceBuildType.ID); err != nil {
		return err
	}

	return setSnapshotDependencyProperties(d, dt.Properties)
}

func resourceSnapshotDependencyDelete(d *schema.ResourceData, meta interface{}) error {
	depService, dt, err := getSnapshotDependency(meta.(*api.Client), d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			// This dependency was deleted out-of-band
			return nil
		}
		return err
	}

	return depService.DeleteSnapshot(dt.ID)
}

func getSnapshotDependency(client *api.Client, id string) (*api.DependencyService, *api.SnapshotDependency, error) {
	var buildConfigID, snapshotID string
	if n, err := fmt.Sscanf(id, snapshotIDFormat, &buildConfigID, &snapshotID); err != nil {
		return nil, nil, fmt.Errorf("invalid snapshot_dependency ID '%s' - %v", id, err)
	} else if n != 2 {
		return nil, nil, fmt.Errorf("invalid snapshot_dependency ID '%s' - Unrecognized format", id)
	}

	depService := client.DependencyService(buildConfigID)

	dt, err := depService.GetSnapshotByID(snapshotID)
	if err != nil {
		return nil, nil, err
	}

	return depService, dt, nil
}

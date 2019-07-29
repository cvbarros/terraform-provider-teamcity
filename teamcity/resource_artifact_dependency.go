package teamcity

import (
	"fmt"

	api "github.com/cvbarros/go-teamcity-sdk/teamcity"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceArtifactDependency() *schema.Resource {
	return &schema.Resource{
		Create: resourceArtifactDependencyCreate,
		Read:   resourceArtifactDependencyRead,
		Delete: resourceArtifactDependencyDelete,
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
				ForceNew: true,
				Required: true,
			},
			"dependency_revision": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  string(api.LatestSuccessfulBuild),
				ValidateFunc: validation.StringInSlice([]string{
					string(api.LatestFinishedBuild),
					string(api.LastBuildFinishedWithTag),
					string(api.LatestPinnedBuild),
					string(api.LatestSuccessfulBuild),
					string(api.BuildWithSpecifiedNumber),
					string(api.BuildFromSameChain),
				}, false),
			},
			"path_rules": {
				Type:     schema.TypeList,
				ForceNew: true,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"revision": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"clean_destination": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceArtifactDependencyCreate(d *schema.ResourceData, meta interface{}) error {
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
	opt, err := expandArtifactDependencyOptions(d)
	if err != nil {
		return err
	}
	dep, err := api.NewArtifactDependency(d.Get("source_build_config_id").(string), opt)
	if err != nil {
		return err
	}

	out, err := depService.AddArtifactDependency(dep)

	if err != nil {
		return err
	}

	d.SetId(out.ID())

	return resourceArtifactDependencyRead(d, meta)
}

func resourceArtifactDependencyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client).DependencyService(d.Get("build_config_id").(string))

	dt, err := getArtifactDependency(client, d.Id())
	if err != nil {
		return err
	}

	if err := d.Set("build_config_id", dt.BuildTypeID()); err != nil {
		return err
	}
	if dt.Options.CleanDestination {
		if err := d.Set("clean_destination", dt.Options.CleanDestination); err != nil {
			return err
		}
	}
	if err := d.Set("dependency_revision", string(dt.Options.ArtifactRevisionType)); err != nil {
		return err
	}
	if err := d.Set("path_rules", flattenStringSlice(dt.Options.PathRules)); err != nil {
		return err
	}
	if dt.Options.ArtifactRevisionType == api.BuildWithSpecifiedNumber || dt.Options.ArtifactRevisionType == api.LastBuildFinishedWithTag {
		if err := d.Set("revision", dt.Options.RevisionNumber); err != nil {
			return err
		}
	} else {
		d.Set("revision", nil)
	}

	return d.Set("source_build_config_id", dt.SourceBuildTypeID)
}

func resourceArtifactDependencyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	dep := client.DependencyService(d.Get("build_config_id").(string))

	return dep.DeleteArtifact(d.Id())
}

func getArtifactDependency(c *api.DependencyService, id string) (*api.ArtifactDependency, error) {

	dt, err := c.GetArtifactByID(id)
	if err != nil {
		return nil, err
	}

	return dt, nil
}

func expandArtifactDependencyOptions(d *schema.ResourceData) (*api.ArtifactDependencyOptions, error) {
	var cleanDestination bool
	var pathRules []string
	var revision string
	var revisionType api.ArtifactDependencyRevision
	if v, ok := d.GetOkExists("clean_destination"); ok {
		cleanDestination = v.(bool)
	}
	if v, ok := d.GetOk("path_rules"); ok {
		pathRules = expandStringSlice(v.([]interface{}))
	}
	if v, ok := d.GetOk("dependency_revision"); ok {
		revisionType = api.ArtifactDependencyRevision(v.(string))
	}

	if revisionType == api.LastBuildFinishedWithTag || revisionType == api.BuildWithSpecifiedNumber {
		if v, ok := d.GetOk("revision"); ok {
			revision = v.(string)
		} else {
			return nil, fmt.Errorf("'revision' property is required if using '%s' or '%s' for 'dependency_revision'", api.LastBuildFinishedWithTag, api.BuildWithSpecifiedNumber)
		}
	} else {
		//Ignore revison and remove from config/state if not used
		if err := d.Set("revision", nil); err != nil {
			return nil, err
		}
	}

	out, err := api.NewArtifactDependencyOptions(pathRules, revisionType, cleanDestination, revision)
	if err != nil {
		return nil, err
	}
	return out, nil
}

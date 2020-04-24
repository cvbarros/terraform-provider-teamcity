package teamcity

import (
	"fmt"
	"log"
	"strings"

	api "github.com/cvbarros/go-teamcity/teamcity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceFeatureGolang() *schema.Resource {
	return &schema.Resource{
		Create: resourceFeatureGolangCreate,
		Read:   resourceFeatureGolangRead,
		Delete: resourceFeatureGolangDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"build_config_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceFeatureGolangCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	buildConfigId := d.Get("build_config_id").(string)

	// validates the Build Configuration exists
	if _, err := client.BuildTypes.GetByID(buildConfigId); err != nil {
		return fmt.Errorf("invalid build_config_id %q - Build configuration does not exist", buildConfigId)
	}

	service := client.BuildFeatureService(buildConfigId)
	feature := api.NewFeatureGolang()
	feature.SetBuildTypeID(buildConfigId)
	createdService, err := service.Create(feature)
	if err != nil {
		return err
	}

	id := fmt.Sprintf("%s|%s", buildConfigId, createdService.ID())
	d.SetId(id)

	return resourceFeatureGolangRead(d, meta)
}

func resourceFeatureGolangRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	id, err := ParseFeatureGolangID(d.Id())
	if err != nil {
		return err
	}

	service := client.BuildFeatureService(id.BuildConfigID)
	if _, err := service.GetByID(id.FeatureID); err != nil {
		// handles this being deleted outside of TF
		if isNotFoundError(err) {
			log.Printf("[DEBUG] Build Feature Golang was not found - removing from state!")
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("build_config_id", id.BuildConfigID)

	return nil
}

func resourceFeatureGolangDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	id, err := ParseFeatureGolangID(d.Id())
	if err != nil {
		return err
	}

	service := client.BuildFeatureService(id.BuildConfigID)
	if err := service.Delete(id.FeatureID); err != nil {
		if !isNotFoundError(err) {
			return err
		}
	}

	return nil
}

type FeatureGolangId struct {
	BuildConfigID string
	FeatureID     string
}

func ParseFeatureGolangID(input string) (*FeatureGolangId, error) {
	// Format: 'BuildConfigID|FeatureID'
	segments := strings.Split(input, "|")
	if len(segments) != 2 {
		return nil, fmt.Errorf("Expected 2 segments but got %d", len(segments))
	}

	id := FeatureGolangId{
		BuildConfigID: segments[0],
		FeatureID:     segments[1],
	}
	return &id, nil
}

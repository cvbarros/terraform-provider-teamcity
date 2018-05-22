package teamcity

import (
	"fmt"
	"strings"

	api "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceFeatureCommitStatusPublisher() *schema.Resource {
	return &schema.Resource{
		Create: resourceFeatureCommitStatusPublisherCreate,
		Read:   resourceFeatureCommitStatusPublisherRead,
		Update: resourceFeatureCommitStatusPublisherUpdate,
		Delete: resourceFeatureCommitStatusPublisherDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"build_config_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"publisher": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"github"}, true),
			},
			"github": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auth_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"token", "password"}, true),
						},
						"host": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "https://api.github.com",
						},
						"username": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"password": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"access_token": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceFeatureCommitStatusPublisherCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	var buildConfigID string

	if v, ok := d.GetOk("build_config_id"); ok {
		buildConfigID = v.(string)
	}

	// validates the Build Configuration exists
	if _, err := client.BuildTypes.GetById(buildConfigID); err != nil {
		return fmt.Errorf("invalid build_config_id '%s' - Build configuration does not exist", buildConfigID)
	}

	srv := client.BuildFeatureService(buildConfigID)

	//Only Github publisher for now - Add support for more publishers later

	dt, err := buildGithubCommitStatusPublisher(d)
	if err != nil {
		return err
	}
	out, err := srv.Create(dt)

	if err != nil {
		return err
	}

	d.SetId(out.ID())

	return resourceFeatureCommitStatusPublisherRead(d, meta)
}

func resourceFeatureCommitStatusPublisherRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceFeatureCommitStatusPublisherUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceFeatureCommitStatusPublisherDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	svr := client.BuildFeatureService(d.Get("build_config_id").(string))

	return svr.Delete(d.Id())
}

func buildGithubCommitStatusPublisher(d *schema.ResourceData) (api.BuildFeature, error) {
	var opt api.StatusPublisherGithubOptions
	// MaxItems ensure at most 1 github element
	local := d.Get("github").(*schema.Set).List()[0].(map[string]interface{})
	host := local["host"].(string)
	authType := local["auth_type"].(string)
	switch strings.ToLower(authType) {
	case "token":
		opt = api.NewCommitStatusPublisherGithubOptionsToken(host, local["access_token"].(string))
	case "password":
		opt = api.NewCommitStatusPublisherGithubOptionsPassword(host, local["username"].(string), local["password"].(string))
	}

	return api.NewFeatureCommitStatusPublisherGithub(opt)
}

// func getBuildFeatureCommitPublisher(c *api.BuildFeatureService, id string) (*api.BuildFeature, error) {
// 	dt, err := c.GetByID(id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return dt, nil
// }

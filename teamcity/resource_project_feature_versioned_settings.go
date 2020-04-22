package teamcity

import (
	"fmt"
	"log"
	"strings"

	api "github.com/cvbarros/go-teamcity/teamcity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceProjectFeatureVersionedSettings() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectFeatureVersionedSettingsCreate,
		Read:   resourceProjectFeatureVersionedSettingsRead,
		Update: resourceProjectFeatureVersionedSettingsUpdate,
		Delete: resourceProjectFeatureVersionedSettingsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vcs_root_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"build_settings": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(api.VersionedSettingsBuildSettingsAlwaysUseCurrent),
					string(api.VersionedSettingsBuildSettingsPreferCurrent),
					string(api.VersionedSettingsBuildSettingsPreferVcs),
				}, false),
			},

			"format": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(api.VersionedSettingsFormatKotlin),
					string(api.VersionedSettingsFormatXML),
				}, false),
			},

			// TODO: once https://github.com/cvbarros/go-teamcity/pull/76 is available
			//"credentials_storage_type": {
			//	Type:     schema.TypeString,
			//	Optional: true,
			//	ValidateFunc: validation.StringInSlice([]string{
			//		string(api.CredentialsStorageTypeCredentialsJSON),
			//		"scrambled", // isn't returned, this is a fake value for convenience in configs
			//	}, false),
			//},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"show_changes": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"use_relative_ids": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceProjectFeatureVersionedSettingsCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	projectId := d.Get("project_id").(string)
	service := client.ProjectFeatureService(projectId)

	feature := api.NewProjectFeatureVersionedSettings(projectId, api.ProjectFeatureVersionedSettingsOptions{
		BuildSettings:  api.VersionedSettingsBuildSettings(d.Get("build_settings").(string)),
		Enabled:        d.Get("enabled").(bool),
		Format:         api.VersionedSettingsFormat(d.Get("format").(string)),
		ShowChanges:    d.Get("show_changes").(bool),
		UseRelativeIds: d.Get("use_relative_ids").(bool),
		VcsRootID:      d.Get("vcs_root_id").(string),
	})

	//if v := d.Get("credentials_storage_type").(string); v == string(api.CredentialsStorageTypeCredentialsJSON) {
	//	feature.Options.CredentialsStorageType = api.CredentialsStorageTypeCredentialsJSON
	//}

	createdFeature, err := service.Create(feature)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s|%s", projectId, createdFeature.ID()))

	return resourceProjectFeatureVersionedSettingsRead(d, meta)
}

func resourceProjectFeatureVersionedSettingsUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	id, err := ParseProjectFeatureId(d.Id())
	if err != nil {
		return err
	}

	service := client.ProjectFeatureService(id.ProjectId)
	feature, err := service.GetByID(id.FeatureId)
	if err != nil {
		return err
	}

	vcsFeature, ok := feature.(*api.ProjectFeatureVersionedSettings)
	if !ok {
		return fmt.Errorf("Expected a VersionedSettings Feature but wasn't!")
	}

	if d.HasChange("build_settings") {
		vcsFeature.Options.BuildSettings = api.VersionedSettingsBuildSettings(d.Get("build_settings").(string))
	}
	//if d.HasChange("credentials_storage_type") {
	//	v := d.Get("credentials_storage_type").(string)
	//	if v == string(api.CredentialsStorageTypeCredentialsJSON) {
	//		feature.Options.CredentialsStorageType = api.CredentialsStorageTypeCredentialsJSON
	//	} else {
	//		feature.Options.CredentialsStorageType = api.CredentialsStorageTypeScrambledInVcs
	//	}
	//}
	if d.HasChange("enabled") {
		vcsFeature.Options.Enabled = d.Get("enabled").(bool)
	}
	if d.HasChange("format") {
		vcsFeature.Options.Format = api.VersionedSettingsFormat(d.Get("format").(string))
	}
	if d.HasChange("show_changes") {
		vcsFeature.Options.ShowChanges = d.Get("show_changes").(bool)
	}
	if d.HasChange("use_relative_ids") {
		vcsFeature.Options.UseRelativeIds = d.Get("use_relative_ids").(bool)
	}
	if d.HasChange("vcs_root_id") {
		vcsFeature.Options.VcsRootID = d.Get("vcs_root_id").(string)
	}

	if _, err := service.Update(vcsFeature); err != nil {
		return err
	}

	return resourceProjectFeatureVersionedSettingsRead(d, meta)
}

func resourceProjectFeatureVersionedSettingsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	id, err := ParseProjectFeatureId(d.Id())
	if err != nil {
		return err
	}

	service := client.ProjectFeatureService(id.ProjectId)
	feature, err := service.GetByID(id.FeatureId)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Printf("[DEBUG] Project Feature Versioned Settings was not found - removing from state!")
			d.SetId("")
			return nil
		}

		return err
	}

	vcsFeature, ok := feature.(*api.ProjectFeatureVersionedSettings)
	if !ok {
		return fmt.Errorf("Expected a VersionedSettings Feature but wasn't!")
	}

	d.Set("build_settings", string(vcsFeature.Options.BuildSettings))
	d.Set("enabled", vcsFeature.Options.Enabled)
	d.Set("format", string(vcsFeature.Options.Format))
	d.Set("project_id", id.ProjectId)
	d.Set("show_changes", vcsFeature.Options.ShowChanges)
	d.Set("use_relative_ids", vcsFeature.Options.UseRelativeIds)
	d.Set("vcs_root_id", vcsFeature.Options.VcsRootID)

	//credentialsStorageType := "scrambled"
	//if vcsFeature.Options.CredentialsStorageType != "" {
	//	credentialsStorageType = string(vcsFeature.Options.CredentialsStorageType)
	//}
	//d.Set("credentials_storage_type", credentialsStorageType)

	return nil
}

func resourceProjectFeatureVersionedSettingsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	id, err := ParseProjectFeatureId(d.Id())
	if err != nil {
		return err
	}

	service := client.ProjectFeatureService(id.ProjectId)
	return service.Delete(id.FeatureId)
}

type projectFeatureId struct {
	FeatureId string
	ProjectId string
}

func ParseProjectFeatureId(input string) (*projectFeatureId, error) {
	segments := strings.Split(input, "|")
	if len(segments) != 2 {
		return nil, fmt.Errorf("Expected a string in the format 'ProjectID|FeatureID' but got %q", input)
	}

	// format is: 'ProjectID|FeatureID'
	id := projectFeatureId{
		ProjectId: segments[0],
		FeatureId: segments[1],
	}
	return &id, nil
}

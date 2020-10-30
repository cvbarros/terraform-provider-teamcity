package teamcity

import (
	"fmt"
	"log"
	"strings"

	api "github.com/cvbarros/go-teamcity/teamcity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceProjectFeatureVaultConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectFeatureVaultConnectionCreate,
		Read:   resourceProjectFeatureVaultConnectionRead,
		Update: resourceProjectFeatureVaultConnectionUpdate,
		Delete: resourceProjectFeatureVaultConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"auth_method": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(api.ConnectionProviderVaultAuthMethodApprole),
				ValidateFunc: validation.StringInSlice([]string{
					string(api.ConnectionProviderVaultAuthMethodIAM),
					string(api.ConnectionProviderVaultAuthMethodApprole),
				}, false),
			},

			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Hashicorp Vault",
			},

			"fail_on_error": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"role_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"secret_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Sensitive:    true,
			},

			"url": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},

			// TODO: Add missing params from api.ConnectionProviderVaultOptions
		},
	}
}

func resourceProjectFeatureVaultConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	projectID := d.Get("project_id").(string)
	service := client.ProjectFeatureService(projectID)

	feature := api.NewProjectConnectionVault(projectID, api.ConnectionProviderVaultOptions{
		AuthMethod:  api.ConnectionProviderVaultAuthMethod(d.Get("auth_method").(string)),
		DisplayName: d.Get("display_name").(string),
		URL:         d.Get("url").(string),
		// TOOO: Add full range of params here
	})

	if v := d.Get("auth_method").(string); v == string(api.ConnectionProviderVaultAuthMethodApprole) {
		feature.Options.RoleID = d.Get("role_id").(string)
		feature.Options.SecretID = d.Get("secret_id").(string)
	}

	// however the ID returned eventually gets overwritten
	// so we need to look it up using the type
	if _, err := service.Create(feature); err != nil {
		return err
	}

	d.SetId(projectID)

	return resourceProjectFeatureVaultConnectionRead(d, meta)
}

// TODO: Support updating resources
func resourceProjectFeatureVaultConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	// client := meta.(*api.Client)

	// projectID := d.Id()
	// service := client.ProjectFeatureService(projectID)
	// feature, err := service.GetByType("versionedSettings")
	// if err != nil {
	// 	return err
	// }

	// vcsFeature, ok := feature.(*api.ProjectFeatureVersionedSettings)
	// if !ok {
	// 	return fmt.Errorf("Expected a VersionedSettings Feature but wasn't")
	// }

	// if d.HasChange("build_settings") {
	// 	vcsFeature.Options.BuildSettings = api.VersionedSettingsBuildSettings(d.Get("build_settings").(string))
	// }
	// if d.HasChange("context_parameters") {
	// 	contextParametersRaw := d.Get("context_parameters").(map[string]interface{})
	// 	vcsFeature.Options.ContextParameters = expandContextParameters(contextParametersRaw)
	// }
	// if d.HasChange("credentials_storage_type") {
	// 	v := d.Get("credentials_storage_type").(string)
	// 	if v == string(api.CredentialsStorageTypeCredentialsJSON) {
	// 		vcsFeature.Options.CredentialsStorageType = api.CredentialsStorageTypeCredentialsJSON
	// 	} else {
	// 		vcsFeature.Options.CredentialsStorageType = api.CredentialsStorageTypeScrambledInVcs
	// 	}
	// }
	// if d.HasChange("enabled") {
	// 	vcsFeature.Options.Enabled = d.Get("enabled").(bool)
	// }
	// if d.HasChange("format") {
	// 	vcsFeature.Options.Format = api.VersionedSettingsFormat(d.Get("format").(string))
	// }
	// if d.HasChange("show_changes") {
	// 	vcsFeature.Options.ShowChanges = d.Get("show_changes").(bool)
	// }
	// if d.HasChange("use_relative_ids") {
	// 	vcsFeature.Options.UseRelativeIds = d.Get("use_relative_ids").(bool)
	// }
	// if d.HasChange("vcs_root_id") {
	// 	vcsFeature.Options.VcsRootID = d.Get("vcs_root_id").(string)
	// }

	// if _, err := service.Update(vcsFeature); err != nil {
	// 	return err
	// }

	return resourceProjectFeatureVaultConnectionRead(d, meta)
}

func resourceProjectFeatureVaultConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	projectID := d.Id()
	service := client.ProjectFeatureService(projectID)
	feature, err := service.GetByTypeAndProvider("OAuthProvider", "teamcity-vault")
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Printf("[DEBUG] Project Feature Vault Connection was not found - removing from state")
			d.SetId("")
			return nil
		}

		return err
	}

	vaultFeature, ok := feature.(*api.ConnectionProviderVault)
	if !ok {
		return fmt.Errorf("Expected a ConnectionProviderVault Feature but wasn't")
	}

	d.Set("project_id", projectID)
	d.Set("url", string(vaultFeature.Options.URL))
	// TODO: Add full set of resource params here.

	return nil
}

func resourceProjectFeatureVaultConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	projectID := d.Id()
	service := client.ProjectFeatureService(projectID)
	feature, err := service.GetByTypeAndProvider("OAuthProvider", "teamcity-vault")
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			// already gone
			return nil
		}

		return err
	}

	return service.Delete(feature.ID())
}

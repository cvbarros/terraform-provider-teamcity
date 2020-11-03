package teamcity

import (
	"errors"
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

			"approle_auth_path": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "approle",
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"approle_role_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"approle_secret_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Sensitive:    true,
			},

			"auth_method": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(api.ConnectionProviderVaultAuthMethodApprole),
				ValidateFunc: validation.StringInSlice([]string{
					string(api.ConnectionProviderVaultAuthMethodIAM),
					string(api.ConnectionProviderVaultAuthMethodApprole),
				}, false),
				Description: "Use Approle or AWS IAM Auth method",
			},

			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Hashicorp Vault",
			},

			"fail_on_error": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Should builds fail in the case of an error resolving parameters",
			},

			"paramater_namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Namespace for use in TeamCity parameters in case of multiple Vault connections",
			},

			"url": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},

			"vault_namespace": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceProjectFeatureVaultConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	projectID := d.Get("project_id").(string)
	service := client.ProjectFeatureService(projectID)

	feature := api.NewProjectConnectionVault(projectID, api.ConnectionProviderVaultOptions{
		AuthMethod:     api.ConnectionProviderVaultAuthMethod(d.Get("auth_method").(string)),
		DisplayName:    d.Get("display_name").(string),
		FailOnError:    d.Get("fail_on_error").(bool),
		URL:            d.Get("url").(string),
		VaultNamespace: d.Get("vault_namespace").(string),
	})

	if v := d.Get("auth_method").(string); v == string(api.ConnectionProviderVaultAuthMethodApprole) {
		approleAuthPath, approleAuthPathOk := d.GetOk("approle_auth_path")
		if approleAuthPathOk {
			feature.Options.Endpoint = approleAuthPath.(string)
		}

		approleRoleID, approleRoleIDOk := d.GetOk("approle_role_id")
		approleSecretID, approleSecretIDOk := d.GetOk("approle_secret_id")

		if !approleRoleIDOk || !approleSecretIDOk {
			return errors.New("both approle_role_id and approle_secret_id must be supplied when using approle auth")
		}

		feature.Options.RoleID = approleRoleID.(string)
		feature.Options.SecretID = approleSecretID.(string)
	}

	if v, ok := d.GetOk("namespace"); ok {
		feature.Options.Namespace = v.(string)
	}

	if v, ok := d.GetOk("vault_namespace"); ok {
		feature.Options.VaultNamespace = v.(string)
	}

	// however the ID returned eventually gets overwritten
	// so we need to look it up using the type
	if _, err := service.Create(feature); err != nil {
		return err
	}

	d.SetId(projectID)

	return resourceProjectFeatureVaultConnectionRead(d, meta)
}

func resourceProjectFeatureVaultConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	projectID := d.Id()
	service := client.ProjectFeatureService(projectID)
	feature, err := service.GetByTypeAndProvider("OAuthProvider", "teamcity-vault")
	if err != nil {
		return err
	}

	vaultFeature, ok := feature.(*api.ConnectionProviderVault)
	if !ok {
		return fmt.Errorf("Expected a ConnectionProviderVault Feature but wasn't")
	}

	if d.HasChange("auth_method") {
		vaultFeature.Options.AuthMethod = api.ConnectionProviderVaultAuthMethod(d.Get("auth_method").(string))
	}

	if d.HasChange("approle_auth_path") {
		vaultFeature.Options.Endpoint = d.Get("approle_auth_path").(string)
	}

	if d.HasChange("approle_role_id") {
		vaultFeature.Options.RoleID = d.Get("approle_role_id").(string)
	}

	if d.HasChange("approle_secret_id") {
		vaultFeature.Options.SecretID = d.Get("approle_secret_id").(string)
	}

	if d.HasChange("fail_on_error") {
		vaultFeature.Options.FailOnError = d.Get("fail_on_error").(bool)
	}

	if d.HasChange("display_name") {
		vaultFeature.Options.DisplayName = d.Get("display_name").(string)
	}

	if d.HasChange("parameter_namespace") {
		vaultFeature.Options.Namespace = d.Get("parameter_namespace").(string)
	}

	if d.HasChange("vault_namespace") {
		vaultFeature.Options.VaultNamespace = d.Get("vault_namespace").(string)
	}

	if d.HasChange("url") {
		vaultFeature.Options.URL = d.Get("url").(string)
	}

	if _, err := service.Update(vaultFeature); err != nil {
		return err
	}

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
	d.Set("auth_method", vaultFeature.Options.AuthMethod)
	d.Set("display_name", vaultFeature.Options.DisplayName)
	d.Set("fail_on_error", vaultFeature.Options.FailOnError)
	d.Set("namespace", vaultFeature.Options.Namespace)
	d.Set("url", string(vaultFeature.Options.URL))
	d.Set("vault_namespace", vaultFeature.Options.VaultNamespace)

	if vaultFeature.Options.AuthMethod == api.ConnectionProviderVaultAuthMethodApprole {
		d.Set("approle_auth_path", string(api.ConnectionProviderVaultAuthMethodApprole))
		d.Set("approle_role_id", vaultFeature.Options.RoleID)
		// approle_secret_id is a non-readable field via API, so can only be created or updated.
		d.Set("approle_secret_id", d.Get("approle_secret_id").(string))
	}

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

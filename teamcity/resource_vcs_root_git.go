package teamcity

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	api "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceVcsRootGit() *schema.Resource {
	return &schema.Resource{
		Create: resourceVcsRootGitCreate,
		Read:   resourceVcsRootGitRead,
		Update: resourceVcsRootGitUpdate,
		Delete: resourceVcsRootGitDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name to identify this Git VCS Root.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID for the parent project for this VCS Root. Required.",
			},
			"fetch_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "URL used to pull source code for this VCS. For HTTP, prefix with http(s)://. For SSH, use user@server.com.",
			},
			"push_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "URL used to push if needed. Assumes the same as fetch_url if not specified.",
			},
			"default_branch": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Branch specification for the default branch to pull/push from/to and inspec changes. Ex: refs/head/master",
			},
			"enable_branch_spec_tags": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Branch specification for the default branch to pull/push from/to and inspec changes. Ex: refs/head/master",
			},
			"submodule_checkout": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "checkout",
				Description: "Branch specification for the default branch to pull/push from/to and inspec changes. Ex: refs/head/master",
				StateFunc: func(v interface{}) string {
					value := v.(string)
					return strings.ToUpper(value)
				},
			},
			"username_style": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "userid",
				ValidateFunc: validation.StringInSlice([]string{"userid", "author_name", "author_email", "author_full"}, true),
				Description:  "Branch specification for the default branch to pull/push from/to and inspec changes. Ex: refs/head/master",
			},
			"auth": {
				Type:        schema.TypeSet,
				MaxItems:    1,
				Optional:    true,
				Description: "Authentication configuration for VCS Root. If not specified, defaults to anonymous auth.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"userpass", "ssh", "anonymous"}, true),
							Required:     true,
						},
						"username": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ssh_type": {
							Type:     schema.TypeString,
							Optional: true,
							//ValidateFunc: validation.StringInSlice([]string{"uploadedKey", "customKey", "defaultKey"}, false),
							Description: "If using SSH, this field specifies how the SSH Key will be sourced.",
						},
						"key_spec": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.NoZeroValues,
							Description:  "For 'customKey' refers to the path to a private key. For 'uploadedKey', corresponds to the name of the SSH Key uploaded into the project. Required if using 'customKey' or 'uploadedKey'.",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Computed:    true,
							Description: "Password if using 'userpass' auth. Private key passphrase if using 'uploadedKey' or 'customKey'. Required if not anonymous auth.",
						},
					},
				},
				Set: gitVcsAuthHash,
			},
		},
	}
}

var expandUsernameStyleMap = map[string]string{
	"userid":       string(api.GitVcsUsernameStyleUserID),
	"author_email": string(api.GitVcsUsernameStyleAuthorEmail),
	"author_name":  string(api.GitVcsUsernameStyleAuthorName),
	"author_full":  string(api.GitVcsUsernameStyleAuthorNameAndEmail),
}

var flattenUsernameStyleMap = reverseMap(expandUsernameStyleMap)

func resourceVcsRootGitCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	projectID := d.Get("project_id").(string)
	var gitVcs *api.GitVcsRoot

	vcsOpts, err := expandGitVcsRootOptions(d)
	if err != nil {
		return err
	}

	if v, ok := d.GetOk("username_style"); ok {
		vcsOpts.UsernameStyle = api.GitVcsUsernameStyle(expandUsernameStyleMap[v.(string)])
	}

	if v, ok := d.GetOk("enable_branch_spec_tags"); ok {
		vcsOpts.EnableTagsInBranchSpec = v.(bool)
	}

	if gitVcs, err = api.NewGitVcsRoot(projectID, d.Get("name").(string), vcsOpts); err != nil {
		return err
	}

	created, err := client.VcsRoots.Create(projectID, gitVcs)
	if err != nil {
		return err
	}

	d.MarkNewResource()
	d.SetId(created.ID)
	return resourceVcsRootGitRead(d, meta)
}

func resourceVcsRootGitRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	vcsID := d.Id()

	vcs, err := client.VcsRoots.GetByID(vcsID)
	if err != nil {
		return err
	}

	dt, ok := vcs.(*api.GitVcsRoot)
	if !ok {
		return fmt.Errorf("VCS with ID = %s has a type mismatch, not a Git VCS. Actual type: %s", vcsID, vcs.VcsName())
	}

	if err := d.Set("project_id", dt.Project.ID); err != nil {
		return err
	}

	if err := d.Set("name", dt.Name); err != nil {
		return err
	}

	if err := d.Set("fetch_url", dt.Options.FetchURL); err != nil {
		return err
	}

	if err := d.Set("push_url", dt.Options.PushURL); err != nil {
		return err
	}

	if err := d.Set("default_branch", dt.Options.DefaultBranch); err != nil {
		return err
	}

	if auth, err := flattenGitVcsRootAuth(d, dt.Options); err != nil {
		if err := d.Set("auth", auth); err != nil {
			return err
		}
	}

	if err := d.Set("submodule_checkout", dt.Options.SubModuleCheckout); err != nil {
		return err
	}

	if err := d.Set("enable_branch_spec_tags", dt.Options.EnableTagsInBranchSpec); err != nil {
		return err
	}

	if d.HasChange("username_style") {
		if err := d.Set("username_style", flattenUsernameStyleMap[string(dt.Options.UsernameStyle)]); err != nil {
			return err
		}
	}

	return nil
}

func expandGitVcsRootOptions(d *schema.ResourceData) (*api.GitVcsRootOptions, error) {
	authType, err := getGitAuthType(d)
	var username, password, fetchURL, pushURL string
	if err != nil {
		return nil, err
	}

	if v, ok := d.GetOk("push_url"); ok {
		pushURL = v.(string)
	}

	if v, ok := d.GetOk("fetch_url"); ok {
		fetchURL = v.(string)
	}

	_, authSpecified := d.GetOkExists("auth")
	if !authSpecified {
		opt, err := api.NewGitVcsRootOptions(d.Get("default_branch").(string), fetchURL, pushURL, authType, username, password)
		if err != nil {
			return nil, err
		}
		return opt, nil
	}

	// Only 1 max permitted
	auth := d.Get("auth").(*schema.Set).List()[0].(map[string]interface{})

	if authType != api.GitAuthMethodAnonymous {
		username = auth["username"].(string)
	}

	if authType == api.GitAuthSSHUploadedKey || authType == api.GitAuthSSHCustomKey || authType == api.GitAuthMethodPassword {
		password = auth["password"].(string)
	}

	opt, err := api.NewGitVcsRootOptions(d.Get("default_branch").(string), fetchURL, pushURL, authType, username, password)
	if err != nil {
		return nil, err
	}
	if v, ok := auth["key_spec"]; ok {
		opt.PrivateKeySource = v.(string)
	}
	return opt, nil
}

func flattenGitVcsRootAuth(d *schema.ResourceData, dt *api.GitVcsRootOptions) ([]map[string]interface{}, error) {
	var optsToSave []map[string]interface{}
	m := make(map[string]interface{})
	authType, sshType := readAuthTypeFromAuthMethod(dt)

	if authType == "" {
		return nil, fmt.Errorf("invalid auth method returned from api: '%s'", dt.AuthMethod)
	}

	m["type"] = authType
	if authType == "ssh" {
		m["ssh_type"] = sshType
		if dt.PrivateKeySource != "" {
			m["key_spec"] = dt.PrivateKeySource
		}
	}

	if dt.Username != "" {
		m["username"] = dt.Username
	}

	//Set back password if contained in state
	if auth, ok := d.GetOk("auth"); ok {
		authSet := auth.(*schema.Set).List()[0].(map[string]interface{})
		if pwd, ok := authSet["password"]; ok {
			m["password"] = pwd.(string)
		}
	}

	optsToSave = append(optsToSave, m)
	return optsToSave, nil
}

func resourceVcsRootGitUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVcsRootGitDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	return client.VcsRoots.Delete(d.Id())
}

func getGitAuthType(d *schema.ResourceData) (api.GitAuthMethod, error) {

	// If no auth specified, assume "anonymous"
	authConfig, ok := d.GetOk("auth")
	if !ok {
		return api.GitAuthMethodAnonymous, nil
	}

	// Only 1 max permitted
	auth := authConfig.(*schema.Set).List()[0].(map[string]interface{})
	authType := auth["type"].(string)

	switch authType {
	case "userpass":
		return api.GitAuthMethodPassword, nil
	case "anonymous":
		return api.GitAuthMethodAnonymous, nil
	case "ssh":
		sshType := auth["ssh_type"].(string)
		switch sshType {
		case "customKey":
			return api.GitAuthSSHCustomKey, nil
		case "uploadedKey":
			return api.GitAuthSSHUploadedKey, nil
		case "defaultKey":
			return api.GitAuthSSHDefaultKey, nil
		}
	}

	return "", fmt.Errorf("unsupported auth type: %s", authType)
}

func readAuthTypeFromAuthMethod(vcsOpt *api.GitVcsRootOptions) (string, string) {
	switch vcsOpt.AuthMethod {
	case api.GitAuthMethodAnonymous:
		return "anonymous", ""
	case api.GitAuthMethodPassword:
		return "userpass", ""
	case api.GitAuthSSHCustomKey:
		return "ssh", "customKey"
	case api.GitAuthSSHDefaultKey:
		return "ssh", "defaultKey"
	case api.GitAuthSSHUploadedKey:
		return "ssh", "uploadedKey"
	}

	return "", ""
}

func gitVcsAuthHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["type"].(string)))
	hash := hashcode.String(buf.String())

	log.Printf("[DEBUG] GitVcs Auth Hash: %d", hash)
	return hash
}

func reverseMap(m map[string]string) map[string]string {
	n := make(map[string]string)
	for k, v := range m {
		n[v] = k
	}
	return n
}

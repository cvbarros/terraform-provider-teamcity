package teamcity

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	api "github.com/cvbarros/go-teamcity-sdk/teamcity"
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
			"modification_check_interval": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntAtLeast(1),
				Optional:     true,
				Description:  "Specifies how often TeamCity polls the VCS repository for VCS changes (in seconds)",
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
			"branches": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Branches to monitor besides the default with a set of rules in the form of +|-:branch_name (with the optional * placeholder)",
			},
			"enable_branch_spec_tags": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, tags can be used in the branch specification",
			},
			"submodule_checkout": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "checkout",
				ValidateFunc: validation.StringInSlice([]string{"checkout", "ignore"}, true),
				Description:  "Defines whether to checkout submodules. Use either 'checkout' or 'ignore'.",
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
				Description:  "Defines a way TeamCity binds VCS changes to the user. Changing username style will affect only newly collected changes. Old changes will continue to be stored with the style that was active at the time of collecting changes. Allowed values: 'userid', 'author_name', 'author_email', 'author_full'",
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
			"agent": {
				Type:        schema.TypeSet,
				MaxItems:    1,
				Optional:    true,
				Description: "Agent settings for the VCS Root",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"git_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"clean_policy": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"branch_change", "always", "never"}, false),
						},
						"clean_files_policy": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"untracked", "ignored_only", "non_ignored_only"}, false),
						},
						"use_mirrors": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
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

var expandCleanPolicyMap = map[string]string{
	"branch_change": string(api.CleanPolicyBranchChange),
	"always":        string(api.CleanPolicyAlways),
	"never":         string(api.CleanPolicyNever),
}

var flattenCleanPolicyMap = reverseMap(expandCleanPolicyMap)

var expandCleanFilesPolicyMap = map[string]string{
	"untracked":        string(api.CleanFilesPolicyAllUntracked),
	"ignored_only":     string(api.CleanFilesPolicyIgnoredOnly),
	"non_ignored_only": string(api.CleanFilesPolicyIgnoredUntracked),
}

var flattenCleanFilesPolicyMap = reverseMap(expandCleanFilesPolicyMap)

func resourceVcsRootGitCreate(d *schema.ResourceData, meta interface{}) error {
	d.MarkNewResource()
	return resourceVcsRootGitUpdate(d, meta)
}

func resourceVcsRootGitUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	projectID := d.Get("project_id").(string)
	var gitVcs *api.GitVcsRoot
	var name string
	var modificationCheckInterval int

	vcsOpts, err := expandGitVcsRootOptions(d)

	if err != nil {
		return err
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}
	if v, ok := d.GetOk("username_style"); ok {
		vcsOpts.UsernameStyle = api.GitVcsUsernameStyle(expandUsernameStyleMap[v.(string)])
	}

	if v, ok := d.GetOk("enable_branch_spec_tags"); ok {
		vcsOpts.EnableTagsInBranchSpec = v.(bool)
	}

	if v, ok := d.GetOk("modification_check_interval"); ok {
		modificationCheckInterval = v.(int)
	}

	if d.IsNewResource() {
		log.Printf("[INFO] detected new VCS Root resource, creating.")
		if gitVcs, err = api.NewGitVcsRoot(projectID, name, vcsOpts); err != nil {
			return err
		}
		if modificationCheckInterval > 0 {
			gitVcs.SetModificationCheckInterval(int32(modificationCheckInterval))
		}
		created, err := client.VcsRoots.Create(projectID, gitVcs)
		if err != nil {
			return err
		}
		d.SetId(created.ID)

		return resourceVcsRootGitRead(d, meta)
	}
	log.Printf("[INFO] Updating VCS Root resource.")
	vcs, err := client.VcsRoots.GetByID(d.Id())
	if err != nil {
		return err
	}

	gitVcs = vcs.(*api.GitVcsRoot)
	if d.HasChange("name") {
		gitVcs.SetName(name)
	}
	if d.HasChange("project_id") {
		gitVcs.SetProjectID(projectID)
	}
	if d.HasChange("modification_check_interval") {
		gitVcs.SetModificationCheckInterval(int32(modificationCheckInterval))
	}
	if d.HasChange("auth") {
		gitVcs.Options.AuthMethod = vcsOpts.AuthMethod
		gitVcs.Options.Password = vcsOpts.Password
		gitVcs.Options.PrivateKeySource = vcsOpts.PrivateKeySource
		gitVcs.Options.Username = vcsOpts.Username
	}

	if d.HasChange("agent") {
		gitVcs.Options.AgentSettings = vcsOpts.AgentSettings
	}

	if d.HasChange("username_style") {
		gitVcs.Options.UsernameStyle = vcsOpts.UsernameStyle
	}

	if d.HasChange("enable_branch_spec_tags") {
		gitVcs.Options.EnableTagsInBranchSpec = vcsOpts.EnableTagsInBranchSpec
	}

	if d.HasChange("push_url") {
		gitVcs.Options.PushURL = vcsOpts.PushURL
	}
	if d.HasChange("fetch_url") {
		gitVcs.Options.FetchURL = vcsOpts.FetchURL
	}
	if d.HasChange("branches") {
		gitVcs.Options.BranchSpec = vcsOpts.BranchSpec
	}
	if d.HasChange("default_branch") {
		gitVcs.Options.DefaultBranch = vcsOpts.DefaultBranch
	}
	if d.HasChange("submodule_checkout") {
		gitVcs.Options.SubModuleCheckout = vcsOpts.SubModuleCheckout
	}

	_, err = client.VcsRoots.Update(gitVcs)
	if err != nil {
		return err
	}

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

	if err := d.Set("name", dt.Name()); err != nil {
		return err
	}

	if dt.ModificationCheckInterval() != nil {
		v := *(dt.ModificationCheckInterval())
		if err := d.Set("modification_check_interval", int(v)); err != nil {
			return err
		}
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

	if len(dt.Options.BranchSpec) > 0 {
		if err := d.Set("branches", dt.Options.BranchSpec); err != nil {
			return err
		}
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

	if agent, err := flattenGitAgentSettings(d, dt.Options.AgentSettings); err != nil {
		if err := d.Set("agent", agent); err != nil {
			return err
		}
	}

	return nil
}

func resourceVcsRootGitDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	return client.VcsRoots.Delete(d.Id())
}

func expandGitVcsRootOptions(d *schema.ResourceData) (*api.GitVcsRootOptions, error) {
	var username, password, fetchURL, pushURL, privateKeySource string
	var opt *api.GitVcsRootOptions

	if v, ok := d.GetOk("push_url"); ok {
		pushURL = v.(string)
	}

	if v, ok := d.GetOk("fetch_url"); ok {
		fetchURL = v.(string)
	}

	authType, err := getGitAuthType(d)
	if err != nil {
		return nil, err
	}

	_, authSpecified := d.GetOkExists("auth")
	if authSpecified {
		// Only 1 max permitted
		auth := d.Get("auth").(*schema.Set).List()[0].(map[string]interface{})

		if authType != api.GitAuthMethodAnonymous {
			username = auth["username"].(string)
		}

		if authType == api.GitAuthSSHUploadedKey || authType == api.GitAuthSSHCustomKey || authType == api.GitAuthMethodPassword {
			password = auth["password"].(string)
		}

		if v, ok := auth["key_spec"]; ok {
			privateKeySource = v.(string)
		}
	}

	agent, err := expandGitVcsAgentSettings(d)
	if err != nil {
		return nil, err
	}

	opt, err = api.NewGitVcsRootOptionsWithAgentSettings(d.Get("default_branch").(string), fetchURL, pushURL, authType, username, password, agent)
	if err != nil {
		return nil, err
	}

	if v, ok := d.GetOk("branches"); ok {
		opt.BranchSpec = expandStringSlice(v.([]interface{}))
	}

	if v, ok := d.GetOk("submodule_checkout"); ok {
		opt.SubModuleCheckout = strings.ToUpper(v.(string))
	}
	if privateKeySource != "" {
		opt.PrivateKeySource = privateKeySource
	}

	return opt, nil
}

func expandGitVcsAgentSettings(d *schema.ResourceData) (*api.GitAgentSettings, error) {
	var gitPath, cleanPolicy, cleanFilesPolicy string
	var useMirrors bool

	v, ok := d.GetOk("agent")
	if !ok {
		return nil, nil
	}

	agent := v.(*schema.Set).List()[0].(map[string]interface{})

	if v, ok = agent["git_path"]; ok {
		gitPath = v.(string)
	}

	if v, ok = agent["clean_policy"]; ok {
		cleanPolicy = expandCleanPolicyMap[v.(string)]
	}

	if v, ok = agent["clean_files_policy"]; ok {
		cleanFilesPolicy = expandCleanFilesPolicyMap[v.(string)]
	}

	if v, ok = agent["use_mirrors"]; ok {
		useMirrors = v.(bool)
	}

	return &api.GitAgentSettings{
		GitPath:          gitPath,
		CleanFilesPolicy: api.GitAgentCleanFilesPolicy(cleanFilesPolicy),
		CleanPolicy:      api.GitAgentCleanPolicy(cleanPolicy),
		UseMirrors:       useMirrors,
	}, nil
}

func flattenGitAgentSettings(d *schema.ResourceData, dt *api.GitAgentSettings) ([]map[string]interface{}, error) {
	if dt == nil {
		return nil, nil
	}
	var optsToSave []map[string]interface{}
	m := make(map[string]interface{})

	if dt.GitPath != "" {
		m["git_path"] = dt.GitPath
	}
	if dt.CleanPolicy != "" {
		m["clean_policy"] = flattenCleanPolicyMap[string(dt.CleanPolicy)]
	}
	if dt.CleanFilesPolicy != "" {
		m["clean_files_policy"] = flattenCleanFilesPolicyMap[string(dt.CleanFilesPolicy)]
	}

	m["use_mirrors"] = dt.UseMirrors

	optsToSave = append(optsToSave, m)
	return optsToSave, nil
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

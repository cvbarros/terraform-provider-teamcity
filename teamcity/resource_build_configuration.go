package teamcity

import (
	"bytes"
	"fmt"
	"strings"

	api "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceBuildConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceBuildConfigurationCreate,
		Read:   resourceBuildConfigurationRead,
		Update: resourceBuildConfigurationUpdate,
		Delete: resourceBuildConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vcs_root": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"checkout_rules": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
				Set: vcsRootHash,
			},
			"step": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateStepType(),
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"file": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"args": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Set: stepSetHash,
			},
			"env_params": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
			"config_params": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
			"sys_params": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceBuildConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	bt := &api.BuildType{}
	var projectID string

	props := api.NewProperties()

	if v, ok := d.GetOk("project_id"); ok {
		projectID = v.(string)
		bt.ProjectID = projectID
	}

	if v, ok := d.GetOk("name"); ok {
		bt.Name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		bt.Description = v.(string)
	}

	if v, ok := d.GetOk("env_params"); ok {
		envParams := v.(map[string]interface{})
		for k, v := range envParams {
			p := &api.Property{
				Name:  fmt.Sprintf("env.%s", k),
				Value: v.(string),
			}
			props.Add(p)
		}
	}

	if v, ok := d.GetOk("sys_params"); ok {
		sysParams := v.(map[string]interface{})
		for k, v := range sysParams {
			p := &api.Property{
				Name:  fmt.Sprintf("system.%s", k),
				Value: v.(string),
			}
			props.Add(p)
		}
	}

	if v, ok := d.GetOk("config_params"); ok {
		configParams := v.(map[string]interface{})
		for k, v := range configParams {
			p := &api.Property{
				Name:  k,
				Value: v.(string),
			}
			props.Add(p)
		}
	}

	created, err := client.BuildTypes.Create(projectID, bt)

	if err != nil {
		return err
	}

	d.MarkNewResource()
	d.SetId(created.ID)
	d.Partial(true)

	err = client.BuildTypeParameterService(created.ID).Add(props.Items...)

	if err != nil {
		return err
	}

	d.SetPartial("env_params")
	d.SetPartial("config_params")
	d.SetPartial("sys_params")

	if v, ok := d.GetOk("vcs_root"); ok {
		vcs := v.(*schema.Set).List()
		for _, raw := range vcs {
			toAttach := buildVcsRootEntry(raw)

			err = client.BuildTypes.AttachVcsRootEntry(created.ID, toAttach)

			if err != nil {
				return err
			}
		}
		d.SetPartial("vcs_root")
	}

	if v, ok := d.GetOk("step"); ok {
		steps := v.(*schema.Set).List()
		for _, raw := range steps {

			newStep, err := buildStep(raw)
			if err != nil {
				return err
			}

			_, err = client.BuildTypes.AddStep(created.ID, newStep)
			if err != nil {
				return err
			}
		}
		d.SetPartial("step")
	}

	d.Partial(false)

	return resourceBuildConfigurationRead(d, meta)
}

func resourceBuildConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	dt, err := getBuildConfiguration(client, d.Id())
	if err != nil {
		return err
	}

	if err := d.Set("name", dt.Name); err != nil {
		return err
	}

	if err := d.Set("description", dt.Description); err != nil {
		return err
	}

	if err := d.Set("project_id", dt.ProjectID); err != nil {
		return err
	}

	params := dt.Parameters
	var envParams, configParams, sysParams = make(map[string]string), make(map[string]string), make(map[string]string)
	if params != nil && params.Count > 0 {
		paramMap := params.Map()
		for k, v := range paramMap {
			switch {
			case strings.HasPrefix(k, "env."):
				envParams[k[4:len(k)]] = v // Strip env. when setting back the key
			case strings.HasPrefix(k, "system."):
				sysParams[k[7:len(k)]] = v // Strip system. when setting back the key
			default:
				configParams[k] = v
			}
		}
	}

	if err := d.Set("env_params", envParams); err != nil {
		return err
	}

	if err := d.Set("sys_params", sysParams); err != nil {
		return err
	}

	if err := d.Set("config_params", configParams); err != nil {
		return err
	}

	vcsRoots := dt.VcsRootEntries

	if vcsRoots != nil && vcsRoots.Count > 0 {
		var vcsToSave []map[string]interface{}
		for _, el := range vcsRoots.Items {
			m := make(map[string]interface{})
			m["id"] = el.ID
			m["checkout_rules"] = strings.Split(el.CheckoutRules, "\\n")
			vcsToSave = append(vcsToSave, m)
		}

		if err := d.Set("vcs_root", vcsToSave); err != nil {
			return err
		}
	}

	steps, err := client.BuildTypes.GetSteps(d.Id())
	if err != nil {
		return err
	}
	if steps != nil && len(steps) > 0 {
		var stepsToSave []map[string]interface{}
		for _, el := range steps {
			l, err := resourceBuildStepPermGather(el)
			if err != nil {
				return err
			}
			stepsToSave = append(stepsToSave, l)
		}

		if err := d.Set("step", stepsToSave); err != nil {
			return err
		}
	}

	return nil
}

var stepTypeMap = map[string]string{
	"jetbrains_powershell": "powershell",
}

func resourceBuildStepPermGather(s api.Step) (map[string]interface{}, error) {
	mapType := stepTypeMap[s.Type()]

	switch mapType {
	case "powershell":
		return resourceBuildStepPowershellGather(s.(*api.StepPowershell)), nil
	default:
		return nil, fmt.Errorf("Build step type '%s' not supported", s.Type())
	}
}

func resourceBuildStepPowershellGather(s *api.StepPowershell) map[string]interface{} {
	m := make(map[string]interface{})
	if s.ScriptFile != "" {
		m["file"] = s.ScriptFile
		if s.ScriptArgs != "" {
			m["args"] = s.ScriptArgs
		}
	}
	if s.Name() != "" {
		m["name"] = s.Name()
	}
	m["type"] = "powershell"

	return m
}

func resourceBuildConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBuildConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	return client.BuildTypes.Delete(d.Id())
}

func getBuildConfiguration(c *api.Client, id string) (*api.BuildType, error) {
	dt, err := c.BuildTypes.GetByID(id)
	if err != nil {
		return nil, err
	}

	return dt, nil
}

func buildStep(raw interface{}) (api.Step, error) {
	localStep := raw.(map[string]interface{})

	t := localStep["type"].(string)
	switch t {
	case "powershell":
		return expandStepPowershell(localStep)
	default:
		return nil, fmt.Errorf("Unsupported step type '%s'", t)
	}
}

func expandStepPowershell(dt map[string]interface{}) (*api.StepPowershell, error) {
	var file, args, name string

	if v, ok := dt["file"]; ok {
		file = v.(string)
	}
	if v, ok := dt["args"]; ok {
		args = v.(string)
	}
	if v, ok := dt["name"]; ok {
		name = v.(string)
	}

	return api.NewStepPowershellScriptFile(name, file, args)
}

func buildVcsRootEntry(raw interface{}) *api.VcsRootEntry {
	localVcs := raw.(map[string]interface{})
	rawRules := localVcs["checkout_rules"].([]interface{})
	var toAttachRules string
	if len(rawRules) > 0 {
		stringRules := make([]string, len(rawRules))
		for i, el := range rawRules {
			stringRules[i] = el.(string)
		}
		toAttachRules = strings.Join(stringRules, "\\n")
	}

	return api.NewVcsRootEntryWithRules(&api.VcsRootReference{ID: localVcs["id"].(string)}, toAttachRules)
}

func vcsRootHash(v interface{}) int {
	raw := v.(map[string]interface{})
	return schema.HashString(raw["id"].(string))
}

func stepSetHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["type"].(string)))

	if v, ok := m["name"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m["file"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m["args"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	return hashcode.String(buf.String())
}

func validateStepType() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"powershell",
	}, true)
}

package teamcity

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"

	api "github.com/cvbarros/go-teamcity/teamcity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceBuildConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceBuildConfigCreate,
		Read:   resourceBuildConfigRead,
		Update: resourceBuildConfigUpdate,
		Delete: resourceBuildConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			if diff.HasChange("settings") {
				o, n := diff.GetChange("settings")

				os := o.(*schema.Set)
				ns := n.(*schema.Set)
				if os.Len() == 0 || ns.Len() == 0 {
					return nil
				}
				osi, err := expandBuildConfigOptionsRaw(os)
				if err != nil {
					return err
				}
				nsi, err := expandBuildConfigOptionsRaw(ns)
				if err != nil {
					return err
				}

				if buildCounterChange(osi, nsi) {
					var setComputed bool

					// If the configuration doesn't specify the build counter, set the value from READ and mark settings as computed
					if nsi.BuildCounter == 0 {
						log.Printf("[INFO] Build counter not defined in config. Setting it to computed: %v after reading.", osi.BuildCounter)
						nsi.BuildCounter = osi.BuildCounter
						setComputed = true
					} else if osi.BuildCounter > nsi.BuildCounter {
						log.Printf("[INFO] Build counter computed is higher, adjusting state. Old: %v, New: %v.", osi.BuildCounter, nsi.BuildCounter)
						nsi.BuildCounter = osi.BuildCounter
						setComputed = true
					}
					if setComputed {
						computed := flattenBuildConfigOptionsRaw(nsi)
						_ = diff.SetNew("settings", []map[string]interface{}{computed})
					}
				}
			}
			return nil
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"is_template": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"step_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"powershell", "cmd_line"}, false),
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
						"code": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"condition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:         schema.TypeString,
										ValidateFunc: validation.StringInSlice(api.OperatorStrings, false),
										Required:     true,
									},
									"parameter_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"env_params": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"config_params": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"sys_params": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"settings": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				MaxItems:   1,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"configuration_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"REGULAR", "DEPLOYMENT", "COMPOSITE"}, false),
							Default:      "REGULAR",
						},
						"build_number_format": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "%build.counter%",
						},
						"build_counter": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(0),
							Computed:     true,
						},
						"allow_personal_builds": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"artifact_paths": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"detect_hanging": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"status_widget": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"concurrent_limit": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(0),
							Default:      0,
						},
					},
				},
			},
			"templates": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceBuildConfigInstanceResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceBuildConfigInstanceStateUpgradeV0,
				Version: 0,
			},
		},
	}
}

func buildCounterChange(o *api.BuildTypeOptions, n *api.BuildTypeOptions) bool {
	return o.AllowPersonalBuildTriggering == n.AllowPersonalBuildTriggering &&
		reflect.DeepEqual(o.ArtifactRules, n.ArtifactRules) &&
		o.BuildConfigurationType == n.BuildConfigurationType &&
		o.BuildNumberFormat == n.BuildNumberFormat &&
		o.EnableHangingBuildsDetection == n.EnableHangingBuildsDetection &&
		o.EnableStatusWidget == n.EnableStatusWidget &&
		o.MaxSimultaneousBuilds == n.MaxSimultaneousBuilds &&
		o.BuildCounter != n.BuildCounter
}

func resourceBuildConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	var projectID, name string
	isTemplate := false

	if err := validateBuildConfig(d); err != nil {
		return err
	}

	if v, ok := d.GetOk("project_id"); ok {
		projectID = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	log.Printf("[DEBUG] resourceBuildConfigCreate: starting create for build configuration named '%v'.", name)

	var bt *api.BuildType
	var err error

	if v, ok := d.GetOk("is_template"); ok {
		log.Printf("[DEBUG] resourceBuildConfigCreate: setting is_template = '%v'.", v.(bool))
		isTemplate = v.(bool)
	}

	if isTemplate {
		bt, err = api.NewBuildTypeTemplate(projectID, name)
	} else {
		bt, err = api.NewBuildType(projectID, name)
	}
	if err != nil {
		return err
	}

	//BuildType templates don't support description
	if v, ok := d.GetOk("description"); ok && !isTemplate {
		bt.Description = v.(string)
	}

	bt.Parameters, err = expandParameterCollection(d)
	if err != nil {
		return err
	}

	opt, err := expandBuildConfigOptions(d)
	if err != nil {
		return err
	}
	if opt != nil {
		bt.Options = opt
		opt.Template = isTemplate
	}

	created, err := client.BuildTypes.Create(projectID, bt)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] resourceBuildConfigCreate: sucessfully created build configuration with id = '%v'. Marking new resource.", created.ID)

	d.MarkNewResource()
	d.SetId(created.ID)

	log.Printf("[DEBUG] resourceBuildConfigCreate: initial creation finished. Calling resourceBuildConfigUpdate to update the rest of resource.")

	return resourceBuildConfigUpdate(d, meta)
}

func resourceBuildConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	dt, err := getBuildConfiguration(client, d.Id())
	log.Printf("[DEBUG] resourceBuildConfigUpdate started for resouceId: %v", d.Id())

	if err != nil {
		return err
	}

	if d.HasChange("is_template") {
		err := validateBuildConfig(d)
		if err != nil {
			return err
		}
	}

	var changed bool
	if d.HasChange("sys_params") || d.HasChange("config_params") || d.HasChange("env_params") {
		log.Printf("[DEBUG] resourceBuildConfigUpdate: change detected for params")
		dt.Parameters, err = expandParameterCollection(d)
		if err != nil {
			return err
		}
		changed = true
	}
	if v, ok := d.GetOk("description"); ok {
		if d.HasChange("description") {
			log.Printf("[DEBUG] resourceBuildConfigUpdate: change detected for description")
			dt.Description = v.(string)
			changed = true
		}
	}
	if d.HasChange("settings") {
		isTemplate := false
		if v, ok := d.GetOk("is_template"); ok {
			log.Printf("[DEBUG] resourceBuildConfigCreate: setting is_template = '%v'.", v.(bool))
			isTemplate = v.(bool)
		}

		log.Printf("[DEBUG] resourceBuildConfigUpdate: change detected for settings")
		if _, ok := d.GetOk("settings"); ok {
			dt.Options, err = expandBuildConfigOptions(d)
			dt.Options.Template = isTemplate
			changed = true
		}
	}

	if changed {
		_, err := client.BuildTypes.Update(dt)
		if err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("vcs_root"); ok {
		vcs := v.(*schema.Set).List()
		for _, raw := range vcs {
			toAttach := buildVcsRootEntry(raw)

			err := client.BuildTypes.AttachVcsRootEntry(dt.ID, toAttach)
			if err != nil {
				return err
			}
			log.Printf("[DEBUG] resourceBuildConfigUpdate: attached vcsRoot '%v' to build configuration", toAttach.ID)
		}
	}

	if d.HasChange("step") {
		log.Printf("[DEBUG] resourceBuildConfigUpdate: change detected for steps")
		o, n := d.GetChange("step")
		os := o.([]interface{})
		ns := n.([]interface{})

		remove, _ := expandBuildSteps(os)
		add, err := expandBuildSteps(ns)

		if err != nil {
			return err
		}
		if len(remove) > 0 {
			for _, s := range remove {
				err := client.BuildTypes.DeleteStep(dt.ID, s.GetID())
				if err != nil {
					return err
				}
			}
		}
		if len(add) > 0 {
			for _, s := range add {
				_, err := client.BuildTypes.AddStep(dt.ID, s)
				if err != nil {
					return err
				}
			}
		}
	}

	if d.HasChange("templates") {
		log.Printf("[DEBUG] resourceBuildConfigUpdate: change detected for templates")
		remove, add := getChangeExpandedStringList(d.GetChange("templates"))
		buildTemplateService := client.BuildTemplateService(d.Id())
		for _, a := range add {
			_, err := buildTemplateService.Attach(a)
			log.Printf("[DEBUG] resourceBuildConfigUpdate: attached template '%v' to build configuration", a)
			if err != nil {
				return err
			}
		}
		for _, r := range remove {
			err := buildTemplateService.Detach(r)
			if err != nil {
				return err
			}
			log.Printf("[DEBUG] resourceBuildConfigUpdate: detached template '%v' from build configuration", r)
		}
	}

	log.Printf("[DEBUG] resourceBuildConfigUpdate: updated finished. Calling 'read' to refresh state.")
	return resourceBuildConfigRead(d, meta)
}

func resourceBuildConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	log.Printf("[DEBUG] resourceBuildConfigDelete: destroying build configuration '%v'.", d.Id())
	return client.BuildTypes.Delete(d.Id())
}

func resourceBuildConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	log.Printf("[DEBUG] resourceBuildConfigRead started for resouceId: %v", d.Id())
	dt, err := getBuildConfiguration(client, d.Id())
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] BuildConfiguration '%v' retrieved successfully", dt.Name)
	if err := d.Set("name", dt.Name); err != nil {
		return err
	}
	if err := d.Set("is_template", dt.IsTemplate); err != nil {
		return err
	}
	//description not supported for templates.
	if !dt.IsTemplate {
		if err := d.Set("description", dt.Description); err != nil {
			return err
		}
	}
	if err := d.Set("project_id", dt.ProjectID); err != nil {
		return err
	}
	err = flattenParameterCollection(d, dt.Parameters)
	if err != nil {
		return err
	}
	err = flattenBuildConfigOptions(d, dt.Options)
	if err != nil {
		return err
	}
	err = flattenTemplates(d, dt.Templates)
	if err != nil {
		return err
	}

	vcsRoots := dt.VcsRootEntries

	if vcsRoots != nil && len(vcsRoots) > 0 {
		var vcsToSave []map[string]interface{}
		for _, el := range vcsRoots {
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
			l, err := flattenBuildStep(el)
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

func validateBuildConfig(d *schema.ResourceData) error {
	if v, ok := d.GetOk("is_template"); ok {
		isTemplate := v.(bool)

		if isTemplate {
			if _, ok := d.GetOk("description"); ok {
				return fmt.Errorf("'description' field is not supported for Build Configuration Templates. See issue https://youtrack.jetbrains.com/issue/TW-63617 for details")
			}
			if _, ok := d.GetOk("settings"); ok {
				opt, err := expandBuildConfigOptions(d)
				if err != nil {
					return err
				}
				// If there's build counter specified in the configuration
				if opt.BuildCounter != 0 {
					return errors.New("'settings.build_counter' field is not supported for Build Configuration Templates, is_template = true")
				}
			}
		}
	}
	return nil
}

func getBuildConfiguration(c *api.Client, id string) (*api.BuildType, error) {
	dt, err := c.BuildTypes.GetByID(id)
	if err != nil {
		return nil, err
	}

	return dt, nil
}

var stepTypeMap = map[string]string{
	api.StepTypePowershell:  "powershell",
	api.StepTypeCommandLine: "cmd_line",
}

func flattenTemplates(d *schema.ResourceData, templates *api.Templates) error {
	if templates == nil {
		return nil
	}
	templateIds := make([]string, templates.Count)
	if templates.Count > 0 {
		for i, v := range templates.Items {
			templateIds[i] = v.ID
		}
	}
	if err := d.Set("templates", templateIds); err != nil {
		return err
	}
	return nil
}

func flattenParameterCollection(d *schema.ResourceData, params *api.Parameters) error {
	var configParams, sysParams, envParams = flattenParameters(params)

	if len(envParams) > 0 {
		if err := d.Set("env_params", envParams); err != nil {
			return err
		}
	}
	if len(sysParams) > 0 {
		if err := d.Set("sys_params", sysParams); err != nil {
			return err
		}
	}
	if len(configParams) > 0 {
		if err := d.Set("config_params", configParams); err != nil {
			return err
		}
	}
	return nil
}

func expandParameterCollection(d *schema.ResourceData) (*api.Parameters, error) {
	var config, system, env *api.Parameters
	if v, ok := d.GetOk("env_params"); ok {
		p, err := expandParameters(v.(map[string]interface{}), api.ParameterTypes.EnvironmentVariable)
		if err != nil {
			return nil, err
		}
		env = p
	}

	if v, ok := d.GetOk("sys_params"); ok {
		p, err := expandParameters(v.(map[string]interface{}), api.ParameterTypes.System)
		if err != nil {
			return nil, err
		}
		system = p
	}

	if v, ok := d.GetOk("config_params"); ok {
		p, err := expandParameters(v.(map[string]interface{}), api.ParameterTypes.Configuration)
		if err != nil {
			return nil, err
		}
		config = p
	}

	out := api.NewParametersEmpty()

	if config != nil {
		out = out.Concat(config)
	}
	if system != nil {
		out = out.Concat(system)
	}
	if env != nil {
		out = out.Concat(env)
	}
	return out, nil
}

func flattenParameters(dt *api.Parameters) (config map[string]string, sys map[string]string, env map[string]string) {
	env, sys, config = make(map[string]string), make(map[string]string), make(map[string]string)
	for _, p := range dt.Items {
		switch p.Type {
		case api.ParameterTypes.Configuration:
			config[p.Name] = p.Value
		case api.ParameterTypes.EnvironmentVariable:
			env[p.Name] = p.Value
		case api.ParameterTypes.System:
			sys[p.Name] = p.Value
		}
	}
	return config, sys, env
}

func expandParameters(raw map[string]interface{}, paramType string) (*api.Parameters, error) {
	out := api.NewParametersEmpty()
	for k, v := range raw {
		p, err := api.NewParameter(paramType, k, v.(string))
		if err != nil {
			return nil, err
		}
		out.AddOrReplaceParameter(p)
	}
	return out, nil
}

func expandBuildConfigOptions(d *schema.ResourceData) (*api.BuildTypeOptions, error) {
	v, ok := d.GetOk("settings")
	if !ok {
		return nil, nil
	}

	return expandBuildConfigOptionsRaw(v.(*schema.Set))
}

func expandBuildConfigOptionsRaw(v *schema.Set) (*api.BuildTypeOptions, error) {
	raw := v.List()[0].(map[string]interface{})
	opt := api.NewBuildTypeOptionsWithDefaults()

	if v, ok := raw["configuration_type"]; ok {
		opt.BuildConfigurationType = strings.ToUpper(v.(string))
	}
	if v, ok := raw["build_number_format"]; ok {
		opt.BuildNumberFormat = v.(string)
	}
	if v, ok := raw["build_counter"]; ok {
		opt.BuildCounter = v.(int)
	}
	if v, ok := raw["allow_personal_builds"]; ok {
		opt.AllowPersonalBuildTriggering = v.(bool)
	}
	if v, ok := raw["artifact_paths"]; ok {
		opt.ArtifactRules = expandStringSlice(v.([]interface{}))
	}
	if v, ok := raw["detect_hanging"]; ok {
		opt.EnableHangingBuildsDetection = v.(bool)
	}
	if v, ok := raw["status_widget"]; ok {
		opt.EnableStatusWidget = v.(bool)
	}
	if v, ok := raw["concurrent_limit"]; ok {
		opt.MaxSimultaneousBuilds = v.(int)
	}

	return opt, nil
}
func flattenBuildConfigOptions(d *schema.ResourceData, dt *api.BuildTypeOptions) error {
	m := flattenBuildConfigOptionsRaw(dt)
	return d.Set("settings", []map[string]interface{}{m})
}

func flattenBuildConfigOptionsRaw(dt *api.BuildTypeOptions) map[string]interface{} {
	m := make(map[string]interface{})

	m["configuration_type"] = dt.BuildConfigurationType
	m["build_number_format"] = dt.BuildNumberFormat
	m["build_counter"] = dt.BuildCounter
	m["allow_personal_builds"] = dt.AllowPersonalBuildTriggering
	m["artifact_paths"] = flattenStringSlice(dt.ArtifactRules)
	m["detect_hanging"] = dt.EnableHangingBuildsDetection
	m["status_widget"] = dt.EnableStatusWidget
	m["concurrent_limit"] = dt.MaxSimultaneousBuilds

	return m
}

func flattenBuildStep(s api.Step) (map[string]interface{}, error) {
	mapType := stepTypeMap[s.Type()]
	var out map[string]interface{}
	var err error
	switch mapType {
	case "powershell":
		out, err = flattenBuildStepPowershell(s.(*api.StepPowershell)), nil
	case "cmd_line":
		out, err = flattenBuildStepCmdLine(s.(*api.StepCommandLine)), nil
	default:
		return nil, fmt.Errorf("build step type '%s' not supported", s.Type())
	}
	out["step_id"] = s.GetID()
	return out, err
}

func flattenBuildStepPowershell(s *api.StepPowershell) map[string]interface{} {
	m := make(map[string]interface{})
	if s.ScriptFile != "" {
		m["file"] = s.ScriptFile
		if s.ScriptArgs != "" {
			m["args"] = s.ScriptArgs
		}
	}
	if s.Code != "" {
		m["code"] = s.Code
	}
	if s.Name != "" {
		m["name"] = s.Name
	}
	if s.ExecuteMode != "" {
		m["mode"] = s.ExecuteMode
	}

	if s.Conditions != "" {
		re := regexp.MustCompile(`\["([^"]+)","([^"]+)"."([^"]*)"]`)
		cn := re.FindAllStringSubmatch(s.Conditions, -1)
		var conditions []interface{}
		for _, expr := range cn {
			cs := map[string]string{
				"operator":       expr[1],
				"parameter_name": expr[2],
				"value":          expr[3],
			}
			conditions = append(conditions, cs)
		}
		m["condition"] = conditions
	}
	m["type"] = "powershell"

	return m
}

func flattenBuildStepCmdLine(s *api.StepCommandLine) map[string]interface{} {
	m := make(map[string]interface{})
	if s.CommandExecutable != "" {
		m["file"] = s.CommandExecutable
		if s.CommandParameters != "" {
			m["args"] = s.CommandParameters
		}
	}
	if s.CustomScript != "" {
		m["code"] = s.CustomScript
	}
	if s.Name != "" {
		m["name"] = s.Name
	}
	if s.ExecuteMode != "" {
		m["mode"] = s.ExecuteMode
	}

	if s.Conditions != "" {
		re := regexp.MustCompile(`\["([^"]+)","([^"]+)"."([^"]*)"]`)
		cn := re.FindAllStringSubmatch(s.Conditions, -1)
		var conditions []interface{}
		for _, expr := range cn {
			cs := map[string]string{
				"operator":       expr[1],
				"parameter_name": expr[2],
				"value":          expr[3],
			}
			conditions = append(conditions, cs)
		}
		m["condition"] = conditions
	}
	m["type"] = "cmd_line"

	return m
}

func expandBuildSteps(list interface{}) ([]api.Step, error) {
	out := make([]api.Step, 0)
	in := list.([]interface{})
	for _, i := range in {
		s, err := expandBuildStep(i)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, nil
}

func expandBuildStep(raw interface{}) (api.Step, error) {
	localStep := raw.(map[string]interface{})

	t := localStep["type"].(string)
	switch t {
	case "powershell":
		return expandStepPowershell(localStep)
	case "cmd_line":
		return expandStepCmdLine(localStep)
	default:
		return nil, fmt.Errorf("unsupported step type '%s'", t)
	}
}

func expandStepCmdLine(dt map[string]interface{}) (*api.StepCommandLine, error) {
	var file, args, name, code, mode, serializedConditions string

	if v, ok := dt["file"]; ok {
		file = v.(string)
	}
	if v, ok := dt["args"]; ok {
		args = v.(string)
	}
	if v, ok := dt["name"]; ok {
		name = v.(string)
	}
	if v, ok := dt["code"]; ok {
		code = v.(string)
	}
	if v, ok := dt["mode"]; ok {
		mode = v.(string)
	}
	if v, ok := dt["condition"]; ok {
		conditions := v.([]interface{})

		serializedConditions += `[`
		for i, cs := range conditions {
			var operator, parameterName, value string

			condition := cs.(map[string]interface{})

			if v, ok := condition["operator"]; ok {
				operator = v.(string)
			}
			if v, ok := condition["parameter_name"]; ok {
				parameterName = v.(string)
			}
			if v, ok := condition["value"]; ok {
				value = v.(string)
			}

			serializedConditions += `["` + operator + `","` + parameterName + `","` + value + `"]`
			if i != len(conditions)-1 {
				serializedConditions += `,`
			}
		}
		serializedConditions += `]`
	}

	var s *api.StepCommandLine
	var err error
	if file != "" {
		s, err = api.NewStepCommandLineExecutable(name, file, args, mode, serializedConditions)
	} else {
		s, err = api.NewStepCommandLineScript(name, code, mode, serializedConditions)
	}
	if err != nil {
		return nil, err
	}

	if v, ok := dt["step_id"]; ok {
		s.ID = v.(string)
	}
	return s, nil
}

func expandStepPowershell(dt map[string]interface{}) (*api.StepPowershell, error) {
	var file, args, name, code, mode, serializedConditions string

	if v, ok := dt["file"]; ok {
		file = v.(string)
	}
	if v, ok := dt["args"]; ok {
		args = v.(string)
	}
	if v, ok := dt["name"]; ok {
		name = v.(string)
	}
	if v, ok := dt["code"]; ok {
		code = v.(string)
	}
	if v, ok := dt["mode"]; ok {
		mode = v.(string)
	}
	if v, ok := dt["condition"]; ok {
		conditions := v.([]interface{})

		serializedConditions += `[`
		for i, cs := range conditions {
			var operator, parameterName, value string
			condition := cs.(map[string]interface{})

			if v, ok := condition["operator"]; ok {
				operator = v.(string)
			}
			if v, ok := condition["parameter_name"]; ok {
				parameterName = v.(string)
			}
			if v, ok := condition["value"]; ok {
				value = v.(string)
			}

			serializedConditions += `["` + operator + `","` + parameterName + `","` + value + `"]`
			if i != len(conditions)-1 {
				serializedConditions += `,`
			}
		}
		serializedConditions += `]`
	}

	var s *api.StepPowershell
	var err error
	if file != "" {
		s, err = api.NewStepPowershellScriptFile(name, file, args, mode, serializedConditions)
	} else {
		s, err = api.NewStepPowershellCode(name, code, mode, serializedConditions)
	}
	if err != nil {
		return nil, err
	}

	if v, ok := dt["step_id"]; ok {
		s.ID = v.(string)
	}
	return s, nil
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

	if v, ok := m["code"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m["mode"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m["condition"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	return hashcode.String(buf.String())
}

func resourceBuildConfigInstanceResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"is_template": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
						"step_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"powershell", "cmd_line"}, false),
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
						"code": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"condition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:         schema.TypeString,
										ValidateFunc: validation.StringInSlice(api.OperatorStrings, false),
										Required:     true,
									},
									"parameter_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
				Set: stepSetHash,
			},
			"env_params": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"config_params": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"sys_params": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"settings": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				MaxItems:   1,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"configuration_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"REGULAR", "DEPLOYMENT", "COMPOSITE"}, false),
							Default:      "REGULAR",
						},
						"build_number_format": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "%build.counter%",
						},
						"build_counter": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(0),
							Computed:     true,
						},
						"allow_personal_builds": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"artifact_paths": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"detect_hanging": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"status_widget": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"concurrent_limit": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(0),
							Default:      0,
						},
					},
				},
			},
			"templates": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

func resourceBuildConfigInstanceStateUpgradeV0(rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	if raw, ok := rawState["steps"]; ok {
		s := raw.(*schema.Set)
		list := s.List()
		rawState["steps"] = list
	}

	return rawState, nil
}

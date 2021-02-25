package teamcity

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	api "github.com/yext/go-teamcity/teamcity"
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

				os := o.([]interface{})
				ns := n.([]interface{})
				if len(os) == 0 || len(ns) == 0 {
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
						diff.SetNew("settings", []map[string]interface{}{computed})
					}
				}
			}

			return nil
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
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
							Default:  "default",
							ValidateFunc: validation.StringInSlice([]string{
								api.StepExecuteModeDefault,
								api.StepExecuteModeOnlyIfBuildIsSuccessful,
								api.StepExecuteModeEvenWhenFailed,
								api.StepExecuteAlways,
							}, false),
						},
					},
				},
			},
			"feature": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"properties": {
							Type:     schema.TypeMap,
							Required: true,
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
			"env_params_specs": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"config_params_specs": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"sys_params_specs": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"settings": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Computed: true,
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

	if v, ok := d.GetOk("project_id"); ok {
		projectID = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	bt, err := api.NewBuildType(projectID, name)
	if err != nil {
		return err
	}

	if v, ok := d.GetOk("description"); ok {
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
	}

	created, err := client.BuildTypes.Create(projectID, bt)

	if err != nil {
		return err
	}

	d.MarkNewResource()
	d.SetId(created.ID)
	d.Partial(true)

	return resourceBuildConfigUpdate(d, meta)
}

func resourceBuildConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	dt, err := getBuildConfiguration(client, d.Id())
	if err != nil {
		return err
	}

	if d.HasChange("name") {
		v := d.Get("name")
		err = client.BuildTypes.Rename(d.Id(), v.(string))
		if err != nil {
			return err
		}
	}

	var changed bool
	if d.HasChange("config_params") || d.HasChange("env_params") || d.HasChange("sys_params") ||
		d.HasChange("config_params_specs") || d.HasChange("env_params_specs") || d.HasChange("sys_params_specs") {
		dt.Parameters, err = expandParameterCollection(d)
		if err != nil {
			return err
		}
		changed = true
	}
	if v, ok := d.GetOk("description"); ok {
		if d.HasChange("description") {
			dt.Description = v.(string)
			changed = true
		}
	}
	if d.HasChange("settings") {
		if _, ok := d.GetOk("settings"); ok {
			dt.Options, err = expandBuildConfigOptions(d)
			changed = true
		}
	}

	if changed {
		_, err := client.BuildTypes.Update(dt)
		d.SetPartial("settings")
		d.SetPartial("description")
		d.SetPartial("config_params")
		d.SetPartial("env_params")
		d.SetPartial("sys_params")
		d.SetPartial("config_params_specs")
		d.SetPartial("env_params_specs")
		d.SetPartial("sys_params_specs")
		if err != nil {
			return err
		}
	}

	if d.HasChange("vcs_root") {
		oldVCSEntries, err := client.BuildTypes.GetVcsRootEntries(dt.ID)

		if err != nil {
			return err
		}

		for _, old := range oldVCSEntries.Items {
			err := client.BuildTypes.DetachVcsRootEntry(dt.ID, old.ID)

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
			}
		}

		d.SetPartial("vcs_root")
	}

	if d.HasChange("step") {
		add, err := expandBuildSteps(d.Get("step").([]interface{}))
		if err != nil {
			return err
		}
		//Remove all existing steps
		remove, err := client.BuildTypes.GetSteps(d.Id())
		for _, s := range remove {
			err := client.BuildTypes.DeleteStep(dt.ID, s.GetID())
			if err != nil {
				return err
			}
		}
		if len(add) > 0 {
			for i, s := range add {
				added, err := client.BuildTypes.AddStep(dt.ID, s)
				log.Printf("[INFO] Adding step '%v' (%v)with order = %v", s.GetName(), added.GetID(), i+1)
				if err != nil {
					return err
				}
			}
		}
		d.SetPartial("step")
	}

	if d.HasChange("feature") {
		srv := client.BuildFeatureService(d.Id())
		err := srv.DeleteAll()
		if err != nil {
			return err
		}
		add, err := expandBuildFeatures(d.Get("feature").([]interface{}))
		if err != nil {
			return err
		}
		if len(add) > 0 {
			for i, s := range add {
				_, err := srv.Create(s)
				log.Printf("[INFO] Adding build feature '%v' with order = %v", s.Type(), i+1)
				if err != nil {
					return err
				}
			}
		}
		d.SetPartial("feature")
	}

	d.Partial(false)

	return resourceBuildConfigRead(d, meta)
}

func resourceBuildConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	return client.BuildTypes.Delete(d.Id())
}

func resourceBuildConfigRead(d *schema.ResourceData, meta interface{}) error {
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
	err = flattenParameterCollection(d, dt.Parameters)
	if err != nil {
		return err
	}
	err = flattenBuildConfigOptions(d, dt.Options)
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

	srv := client.BuildFeatureService(d.Id())
	buildFeatures, err := srv.GetBuildFeatures()
	buildFeaturesToSave, err := flattenBuildFeatures(buildFeatures)
	if err := d.Set("feature", buildFeaturesToSave); err != nil {
		return err
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

func flattenParameterCollection(d *schema.ResourceData, params *api.Parameters) error {
	var configParams, envParams, sysParams, configSpecs, envSpecs, sysSpecs = flattenParameters(params)

	if len(configParams) > 0 {
		if err := d.Set("config_params", configParams); err != nil {
			return err
		}
	}
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
	if len(configSpecs) > 0 {
		if err := d.Set("config_params_specs", configSpecs); err != nil {
			return err
		}
	}
	if len(envSpecs) > 0 {
		if err := d.Set("env_params_specs", envSpecs); err != nil {
			return err
		}
	}
	if len(sysSpecs) > 0 {
		if err := d.Set("sys_params_specs", sysSpecs); err != nil {
			return err
		}
	}
	return nil
}

func expandParameterCollection(d *schema.ResourceData) (*api.Parameters, error) {
	var config, system, env *api.Parameters
	if v, ok := d.GetOk("env_params"); ok {
		specs := d.Get("env_params_specs")
		p, err := expandParameters(v.(map[string]interface{}), specs, api.ParameterTypes.EnvironmentVariable)
		if err != nil {
			return nil, err
		}
		env = p
	}

	if v, ok := d.GetOk("sys_params"); ok {
		specs := d.Get("sys_params_specs")
		p, err := expandParameters(v.(map[string]interface{}), specs, api.ParameterTypes.System)
		if err != nil {
			return nil, err
		}
		system = p
	}

	if v, ok := d.GetOk("config_params"); ok {
		specs := d.Get("config_params_specs")
		p, err := expandParameters(v.(map[string]interface{}), specs, api.ParameterTypes.Configuration)
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

func flattenParameters(dt *api.Parameters) (config, env, sys, configSpecs, envSpecs, sysSpecs map[string]string) {
	config, env, sys, configSpecs, envSpecs, sysSpecs = make(map[string]string), make(map[string]string),
		make(map[string]string), make(map[string]string), make(map[string]string), make(map[string]string)
	for _, p := range dt.Items {
		switch p.Type {
		case api.ParameterTypes.Configuration:
			config[p.Name] = p.Value
			if p.Specs != nil {
				configSpecs[p.Name] = p.Specs.RawValue
			}
		case api.ParameterTypes.EnvironmentVariable:
			env[p.Name] = p.Value
			if p.Specs != nil {
				envSpecs[p.Name] = p.Specs.RawValue
			}
		case api.ParameterTypes.System:
			sys[p.Name] = p.Value
			if p.Specs != nil {
				sysSpecs[p.Name] = p.Specs.RawValue
			}
		}
	}
	return config, env, sys, configSpecs, envSpecs, sysSpecs
}

func expandParameters(raw map[string]interface{}, specs interface{}, paramType string) (*api.Parameters, error) {
	var specMap map[string]interface{}
	if specs != nil {
		specMap = specs.(map[string]interface{})
	}

	for k := range specMap {
		if _, ok := raw[k]; !ok {
			return nil, fmt.Errorf("invalid spec %v", k)
		}
	}

	out := api.NewParametersEmpty()
	for k, v := range raw {
		p, err := api.NewParameter(paramType, k, v.(string))
		if err != nil {
			return nil, err
		}
		if spec, ok := specMap[k]; ok {
			p.AddParameterSpecs(spec.(string))
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

	return expandBuildConfigOptionsRaw(v.([]interface{}))
}

func expandBuildConfigOptionsRaw(v []interface{}) (*api.BuildTypeOptions, error) {
	raw := v[0].(map[string]interface{})
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
		return nil, fmt.Errorf("Build step type '%s' not supported", s.Type())
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
	if s.ExecuteMode != "" {
		m["mode"] = s.ExecuteMode
	}
	if s.Name != "" {
		m["name"] = s.Name
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
	if s.ExecuteMode != "" {
		m["mode"] = s.ExecuteMode
	}
	if s.Name != "" {
		m["name"] = s.Name
	}
	m["type"] = "cmd_line"

	return m
}

func flattenBuildFeatures(bfs []api.BuildFeature) ([]map[string]interface{}, error) {
	var bfsToSave []map[string]interface{}
	for _, bf := range bfs {
		bfToSave := make(map[string]interface{})
		gbf := bf.(*api.GenericBuildFeature)
		bfToSave["type"] = gbf.Type()

		props := gbf.Properties()
		bfToSave["properties"] = make(map[string]string)
		if props != nil && props.Count > 0 {
			propertyMap := bfToSave["properties"].(map[string]string)
			for _, property := range props.Items {
				propertyMap[property.Name] = property.Value
			}
		}
		bfsToSave = append(bfsToSave, bfToSave)
	}
	return bfsToSave, nil
}

func expandBuildFeatures(list interface{}) ([]api.BuildFeature, error) {
	var out []api.BuildFeature
	rawBuildFeatures := list.([]interface{})
	for _, rawBF := range rawBuildFeatures {
		expandedBF, err := expandBuildFeature(rawBF)
		if err != nil {
			return nil, err
		}
		out = append(out, expandedBF)
	}

	return out, nil
}

func expandBuildFeature(raw interface{}) (api.BuildFeature, error) {
	feature := raw.(map[string]interface{})
	var featureType string
	var properties map[string]interface{}
	if v, _ := feature["type"]; len(v.(string)) > 0 {
		featureType = v.(string)
	} else {
		return nil, errors.New("feature type cannot be empty")
	}
	if v, ok := feature["properties"]; ok {
		properties = v.(map[string]interface{})
	}
	bf, err := api.NewGenericBuildFeature(featureType, properties)
	if err != nil {
		return nil, err
	}
	return bf, nil
}

func expandBuildSteps(list interface{}) ([]api.Step, error) {
	out := make([]api.Step, 0)
	in := list.([]interface{})
	names := make(map[string]struct{})
	for _, i := range in {
		s, err := expandBuildStep(i)
		if err != nil {
			return nil, err
		}

		if _, exist := names[s.GetName()]; exist {
			return nil, fmt.Errorf("Duplicate build step names '%s'", s.GetName())
		} else {
			names[s.GetName()] = struct{}{}
		}

		out = append(out, s)
	}

	return out, nil
}

func expandBuildStep(raw interface{}) (step api.Step, err error) {
	localStep := raw.(map[string]interface{})
	err = nil
	t := localStep["type"].(string)
	switch t {
	case "powershell":
		step, err = expandStepPowershell(localStep)
	case "cmd_line":
		step, err = expandStepCmdLine(localStep)
	default:
		return nil, fmt.Errorf("Unsupported step type '%s'", t)
	}
	return step, err
}

func expandStepCmdLine(dt map[string]interface{}) (*api.StepCommandLine, error) {
	var file, args, name, code string

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

	var s *api.StepCommandLine
	var err error
	if file != "" {
		s, err = api.NewStepCommandLineExecutable(name, file, args)
	} else {
		s, err = api.NewStepCommandLineScript(name, code)
	}
	if err != nil {
		return nil, err
	}

	if v, ok := dt["step_id"]; ok {
		s.ID = v.(string)
	}
	if v, ok := dt["mode"]; ok {
		s.ExecuteMode = v.(string)
	}
	return s, nil
}

func expandStepPowershell(dt map[string]interface{}) (*api.StepPowershell, error) {
	var file, args, name, code string

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

	var s *api.StepPowershell
	var err error
	if file != "" {
		s, err = api.NewStepPowershellScriptFile(name, file, args)
	} else {
		s, err = api.NewStepPowershellCode(name, code)
	}
	if err != nil {
		return nil, err
	}

	if v, ok := dt["step_id"]; ok {
		s.ID = v.(string)
	}
	if v, ok := dt["mode"]; ok {
		s.ExecuteMode = v.(string)
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

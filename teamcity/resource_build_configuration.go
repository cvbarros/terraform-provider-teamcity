package teamcity

import (
	"fmt"
	"strings"

	api "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/hashicorp/terraform/helper/schema"
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
	var projectId string

	props := api.NewProperties()

	if v, ok := d.GetOk("project_id"); ok {
		projectId = v.(string)
		bt.ProjectID = projectId
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

	created, err := client.BuildTypes.Create(projectId, bt)

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

	return nil
}

func resourceBuildConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBuildConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	return client.BuildTypes.Delete(d.Id())
}

func getBuildConfiguration(c *api.Client, id string) (*api.BuildType, error) {
	dt, err := c.BuildTypes.GetById(id)
	if err != nil {
		return nil, err
	}

	return dt, nil
}

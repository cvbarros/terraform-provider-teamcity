package teamcity

import (
	"fmt"
	"log"

	api "github.com/cvbarros/go-teamcity/teamcity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,
		Importer: &schema.ResourceImporter{
			State: resourceProjectImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
		},
	}
}

func resourceProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	name := d.Get("name").(string)
	parentProjectId := d.Get("parent_id").(string)

	newProj, err := api.NewProject(name, "", parentProjectId)
	if err != nil {
		return err
	}

	created, err := client.Projects.Create(newProj)
	if err != nil {
		return err
	}

	d.SetId(created.ID)

	return resourceProjectUpdate(d, client)
}

func resourceProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	dt, err := client.Projects.GetByID(d.Id())
	if err != nil {
		return err
	}

	if v, ok := d.GetOk("description"); ok {
		dt.Description = v.(string)
	}

	if v, ok := d.GetOk("parent_id"); ok {
		if v != "" {
			dt.SetParentProject(v.(string))
		}
	}

	dt.Parameters, err = expandParameterCollection(d)
	if err != nil {
		return err
	}

	_, err = client.Projects.Update(dt)
	if err != nil {
		return nil
	}
	return resourceProjectRead(d, meta)
}

func resourceProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	dt, err := client.Projects.GetByID(d.Id())
	if err != nil {
		// handles this being deleted outside of TF
		if isNotFoundError(err) {
			log.Printf("[DEBUG] Project was not found - removing from state!")
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("name", dt.Name)
	d.Set("description", dt.Description)
	d.Set("parent_id", dt.ParentProjectID)

	return flattenParameterCollection(d, dt.Parameters)
}

func resourceProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	log.Print(fmt.Sprintf("[DEBUG]: resourceProjectDelete - Destroying project %v", d.Id()))
	err := client.Projects.Delete(d.Id())
	log.Print(fmt.Sprintf("[INFO]: resourceProjectDelete - Destroyed project %v", d.Id()))
	return err
}

func resourceProjectImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceProjectRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

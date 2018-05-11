package teamcity

import (
	"log"

	api "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/hashicorp/terraform/helper/schema"
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
			},
		},
	}
}

func resourceProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	newProj := &api.Project{
		Name: d.Get("name").(string),
	}

	created, err := client.Projects.Create(newProj)
	if err != nil {
		return err
	}

	d.MarkNewResource()
	d.SetId(created.ID)
	d.Set("name", created.Name)
	return nil
}

func resourceProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	dt, err := getProject(client, d.Id())
	if err != nil {
		return err
	}

	if err := d.Set("name", dt.Name); err != nil {
		return err
	}

	log.Printf("[DEBUG] Project: %v", dt)
	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceProjectRead(d, meta)
}

func resourceProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	return client.Projects.Delete(d.Id())
	return nil
}

func resourceProjectImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceProjectRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func getProject(c *api.Client, id string) (*api.Project, error) {
	dt, err := c.Projects.GetById(id)
	if err != nil {
		return nil, err
	}

	return dt, nil
}

package teamcity

import (
	"log"

	api "github.com/cvbarros/go-teamcity-sdk/client"
	project "github.com/cvbarros/go-teamcity-sdk/client/project"
	models "github.com/cvbarros/go-teamcity-sdk/models"
	"github.com/hashicorp/terraform/helper/schema"
)

// GetProject Retrieves a project by id from the REST API
func GetProject(c *api.TeamCityREST, id string) (*models.Project, error) {
	getParams := project.NewServeProjectParams()
	getParams.ProjectLocator = projectLocatorByID(id)
	dt, err := c.Project.ServeProject(getParams)
	if err != nil {
		return nil, err
	}

	return dt.Payload, nil
}

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
	client := meta.(*api.TeamCityREST)

	params := project.NewCreateProjectParams()
	params.WithBody(&models.NewProjectDescription{
		Name: d.Get("name").(string),
	})

	resp, err := client.Project.CreateProject(params)
	if err != nil {
		return err
	}

	d.MarkNewResource()
	d.SetId(resp.Payload.ID)
	d.Set("name", resp.Payload.Name)

	return nil
}

func resourceProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.TeamCityREST)

	dt, err := GetProject(client, d.Id())
	if err != nil {
		return err
	}

	if err := d.Set("name", dt.Name); err != nil {
		return err
	}

	log.Printf("[DEBUG] project: %v", dt)

	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.TeamCityREST).Project

	if attr, ok := d.GetOk("name"); ok {
		putParams := project.NewSetProjectFieldParams().
			WithProjectLocator(d.Id()).
			WithField("name").
			WithBody(attr.(string))

		if _, err := client.SetProjectField(putParams); err != nil {
			return err
		}
	}

	return resourceProjectRead(d, meta)
}

func resourceProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.TeamCityREST)
	var id = projectLocatorByID(d.Id())

	deleteParams := project.NewDeleteProjectParams().
		WithProjectLocator(id)

	return client.Project.DeleteProject(deleteParams)
}

func resourceProjectImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceProjectRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func projectLocatorByID(projectID string) string {
	return projectID
}

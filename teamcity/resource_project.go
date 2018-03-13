package teamcity

import (
	client "github.com/cvbarros/go-teamcity-sdk/client"
	project "github.com/cvbarros/go-teamcity-sdk/client/project"
	models "github.com/cvbarros/go-teamcity-sdk/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.TeamCityREST)

	params := project.NewCreateProjectParams()
	params.WithBody(&models.NewProjectDescription{
		Name: d.Get("name").(string),
	})

	resp, err := client.Project.CreateProject(params)
	if err != nil {
		return err
	}

	d.SetId(resp.Payload.ID)
	d.Set("name", resp.Payload.Name)

	return nil
}

func resourceProjectRead(d *schema.ResourceData, meta interface{}) error {
	// client := meta.(*api.REST)

	// id, err := strconv.Atoi(d.Id())
	// if err != nil{
	// 	return err
	// }
	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceProjectDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

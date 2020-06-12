package teamcity

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"

	api "github.com/cvbarros/go-teamcity/teamcity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	ResourceProjectRootId   = "_Root"
	ResourceProjectRootName = "<Root project>"
	ResourceProjectRootDescription = "Contains all other projects"
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
			"root": {
				Type: schema.TypeBool,
				Optional: true,
				Default: false,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
		CustomizeDiff: func(diff *schema.ResourceDiff, i interface{}) error {
			root := diff.Get("root").(bool)
			if root {
				if v, ok := diff.GetOk("name"); ok && v.(string) != ResourceProjectRootName {
					return errors.New("'name' cannot be defined for the root project")
				}
				if v, ok := diff.GetOk("description"); ok && v.(string) != ResourceProjectRootDescription {
					return errors.New("'description' cannot be defined for the root project")
				}
			}
			// Validate required name if root = false
			if _, ok := diff.GetOk("name"); !ok && !root {
				return errors.New("'name' is required for non-root project")
			}

			return nil
		},
	}
}

func resourceProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	var name, parentID string

	if isRoot, ok := d.GetOk("root"); ok && isRoot.(bool) {
		// Skip creation altogether
		return resourceProjectUpdate(d, client)
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	if v, ok := d.GetOk("parent_id"); ok {
		if v != "" {
			parentID = v.(string)
		}
	}

	newProj, err := api.NewProject(name, "", parentID)
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
	if isRoot, ok := d.GetOk("root"); ok && isRoot.(bool) {
		// Skip creation altogether
		d.SetId(ResourceProjectRootId)
	}

	client := meta.(*api.Client)
	dt, err := client.Projects.GetByID(d.Id())
	if err != nil {
		return err
	}

	if v, ok := d.GetOk("description"); ok {
		dt.Description = v.(string)
	}

	if d.HasChange("parent_id") {
		parentId := d.Get("parent_id").(string)
		if parentId == "" {
			parentId = ResourceProjectRootId
		}
		dt.SetParentProject(parentId)
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

	dt, err := getProject(client, d.Id())
	if err != nil {
		return err
	}

	root := d.Id() == ResourceProjectRootId
	if err := d.Set("root", root); err != nil {
		return err
	}
	if err := d.Set("name", dt.Name); err != nil {
		return err
	}
	if err := d.Set("description", dt.Description); err != nil {
		return err
	}
	if dt.ParentProject != nil {
		if err := d.Set("parent_id", dt.ParentProject.ID); err != nil {
			return err
		}
	} else {
		d.Set("parent_id","")
	}

	flattenParameterCollection(d, dt.Parameters)
	return nil
}

func resourceProjectDelete(d *schema.ResourceData, meta interface{}) error {
	if root, ok := d.GetOk("root"); ok && root.(bool) {
		//Skip destruction of _Root project
		d.SetId("")
		return nil
	}
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

func getProject(c *api.Client, id string) (*api.Project, error) {
	dt, err := c.Projects.GetByID(id)
	if err != nil {
		return nil, err
	}

	return dt, nil
}

package teamcity

import (
	"fmt"
	"hash/crc32"
	"regexp"
	"strings"

	api "github.com/cvbarros/go-teamcity/teamcity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Delete: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGroupImport,
		},

		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	var key, name, description string

	if v, ok := d.GetOk("key"); ok {
		key = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		description = v.(string)
	}

	if key == "" {
		generateKey, err := generateKey(name)

		if err != nil {
			return err
		}
		key = *generateKey
	}

	newGroup, err := api.NewGroup(key, name, description)
	if err != nil {
		return err
	}

	created, err := client.Groups.Create(newGroup)
	if err != nil {
		return err
	}

	d.MarkNewResource()
	d.SetId(created.Key)

	return resourceGroupRead(d, meta)
}

func generateKey(name string) (*string, error) {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")

	if err != nil {
		return nil, err
	}

	processedName := reg.ReplaceAllString(strings.ToUpper(name), "")
	generatedKey := fmt.Sprintf("%0.7s_%X", processedName, crc32.ChecksumIEEE([]byte(name)))
	return &generatedKey, nil
}

func resourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	dt, err := client.Groups.GetByKey(d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return nil
		}
		return err
	}
	if err := d.Set("key", dt.Key); err != nil {
		return err
	}
	if err := d.Set("name", dt.Name); err != nil {
		return err
	}
	if err := d.Set("description", dt.Description); err != nil {
		return err
	}

	return nil
}

func resourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	return client.Groups.Delete(d.Id())
}

func resourceGroupImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceGroupRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

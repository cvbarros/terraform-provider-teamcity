package teamcity

import (
	"fmt"
	"strings"

	api "github.com/cvbarros/go-teamcity/teamcity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceGroupRoleAssignment() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupRoleAssignmentCreate,
		Read:   resourceGroupRoleAssignmentRead,
		Delete: resourceGroupRoleAssignmentDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGroupRoleAssignmentImport,
		},

		Schema: map[string]*schema.Schema{
			"group_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceGroupRoleAssignmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	var groupKey, roleID, projectID string

	if v, ok := d.GetOk("group_key"); ok {
		groupKey = v.(string)
	}

	if v, ok := d.GetOk("role_id"); ok {
		roleID = v.(string)
	}

	if v, ok := d.GetOk("project_id"); ok {
		projectID = v.(string)
	}

	newGroupRoleAssignment, err := api.NewGroupRoleAssignment(groupKey, roleID, projectID)
	if err != nil {
		return err
	}

	_, err = client.RoleAssignments.AssignToGroup(newGroupRoleAssignment)
	if err != nil {
		return err
	}

	d.MarkNewResource()
	d.SetId(fmt.Sprintf("%s/%s/%s", groupKey, roleID, projectID))

	return resourceGroupRoleAssignmentRead(d, meta)
}

func resourceGroupRoleAssignmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	newGroupRoleAssignment, err := createGroupRoleAssignmentFromResourceData(d)
	if err != nil {
		return err
	}

	dt, err := client.RoleAssignments.GetForGroup(newGroupRoleAssignment)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return nil
		}
		return err
	}
	if err := d.Set("group_key", newGroupRoleAssignment.GroupKey); err != nil {
		return err
	}
	if err := d.Set("role_id", dt.RoleID); err != nil {
		return err
	}
	if err := d.Set("project_id", dt.Scope); err != nil {
		return err
	}

	return nil
}

func resourceGroupRoleAssignmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	newGroupRoleAssignment, err := createGroupRoleAssignmentFromResourceData(d)
	if err != nil {
		return err
	}

	return client.RoleAssignments.UnassignFromGroup(newGroupRoleAssignment)
}

func resourceGroupRoleAssignmentImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceGroupRoleAssignmentRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func createGroupRoleAssignmentFromResourceData(d *schema.ResourceData) (*api.GroupRoleAssignment, error) {
	parts := strings.Split(d.Id(), "/")
	groupKey := parts[0]
	roleID := parts[1]
	scope := parts[2]

	newGroupRoleAssignment, err := api.NewGroupRoleAssignment(groupKey, roleID, scope)
	if err != nil {
		return nil, err
	}

	return newGroupRoleAssignment, nil
}

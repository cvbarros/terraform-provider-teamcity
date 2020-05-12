package teamcity_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	api "github.com/cvbarros/go-teamcity/teamcity"
)

func TestAccGroupRoleAssignmentAssign_SysAdmin(t *testing.T) {
	var r api.RoleAssignmentReference
	resName := "teamcity_group_role_assignment.test_group_1_sys_admin_global"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGroupRoleAssignmentDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccGroupRoleAssignmentConfigSysAdmin,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupRoleAssignmentExists(resName, &r),
					resource.TestCheckResourceAttr(resName, "group_key", generateKey("Test Group #1")),
					resource.TestCheckResourceAttr(resName, "role_id", "SYSTEM_ADMIN"),
					resource.TestCheckResourceAttr(resName, "project_id", "g"),
				),
			},
		},
	})
}

func TestAccGroupRoleAssignmentAssign_ProjDev(t *testing.T) {
	var r api.RoleAssignmentReference
	resName := "teamcity_group_role_assignment.test_group_2_project_dev_test_project"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGroupRoleAssignmentDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccGroupRoleAssignmentConfigProjDev,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupRoleAssignmentExists(resName, &r),
					resource.TestCheckResourceAttr(resName, "group_key", generateKey("Test Group #2")),
					resource.TestCheckResourceAttr(resName, "role_id", "PROJECT_DEVELOPER"),
					resource.TestCheckResourceAttr(resName, "project_id", "p:TestProject"),
				),
			},
		},
	})
}

func testAccCheckGroupRoleAssignmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	return groupRoleAssignmentDestroyHelper(s, client)
}

func groupRoleAssignmentDestroyHelper(s *terraform.State, client *api.Client) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != "teamcity_group_role_assignment" {
			continue
		}

		groupRoleAssignment, _ := createGroupRoleAssignmentFromResourceData(r.Primary.ID)
		_, err := client.RoleAssignments.GetForGroup(groupRoleAssignment)

		if err != nil {
			if strings.Contains(err.Error(), "404") {
				continue
			}
			return fmt.Errorf("Received an error retrieving the GroupRoleAssignment: %s", err)
		}

		return fmt.Errorf("GroupRoleAssignment still exists")
	}
	return nil
}

func testAccCheckGroupRoleAssignmentExists(n string, out *api.RoleAssignmentReference) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		return groupRoleAssignmentExistsHelper(n, s, client, out)
	}
}

func groupRoleAssignmentExistsHelper(n string, s *terraform.State, client *api.Client, out *api.RoleAssignmentReference) error {
	rs, ok := s.RootModule().Resources[n]
	if !ok {
		return fmt.Errorf("Not found: %s", n)
	}

	if rs.Primary.ID == "" {
		return fmt.Errorf("No id for %s is set", n)
	}

	groupRoleAssignment, _ := createGroupRoleAssignmentFromResourceData(rs.Primary.ID)
	resp, err := client.RoleAssignments.GetForGroup(groupRoleAssignment)

	if err != nil {
		return fmt.Errorf("Received an error retrieving Groups: %s", err)
	}

	*out = *resp

	return nil
}

func createGroupRoleAssignmentFromResourceData(id string) (*api.GroupRoleAssignment, error) {
	parts := strings.Split(id, "/")
	groupKey := parts[0]
	roleID := parts[1]
	scope := parts[2]

	newGroupRoleAssignment, err := api.NewGroupRoleAssignment(groupKey, roleID, scope)
	if err != nil {
		return nil, err
	}

	return newGroupRoleAssignment, nil
}

const TestAccGroupRoleAssignmentConfigSysAdmin = `
resource "teamcity_group" "test_group_1" {
  name = "Test Group #1"
}

resource "teamcity_group_role_assignment" "test_group_1_sys_admin_global" {
  group_key  = teamcity_group.test_group_1.id
  role_id    = "SYSTEM_ADMIN"
  project_id = "g"
}
`

const TestAccGroupRoleAssignmentConfigProjDev = `
resource "teamcity_group" "test_group_2" {
  name = "Test Group #2"
}

resource "teamcity_project" "test_project" {
  name = "Test Project"
}

resource "teamcity_group_role_assignment" "test_group_2_project_dev_test_project" {
  group_key  = teamcity_group.test_group_2.id
  role_id    = "PROJECT_DEVELOPER"
  project_id = "p:${teamcity_project.test_project.id}"
}
`

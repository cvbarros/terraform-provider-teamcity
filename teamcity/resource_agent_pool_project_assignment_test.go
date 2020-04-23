package teamcity_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	api "github.com/cvbarros/go-teamcity/teamcity"
	"github.com/cvbarros/terraform-provider-teamcity/teamcity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTeamCityAgentPoolProjectAssignment_Basic(t *testing.T) {
	resName := "teamcity_agent_pool_project_assignment.test"
	ri := int(time.Now().Unix())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamCityAgentPoolProjectAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamCityAgentPoolProjectAssignmentBasic(ri),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamCityAgentPoolProjectAssignmentExists(resName),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"disassociate_from_other_pools", // a create-only field which isn't set into the state
				},
			},
		},
	})
}

func TestAccTeamCityAgentPoolProjectAssignment_DisassocateOthers(t *testing.T) {
	resName := "teamcity_agent_pool_project_assignment.test"
	ri := int(time.Now().Unix())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamCityAgentPoolProjectAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamCityAgentPoolProjectAssignmentDisassocateFromOtherPools(ri),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamCityAgentPoolProjectAssignmentExists(resName),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"disassociate_from_other_pools", // a create-only field which isn't set into the state
				},
			},
			{
				// then remove the association, but leave the Project and Agent Pool
				// which means we can check this has been associated back to the Default pool
				Config: testAccTeamCityAgentPoolProjectAssignmentTemplate(ri),
				Check: resource.ComposeTestCheckFunc(
					// NOTE: this is intentionally checking the Project resource here, since the assignment no longer exists
					testAccCheckTeamCityAgentPoolProjectAssignmentOnlyContains("teamcity_project.test", "Default"),
				),
			},
		},
	})
}

func testAccCheckTeamCityAgentPoolProjectAssignmentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := teamcity.ParseAgentPoolProjectAssignmentID(rs.Primary.ID)
		if err != nil {
			return err
		}

		agentPoolsForProject, err := client.AgentPools.ListForProject(id.ProjectId)
		if err != nil {
			return fmt.Errorf("Received an error retrieving agent pools for project: %s", err)
		}

		exists := false
		for _, pool := range agentPoolsForProject.AgentPools {
			if pool.Id == id.AgentPoolId {
				exists = true
				break
			}
		}

		if !exists {
			return fmt.Errorf("Agent Pool - Project Association was not found!")
		}

		return nil
	}
}

func testAccCheckTeamCityAgentPoolProjectAssignmentOnlyContains(resourceName string, agentPoolName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		agentPoolsForProject, err := client.AgentPools.ListForProject(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Received an error retrieving agent pools for project: %s", err)
		}

		if len(agentPoolsForProject.AgentPools) != 1 {
			return fmt.Errorf("Expected 1 assignment but got %d", len(agentPoolsForProject.AgentPools))
		}

		if agentPoolsForProject.AgentPools[0].Name != agentPoolName {
			return fmt.Errorf("Expected the agent pool %q to be assigned but got %q", agentPoolName, agentPoolsForProject.AgentPools[0].Name)
		}

		return nil
	}
}

func testAccCheckTeamCityAgentPoolProjectAssignmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "teamcity_agent_pool_project_assignment" {
			continue
		}

		id, err := teamcity.ParseAgentPoolProjectAssignmentID(rs.Primary.ID)
		if err != nil {
			return err
		}

		agentPoolsForProject, err := client.AgentPools.ListForProject(id.ProjectId)
		if err != nil {
			if strings.Contains(err.Error(), "400 (Bad Request)") {
				// gone
				return nil
			}

			return fmt.Errorf("Received an error retrieving agent pools for project: %s", err)
		}

		exists := false
		for _, pool := range agentPoolsForProject.AgentPools {
			if pool.Id == id.AgentPoolId {
				exists = true
				break
			}
		}

		if exists {
			return fmt.Errorf("Agent Pool - Project Association still exists!")
		}
	}
	return nil
}

func testAccTeamCityAgentPoolProjectAssignmentBasic(rInt int) string {
	template := testAccTeamCityAgentPoolProjectAssignmentTemplate(rInt)
	return fmt.Sprintf(`
%s

resource "teamcity_agent_pool_project_assignment" "test" {
  agent_pool_id = teamcity_agent_pool.test.id
  project_id    = teamcity_project.test.id
}
`, template)
}

func testAccTeamCityAgentPoolProjectAssignmentDisassocateFromOtherPools(rInt int) string {
	template := testAccTeamCityAgentPoolProjectAssignmentTemplate(rInt)
	return fmt.Sprintf(`
%s

resource "teamcity_agent_pool_project_assignment" "test" {
  agent_pool_id                 = teamcity_agent_pool.test.id
  project_id                    = teamcity_project.test.id
  disassociate_from_other_pools = true
}
`, template)
}

func testAccTeamCityAgentPoolProjectAssignmentTemplate(rInt int) string {
	return fmt.Sprintf(`
resource "teamcity_project" "test" {
  name = "Project %d"
}

resource "teamcity_agent_pool" "test" {
  name = "Pool %d"
}
`, rInt, rInt)
}

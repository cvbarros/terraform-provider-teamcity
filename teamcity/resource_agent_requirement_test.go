package teamcity_test

import (
	"fmt"
	"strings"
	"testing"

	api "github.com/cvbarros/go-teamcity-sdk/teamcity"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTeamcityAgentRequirement_Basic(t *testing.T) {
	resName := "teamcity_agent_requirement.test"
	var out api.AgentRequirement
	var bc api.BuildType

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamcityAgentRequirementDestroy(&out.BuildTypeID),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccAgentRequirementBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists("teamcity_build_config.config", &bc),
					testAccCheckTeamcityAgentRequirementExists(resName, &bc.ID, &out),
					resource.TestCheckResourceAttrPtr(resName, "build_config_id", &out.BuildTypeID),
					resource.TestCheckResourceAttr(resName, "condition", api.Conditions.Equals),
					resource.TestCheckResourceAttr(resName, "name", "agent_condition"),
					resource.TestCheckResourceAttr(resName, "value", "somevalue"),
				),
			},
		},
	})
}

func TestAccTeamcityAgentRequirement_Update(t *testing.T) {
	resName := "teamcity_agent_requirement.test"
	var out api.AgentRequirement
	var bc api.BuildType

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamcityAgentRequirementDestroy(&out.BuildTypeID),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccAgentRequirementBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists("teamcity_build_config.config", &bc),
					testAccCheckTeamcityAgentRequirementExists(resName, &bc.ID, &out),
					resource.TestCheckResourceAttrPtr(resName, "build_config_id", &out.BuildTypeID),
					resource.TestCheckResourceAttr(resName, "condition", api.Conditions.Equals),
					resource.TestCheckResourceAttr(resName, "name", "agent_condition"),
					resource.TestCheckResourceAttr(resName, "value", "somevalue"),
				),
			},
			resource.TestStep{
				Config: TestAccAgentRequirementUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists("teamcity_build_config.config", &bc),
					testAccCheckTeamcityAgentRequirementExists(resName, &bc.ID, &out),
					resource.TestCheckResourceAttrPtr(resName, "build_config_id", &out.BuildTypeID),
					resource.TestCheckResourceAttr(resName, "condition", api.Conditions.DoesNotEqual),
					resource.TestCheckResourceAttr(resName, "name", "updated_condition"),
					resource.TestCheckResourceAttr(resName, "value", "updated_value"),
				),
			},
		},
	})
}

func testAccCheckTeamcityAgentRequirementDestroy(bt *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		return AgentRequirementDestroyHelper(s, bt, client)
	}
}

func AgentRequirementDestroyHelper(s *terraform.State, bt *string, client *api.Client) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != "teamcity_agent_requirement" {
			continue
		}

		srv := client.AgentRequirementService(*bt)
		_, err := srv.GetByID(r.Primary.ID)

		if err != nil {
			if strings.Contains(err.Error(), "404") {
				continue
			}
			return fmt.Errorf("Received an error retrieving the AgentRequirement: %s", err)
		}

		return fmt.Errorf("AgentRequirement still exists")
	}
	return nil
}

func testAccCheckTeamcityAgentRequirementExists(n string, bt *string, snap *api.AgentRequirement) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		return teamcityAgentRequirementExistsHelper(n, bt, s, client, snap)
	}
}

func teamcityAgentRequirementExistsHelper(n string, bt *string, s *terraform.State, client *api.Client, snap *api.AgentRequirement) error {
	rs, ok := s.RootModule().Resources[n]
	if !ok {
		return fmt.Errorf("Not found: %s", n)
	}

	if rs.Primary.ID == "" {
		return fmt.Errorf("No ID is set")
	}

	srv := client.AgentRequirementService(*bt)
	out, err := srv.GetByID(rs.Primary.ID)
	if err != nil {
		return fmt.Errorf("Received an error retrieving Agent Requirement: %s", err)
	}

	*snap = *out
	return nil
}

const TestAccAgentRequirementBasic = `
resource "teamcity_project" "agentrequirement_project_test" {
  name = "AgentRequirementProject"
}

resource "teamcity_build_config" "config" {
	name = "BuildConfig"
	project_id = "${teamcity_project.agentrequirement_project_test.id}"
}

resource "teamcity_agent_requirement" "test" {
	build_config_id = "${teamcity_build_config.config.id}"
	condition = "equals"
	name = "agent_condition"
	value = "somevalue"
}
`

const TestAccAgentRequirementUpdated = `
resource "teamcity_project" "agentrequirement_project_test" {
  name = "AgentRequirementProject"
}

resource "teamcity_build_config" "config" {
	name = "BuildConfig"
	project_id = "${teamcity_project.agentrequirement_project_test.id}"
}

resource "teamcity_agent_requirement" "test" {
	build_config_id = "${teamcity_build_config.config.id}"
	condition = "does-not-equal"
	name = "updated_condition"
	value = "updated_value"
}
`

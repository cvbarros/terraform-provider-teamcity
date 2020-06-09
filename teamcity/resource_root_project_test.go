package teamcity_test

import (
	"fmt"
	"strings"
	"testing"

	api "github.com/cvbarros/go-teamcity/teamcity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccRootProject_Basic(t *testing.T) {
	resName := "teamcity_root_project.root"
	var p api.Project

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRootProjectDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccTeamcityRootProjectConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamcityProjectExists(resName, &p),
					resource.TestCheckResourceAttr(resName, "name", "<Root project>"),
					resource.TestCheckResourceAttr(resName, "description", "Contains all other projects"),
				),
			},
		},
	})
}

func TestAccRootProject_Full(t *testing.T) {
	resName := "teamcity_root_project.root"
	var p api.Project

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRootProjectDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccTeamcityRootProjectFull,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamcityProjectExists(resName, &p),
					resource.TestCheckResourceAttr(resName, "name", "<Root project>"),
					resource.TestCheckResourceAttr(resName, "description", "Contains all other projects"),
					resource.TestCheckResourceAttr(resName, "config_params.param1", "config_value1"),
					resource.TestCheckResourceAttr(resName, "config_params.param2", "config_value2"),
					resource.TestCheckResourceAttr(resName, "env_params.param3", "env_value1"),
					resource.TestCheckResourceAttr(resName, "env_params.param4", "env_value2"),
					resource.TestCheckResourceAttr(resName, "sys_params.param5", "sys_value1"),
					resource.TestCheckResourceAttr(resName, "sys_params.param6", "sys_value2"),
					testAccCheckProjectParameter(&p, api.ParameterTypes.Configuration, "param1", "config_value1"),
					testAccCheckProjectParameter(&p, api.ParameterTypes.Configuration, "param2", "config_value2"),
					testAccCheckProjectParameter(&p, api.ParameterTypes.EnvironmentVariable, "param3", "env_value1"),
					testAccCheckProjectParameter(&p, api.ParameterTypes.EnvironmentVariable, "param4", "env_value2"),
					testAccCheckProjectParameter(&p, api.ParameterTypes.System, "param5", "sys_value1"),
					testAccCheckProjectParameter(&p, api.ParameterTypes.System, "param6", "sys_value2"),
				),
			},
		},
	})
}

func TestAccRootProject_FullUpdate(t *testing.T) {
	resName := "teamcity_root_project.root"
	var p api.Project
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRootProjectDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccTeamcityRootProjectFull,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamcityProjectExists(resName, &p),
					resource.TestCheckResourceAttr(resName, "config_params.param1", "config_value1"),
					resource.TestCheckResourceAttr(resName, "config_params.param2", "config_value2"),
					resource.TestCheckResourceAttr(resName, "env_params.param3", "env_value1"),
					resource.TestCheckResourceAttr(resName, "env_params.param4", "env_value2"),
					resource.TestCheckResourceAttr(resName, "sys_params.param5", "sys_value1"),
					resource.TestCheckResourceAttr(resName, "sys_params.param6", "sys_value2"),
					testAccCheckProjectParameter(&p, api.ParameterTypes.Configuration, "param1", "config_value1"),
					testAccCheckProjectParameter(&p, api.ParameterTypes.Configuration, "param2", "config_value2"),
					testAccCheckProjectParameter(&p, api.ParameterTypes.EnvironmentVariable, "param3", "env_value1"),
					testAccCheckProjectParameter(&p, api.ParameterTypes.EnvironmentVariable, "param4", "env_value2"),
					testAccCheckProjectParameter(&p, api.ParameterTypes.System, "param5", "sys_value1"),
					testAccCheckProjectParameter(&p, api.ParameterTypes.System, "param6", "sys_value2"),
				),
			},
			resource.TestStep{
				Config: testAccTeamcityRootProjectFullUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamcityProjectExists(resName, &p),
					resource.TestCheckResourceAttr(resName, "config_params.param1", "config_value1_updated"),
					resource.TestCheckResourceAttr(resName, "config_params.param2", "config_value2_updated"),
					testAccCheckProjectParameter(&p, api.ParameterTypes.Configuration, "param1", "config_value1_updated"),
					testAccCheckProjectParameter(&p, api.ParameterTypes.Configuration, "param2", "config_value2_updated"),
				),
			},
		},
	})
}

func testAccCheckRootProjectDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	return teamcityRootProjectDestroyHelper(s, client)
}

func teamcityRootProjectDestroyHelper(s *terraform.State, client *api.Client) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != "teamcity_root_project" {
			continue
		}

		dt, err := client.Projects.GetByID(r.Primary.ID)
		// Empty all the parameters, so that it will be destroyed
		dt.Parameters.Items = nil

		_, err = client.Projects.Update(dt)

		if err != nil {
			if strings.Contains(err.Error(), "404") {
				continue
			}
			return fmt.Errorf("Received an error retrieving the project: %s", err)
		}
	}
	return nil
}

const testAccTeamcityRootProjectConfig = `
resource "teamcity_root_project" "root" {
}
`
const testAccTeamcityRootProjectFull = `
resource "teamcity_root_project" "root" {

	config_params = {
		param1 = "config_value1"
		param2 = "config_value2"
	}

	env_params = {
		param3 = "env_value1"
		param4 = "env_value2"
	}

	sys_params = {
		param5 = "sys_value1"
		param6 = "sys_value2"
	}
}
`

const testAccTeamcityRootProjectFullUpdated = `
resource "teamcity_root_project" "root" {

	config_params = {
		param1 = "config_value1_updated"
		param2 = "config_value2_updated"
	}
}
`

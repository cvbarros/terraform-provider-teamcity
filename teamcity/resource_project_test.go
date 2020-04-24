package teamcity_test

import (
	"fmt"
	"strings"
	"testing"

	api "github.com/cvbarros/go-teamcity/teamcity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTeamcityProject_Basic(t *testing.T) {
	resName := "teamcity_project.testproj"
	var p api.Project

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamcityProjectDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccTeamcityProjectConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamcityProjectExists(resName, &p),
					resource.TestCheckResourceAttr(resName, "name", "testproj"),
				),
			},
		},
	})
}

func TestAccTeamcityProject_Full(t *testing.T) {
	resName := "teamcity_project.testproj"
	var p api.Project

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		Providers:     testAccProviders,
		IDRefreshName: resName,
		CheckDestroy:  testAccCheckTeamcityProjectDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccTeamcityProjectFull,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamcityProjectExists(resName, &p),
					resource.TestCheckResourceAttr(resName, "name", "test_project"),
					resource.TestCheckResourceAttr(resName, "description", "Test Project"),
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

func TestAccTeamcityProject_Parent(t *testing.T) {
	parentRes := "teamcity_project.parent"
	childRes := "teamcity_project.child"
	var child, parent api.Project

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamcityProjectDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccTeamcityProjectParent,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamcityProjectExists(parentRes, &parent),
					testAccCheckTeamcityProjectExists(childRes, &child),
					resource.TestCheckResourceAttrPtr(childRes, "parent_id", &parent.ID),
				),
			},
		},
	})
}

func TestAccTeamcityProject_Update(t *testing.T) {
	resName := "teamcity_project.testproj"
	var p api.Project
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamcityProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamcityProjectFull,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamcityProjectExists(resName, &p),
					resource.TestCheckResourceAttr(resName, "description", "Test Project"),
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
			{
				Config: testAccTeamcityProjectFullUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamcityProjectExists(resName, &p),
					resource.TestCheckResourceAttr(resName, "description", "updated project"),
					resource.TestCheckResourceAttr(resName, "description", "Test Project Updated"),
					resource.TestCheckResourceAttr(resName, "config_params.param1", "config_value1"),
					resource.TestCheckResourceAttr(resName, "config_params.param2", "config_value2"),
					testAccCheckProjectParameter(&p, api.ParameterTypes.Configuration, "param1", "config_value1"),
					testAccCheckProjectParameter(&p, api.ParameterTypes.Configuration, "param2", "config_value2"),
				),
			},
		},
	})
}

func testAccCheckProjectParameter(dt *api.Project, paramType string, paramName string, paramValue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, i := range dt.Parameters.Items {
			if i.Type != paramType {
				continue
			}
			if i.Name == paramName {
				if i.Value == paramValue {
					return nil
				} else {
					return fmt.Errorf("param '%s' has a wrong value. expected: %s, actual: %s", paramName, paramValue, i.Value)
				}
			}
		}
		return fmt.Errorf("parameter named '%s' not found with type '%s'", paramName, paramType)
	}
}

func testAccCheckTeamcityProjectExists(n string, project *api.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		return teamcityProjectExistsHelper(n, s, client, project)
	}
}

// func testAccCheckTeamcityProject

func teamcityProjectExistsHelper(n string, s *terraform.State, client *api.Client, p *api.Project) error {
	rs, ok := s.RootModule().Resources[n]
	if !ok {
		return fmt.Errorf("Not found: %s", n)
	}

	if rs.Primary.ID == "" {
		return fmt.Errorf("No ID is set")
	}

	proj, err := client.Projects.GetByID(rs.Primary.ID)
	if err != nil {
		return fmt.Errorf("Received an error retrieving project: %s", err)
	}

	*p = *proj
	return nil
}

func testAccCheckTeamcityProjectDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	return teamcityProjectDestroyHelper(s, client)
}

func teamcityProjectDestroyHelper(s *terraform.State, client *api.Client) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != "teamcity_project" {
			continue
		}

		_, err := client.Projects.GetByID(r.Primary.ID)

		if err != nil {
			if strings.Contains(err.Error(), "404") {
				continue
			}
			return fmt.Errorf("Received an error retrieving the project: %s", err)
		}

		return fmt.Errorf("Project still exists")
	}
	return nil
}

const testAccTeamcityProjectConfig = `
resource "teamcity_project" "testproj" {
  name = "testproj"
}
`

const testAccTeamcityProjectParent = `
resource "teamcity_project" "parent" {
	name = "parent"
}

resource "teamcity_project" "child" {
	name = "child"
	parent_id = "${teamcity_project.parent.id}"
}
`

const testAccTeamcityProjectFull = `
resource "teamcity_project" "testproj" {
	name = "test_project"
	description = "Test Project"

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

const testAccTeamcityProjectFullUpdated = `
resource "teamcity_project" "testproj" {
	name = "updated project"
	description = "Test Project Updated"

	config_params = {
		param1 = "config_value1"
		param2 = "config_value2"
	}
}
`

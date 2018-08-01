package teamcity_test

import (
	"fmt"
	"strings"
	"testing"

	api "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccBuildConfig_Basic(t *testing.T) {
	var bc api.BuildType

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildConfigDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccBuildConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists("teamcity_build_config.build_configuration_test", &bc),
					resource.TestCheckResourceAttr(
						"teamcity_build_config.build_configuration_test", "name", "build config test",
					),
					resource.TestCheckResourceAttr(
						"teamcity_build_config.build_configuration_test", "description", "build config test desc",
					),
					resource.TestCheckResourceAttr(
						"teamcity_build_config.build_configuration_test", "project_id", "BuildConfigProjectTest",
					),
				),
			},
		},
	})
}

func TestAccBuildConfig_Steps(t *testing.T) {
	var bc api.BuildType
	resName := "teamcity_build_config.build_configuration_test"
	expected := map[string]string{
		"name": "build release",
		"type": api.StepTypePowershell,
		"file": "build.ps1",
		"args": "-Target buildrelease",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildConfigDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccBuildConfigSteps,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists(resName, &bc),
					testAccCheckStepExists(&bc.ID, expected),
				),
			},
		},
	})
}

func TestAccBuildConfig_Parameters(t *testing.T) {
	var bc api.BuildType
	resName := "teamcity_build_config.build_configuration_test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildConfigDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccBuildConfigParams,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists(resName, &bc),
					testAccCheckProperties(&bc.Parameters, "env.DEPLOY_SERVER", "server.com"),
					testAccCheckProperties(&bc.Parameters, "env.some_variable", "hello"),
				),
			},
		},
	})
}

func TestAccBuildConfig_VcsRoot(t *testing.T) {
	var bc api.BuildType
	resName := "teamcity_build_config.build_configuration_test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildConfigDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccBuildConfigVcsRoot,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists(resName, &bc),
					testAccCheckVcsRootAttached(&bc.VcsRootEntries, "application", "+:*\\n-:README.MD"),
				),
			},
		},
	})
}

func testAccCheckStepExists(buildTypeID *string, stepExpected map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		steps, err := client.BuildTypes.GetSteps(*buildTypeID)
		if err != nil {
			return fmt.Errorf("error when checking steps: %s", err)
		}

		for _, v := range steps {
			if v.Name() == stepExpected["name"] {
				return assertStepProperties(v, stepExpected)
			}
		}

		return fmt.Errorf("Step named '%s' was not found", stepExpected["name"])
	}
}

func assertStepProperties(actual api.Step, expected map[string]string) error {
	stepType := actual.Type()
	if actual.Type() != expected["type"] {
		return fmt.Errorf("Found step %s but types differ, actual: %s, expected: %s", expected["name"], actual.Type(), expected["type"])
	}

	if stepType == "powershell" {
		psStep := actual.(*api.StepPowershell)
		if p, ok := expected["file"]; ok {
			if p != psStep.ScriptFile {
				return fmt.Errorf("Property 'file' differs, actual: %s, expected: %s", psStep.ScriptFile, p)
			}
		}

		if p, ok := expected["args"]; ok {
			if p != psStep.ScriptArgs {
				return fmt.Errorf("Property 'args' differs, actual: %s, expected: %s", psStep.ScriptArgs, p)
			}
		}
	}

	return nil
}

func getPropertyOk(p *api.Properties, key string) (string, bool) {
	if len(p.Items) == 0 {
		return "", false
	}

	for _, v := range p.Items {
		if v.Name == key {
			return v.Value, true
		}
	}

	return "", false
}

func testAccCheckVcsRootAttached(vcs **api.VcsRootEntries, n string, co string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *vcs == nil {
			return fmt.Errorf("VcsRootEntries must not be nil")
		}

		for _, v := range (*vcs).Items {
			if v.VcsRoot.Name == n {
				if v.CheckoutRules == co {
					return nil
				}
			}
		}

		return fmt.Errorf("VCS Root with name '%s' and checkout rules '%s' was not found", n, co)
	}
}

func testAccCheckBuildConfigExists(n string, out *api.BuildType) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		return buildConfigExistsHelper(n, s, client, out)
	}
}

func buildConfigExistsHelper(n string, s *terraform.State, client *api.Client, out *api.BuildType) error {
	rs, ok := s.RootModule().Resources[n]
	if !ok {
		return fmt.Errorf("Not found: %s", n)
	}

	if rs.Primary.ID == "" {
		return fmt.Errorf("No id for %s is set", n)
	}

	resp, err := client.BuildTypes.GetByID(rs.Primary.ID)

	if err != nil {
		return fmt.Errorf("Received an error retrieving Build Configurationt: %s", err)
	}

	*out = *resp

	return nil
}

func testAccCheckBuildConfigDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	return buildConfigDestroyHelper(s, client)
}

func buildConfigDestroyHelper(s *terraform.State, client *api.Client) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != "teamcity_build_config" {
			continue
		}

		_, err := client.BuildTypes.GetByID(r.Primary.ID)

		if err != nil {
			if strings.Contains(err.Error(), "404") {
				continue
			}
			return fmt.Errorf("Received an error retrieving the Build Configuration: %s", err)
		}

		return fmt.Errorf("Build Configuration still exists")
	}
	return nil
}

// testAccCheckProperties can be used to check the property value for a resource
func testAccCheckProperties(
	props **api.Properties, key string, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if props == nil {
			return fmt.Errorf("Properties must not be nil")
		}

		m := (*props).Map()
		v, ok := m[key]
		if value != "" && !ok {
			return fmt.Errorf("Missing property: %s", key)
		} else if value == "" && ok {
			return fmt.Errorf("Extra property: %s", key)
		}
		if value == "" {
			return nil
		}

		if v != value {
			return fmt.Errorf("%s: bad value: %s", key, v)
		}

		return nil
	}
}

const TestAccBuildConfigBasic = `
resource "teamcity_project" "build_config_project_test" {
  name = "build_config_project_test"
}

resource "teamcity_build_config" "build_configuration_test" {
	name = "build config test"
	project_id = "${teamcity_project.build_config_project_test.id}"
	description = "build config test desc"
}
`

const TestAccBuildConfigParams = `
resource "teamcity_project" "build_config_project_test" {
  name = "build_config_project_test"
}

resource "teamcity_build_config" "build_configuration_test" {
	name = "build config test"
	project_id = "${teamcity_project.build_config_project_test.id}"
	
	env_params {
		DEPLOY_SERVER = "server.com"
		some_variable = "hello"
	}

	config_params {
		github.repository = "nocode"
	}
}
`

const TestAccBuildConfigVcsRoot = `
resource "teamcity_project" "build_config_project_test" {
  name = "build_config_project_test"
}

resource "teamcity_vcs_root_git" "build_config_vcsroot_test" {
	name = "application"
	project_id = "${teamcity_project.build_config_project_test.id}"
	fetch_url = "https://github.com/kelseyhightower/nocode"
	default_branch = "refs/head/master"
}

resource "teamcity_build_config" "build_configuration_test" {
	name = "build config test"
	project_id = "${teamcity_project.build_config_project_test.id}"
	
	vcs_root {
		id = "${teamcity_vcs_root_git.build_config_vcsroot_test.id}"
		checkout_rules = ["+:*", "-:README.MD"]
	}
}
`

const TestAccBuildConfigSteps = `
resource "teamcity_project" "build_config_project_test" {
  name = "build_config_project_test"
}

resource "teamcity_build_config" "build_configuration_test" {
	name = "build config test"
	project_id = "${teamcity_project.build_config_project_test.id}"
	
	step {
		type = "powershell"
		name = "build release"
		file = "build.ps1"
		args = "-Target buildrelease"
	}

	step {
		type = "powershell"
		file = "another.ps1"
	}
}
`

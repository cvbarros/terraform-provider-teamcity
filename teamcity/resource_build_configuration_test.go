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
					resource.TestCheckResourceAttr("teamcity_build_config.build_configuration_test", "name", "build config test"),
					resource.TestCheckResourceAttr("teamcity_build_config.build_configuration_test", "description", "build config test desc"),
					resource.TestCheckResourceAttr("teamcity_build_config.build_configuration_test", "project_id", "BuildConfigProjectTest"),
				),
			},
		},
	})
}

func TestAccBuildConfig_StepsPowershell(t *testing.T) {
	var bc api.BuildType
	resName := "teamcity_build_config.build_configuration_test"
	scriptStep := map[string]string{
		"name": "build_script",
		"type": api.StepTypePowershell,
		"file": "build.ps1",
		"args": "-Target buildrelease",
	}

	codeStep := map[string]string{
		"type": api.StepTypePowershell,
		"name": "build_code",
		"code": "Get-Date",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildConfigDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccBuildConfigStepsPowershell,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists(resName, &bc),
					testAccCheckStepExists(&bc.ID, scriptStep),
					testAccCheckStepExists(&bc.ID, codeStep),
				),
			},
		},
	})
}

func TestAccBuildConfig_StepsCmdLine(t *testing.T) {
	var bc api.BuildType
	resName := "teamcity_build_config.build_configuration_test"

	codeStep := map[string]string{
		"type": api.StepTypeCommandLine,
		"name": "build_script",
		"code": "echo \"Hello World\"",
	}

	exeStep := map[string]string{
		"type": api.StepTypeCommandLine,
		"name": "build_executable",
		"file": "./build.sh",
		"args": "default_target --verbose",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildConfigDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccBuildConfigStepsCmdLine,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists(resName, &bc),
					testAccCheckStepExists(&bc.ID, codeStep),
					testAccCheckStepExists(&bc.ID, exeStep),
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
					resource.TestCheckResourceAttr(resName, "env_params.DEPLOY_SERVER", "server.com"),
					resource.TestCheckResourceAttr(resName, "env_params.some_variable", "hello"),
					resource.TestCheckResourceAttr(resName, "config_params.github.repository", "nocode"),
					resource.TestCheckResourceAttr(resName, "sys_params.system_param", "system_value"),
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

	if stepType == string(api.StepTypePowershell) {
		dt := actual.(*api.StepPowershell)
		if p, ok := expected["file"]; ok {
			if p != dt.ScriptFile {
				return fmt.Errorf("Property 'file' differs, actual: %s, expected: %s", dt.ScriptFile, p)
			}
		}

		if p, ok := expected["args"]; ok {
			if p != dt.ScriptArgs {
				return fmt.Errorf("Property 'args' differs, actual: %s, expected: %s", dt.ScriptArgs, p)
			}
		}

		if p, ok := expected["code"]; ok {
			if p != dt.Code {
				return fmt.Errorf("Property 'code' differs, actual: %s, expected: %s", dt.Code, p)
			}
		}
		return nil
	}

	if stepType == string(api.StepTypeCommandLine) {
		dt := actual.(*api.StepCommandLine)
		if p, ok := expected["file"]; ok {
			if p != dt.CommandExecutable {
				return fmt.Errorf("Property 'file' differs, actual: %s, expected: %s", dt.CommandExecutable, p)
			}
		}

		if p, ok := expected["args"]; ok {
			if p != dt.CommandParameters {
				return fmt.Errorf("Property 'args' differs, actual: %s, expected: %s", dt.CommandParameters, p)
			}
		}

		if p, ok := expected["code"]; ok {
			if p != dt.CustomScript {
				return fmt.Errorf("Property 'code' differs, actual: %s, expected: %s", dt.CustomScript, p)
			}
		}
		return nil
	}

	return fmt.Errorf("Unexpected step type found: %s", stepType)
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

func testAccCheckVcsRootAttached(vcs *[]*api.VcsRootEntry, n string, co string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *vcs == nil {
			return fmt.Errorf("VcsRootEntries must not be nil")
		}

		for _, v := range *vcs {
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
	props **api.Parameters, key string, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if props == nil {
			return fmt.Errorf("Parameters must not be nil")
		}

		m := (*props).Properties().Map()
		v, ok := m[key]
		if value != "" && !ok {
			return fmt.Errorf("Missing parameter: %s", key)
		} else if value == "" && ok {
			return fmt.Errorf("Extra parameter: %s", key)
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

	sys_params {
		system_param = "system_value"
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

const TestAccBuildConfigStepsPowershell = `
resource "teamcity_project" "build_config_project_test" {
  name = "build_config_project_test"
}

resource "teamcity_build_config" "build_configuration_test" {
	name = "build config test"
	project_id = "${teamcity_project.build_config_project_test.id}"
	
	step {
		type = "powershell"
		name = "build_script"
		file = "build.ps1"
		args = "-Target buildrelease"
	}

	step {
		type = "powershell"
		name = "build_code"
		code = "Get-Date"
	}
}
`

const TestAccBuildConfigStepsCmdLine = `
resource "teamcity_project" "build_config_project_test" {
  name = "build_config_project_test"
}

resource "teamcity_build_config" "build_configuration_test" {
	name = "build config test"
	project_id = "${teamcity_project.build_config_project_test.id}"
	
	step {
		type = "cmd_line"
		name = "build_script"
		code = "echo \"Hello World\""
	}

	step {
		type = "cmd_line"
		name = "build_executable"
		file = "./build.sh"
		args = "default_target --verbose"
	}
}
`

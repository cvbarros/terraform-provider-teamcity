package teamcity_test

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	api "github.com/cvbarros/go-teamcity/teamcity"
)

func TestAccBuildConfig_Basic(t *testing.T) {
	var bc api.BuildType
	resName := "teamcity_build_config.build_configuration_test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccBuildConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists(resName, &bc),
					resource.TestCheckResourceAttr(resName, "name", "build config test"),
					resource.TestCheckResourceAttr(resName, "description", "build config test desc"),
					resource.TestCheckResourceAttr(resName, "project_id", "BuildConfigProjectTest"),
					resource.TestCheckResourceAttr(resName, "is_template", "false"),
				),
			},
		},
	})
}

func TestAccBuildConfig_TemplateDoesNotSupportDescription(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      TestAccBuildConfigTemplateWithDescription,
				ExpectError: regexp.MustCompile("'description' field is not supported for Build Configuration Templates. See issue https://youtrack.jetbrains.com/issue/TW-63617 for details"),
			},
		},
	})
}

func TestAccBuildConfig_Template(t *testing.T) {
	var bc api.BuildType
	resName := "teamcity_build_config.build_configuration_test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccBuildConfigTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists(resName, &bc),
					resource.TestCheckResourceAttr(resName, "name", "build config template"),
					resource.TestCheckResourceAttr(resName, "project_id", "BuildConfigProjectTest"),
					resource.TestCheckResourceAttr(resName, "is_template", "true"),
				),
			},
		},
	})
}

func TestAccBuildConfig_BasicBuildCounter(t *testing.T) {
	var bc api.BuildType
	resName := "teamcity_build_config.build_configuration_test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccBuildConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists(resName, &bc),
					resource.TestCheckResourceAttr(resName, "name", "build config test"),
					resource.TestCheckResourceAttr(resName, "description", "build config test desc"),
					resource.TestCheckResourceAttr(resName, "project_id", "BuildConfigProjectTest"),
				),
			},
			{
				PreConfig:          func() { updateBuildCounter(&bc, 10) }, //Simulate external computed
				Config:             TestAccBuildConfigBasic,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}

func TestAccBuildConfig_UpdateBuildCounter(t *testing.T) {
	var bc api.BuildType
	resName := "teamcity_build_config.build_config"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccBuildConfigBuildCounter,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists(resName, &bc),
					resource.TestCheckResourceAttr(resName, "settings.3691833092.build_counter", "2"),
				),
			},
			{
				Config: TestAccBuildConfigBuildCounterUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists(resName, &bc),
					resource.TestCheckResourceAttr(resName, "settings.1862934543.build_counter", "10"),
				),
			},
			{
				PreConfig:          func() { updateBuildCounter(&bc, 20) }, //Simulate external computed
				Config:             TestAccBuildConfigBuildCounterUpdated,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}

func TestAccBuildConfig_UpdateOtherSetting(t *testing.T) {
	var bc api.BuildType
	resName := "teamcity_build_config.build_configuration_test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccBuildConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists(resName, &bc),
					resource.TestCheckResourceAttr(resName, "settings.2597056888.build_number_format", "2.0.%build.counter%"),
					resource.TestCheckResourceAttr(resName, "settings.2597056888.build_counter", "0"),
				),
			},
			{
				Config: TestAccBuildConfigBasicUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists(resName, &bc),
					resource.TestCheckResourceAttr(resName, "settings.1873735866.build_counter", "0"),
					resource.TestCheckResourceAttr(resName, "settings.1873735866.build_number_format", "3.0.%build.counter%"),
				),
			},
		},
	})
}

func TestAccBuildConfig_NestedProject(t *testing.T) {
	var bc api.BuildType
	resName := "teamcity_build_config.build_config"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccBuildConfigurationIdWithParent,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists(resName, &bc),
					resource.TestCheckResourceAttr(resName, "id", "Parent_Child_BuildConfig"),
				),
			},
		},
	})
}

func TestAccBuildConfig_UpdateBasic(t *testing.T) {
	var bc api.BuildType
	resName := "teamcity_build_config.build_configuration_test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccBuildConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists(resName, &bc),
					resource.TestCheckResourceAttr(resName, "name", "build config test"),
					resource.TestCheckResourceAttr(resName, "description", "build config test desc"),
					resource.TestCheckResourceAttr(resName, "project_id", "BuildConfigProjectTest"),
				),
			},
			{
				Config: TestAccBuildConfigBasicUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists(resName, &bc),
					resource.TestCheckResourceAttr(resName, "name", "build config updated"),
					resource.TestCheckResourceAttr(resName, "description", "build config test desc updated"),
				),
			},
		},
	})
}

const TestAccBuildConfigAddAndRemoveSteps = `
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
		type = "cmd_line"
		name = "removed_step"
		code = "echo"
	}
}
`

const TestAccBuildConfigAddAndRemoveStepsUpdated = `
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

func TestAccBuildConfig_StepsAddAndRemove(t *testing.T) {
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
	removedStep := map[string]string{
		"type": api.StepTypeCommandLine,
		"name": "removed_step",
		"code": "echo",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccBuildConfigAddAndRemoveSteps,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists(resName, &bc),
					resource.TestCheckResourceAttrSet(resName, "step.0.step_id"),
					resource.TestCheckResourceAttrSet(resName, "step.1.step_id"),
					testAccCheckStepExists(&bc.ID, scriptStep),
					testAccCheckStepExists(&bc.ID, removedStep)),
			},
			{
				Config: TestAccBuildConfigAddAndRemoveStepsUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists(resName, &bc),
					resource.TestCheckResourceAttrSet(resName, "step.0.step_id"),
					resource.TestCheckResourceAttr(resName, "step.0.name", "build_script"),
					resource.TestCheckResourceAttrSet(resName, "step.1.step_id"),
					resource.TestCheckResourceAttr(resName, "step.1.name", "build_code"),
					testAccCheckStepExists(&bc.ID, scriptStep),
					testAccCheckStepExists(&bc.ID, codeStep),
					testAccCheckStepRemoved(&bc.ID, removedStep)),
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
			{
				Config: TestAccBuildConfigStepsPowershell,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists(resName, &bc),
					resource.TestCheckResourceAttrSet(resName, "step.0.step_id"),
					resource.TestCheckResourceAttrSet(resName, "step.1.step_id"),
					testAccCheckStepExists(&bc.ID, scriptStep),
					testAccCheckStepExists(&bc.ID, codeStep),
				),
			},
		},
	})
}

func TestAccBuildConfig_UpdateSteps(t *testing.T) {
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

	scriptStepUpdate := map[string]string{
		"name": "updated_script",
		"type": api.StepTypePowershell,
		"file": "updated.ps1",
		"args": "-Target pullrequest",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccBuildConfigStepsPowershell,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists(resName, &bc),
					testAccCheckStepExists(&bc.ID, scriptStep),
					testAccCheckStepExists(&bc.ID, codeStep),
				),
			},
			{
				Config: TestAccBuildConfigStepsPowershellUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStepExists(&bc.ID, scriptStepUpdate),
					testAccCheckStepRemoved(&bc.ID, codeStep),
					resource.TestCheckResourceAttr(resName, "step.0.file", "updated.ps1"),
					resource.TestCheckResourceAttr(resName, "step.0.name", "updated_script"),
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
			{
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

func TestAccBuildConfig_StepsCmdLineUpdateSteps(t *testing.T) {
	var bc api.BuildType
	resName := "teamcity_build_config.build_configuration_test"
	scriptStep := map[string]string{
		"name": "build_script",
		"type": api.StepTypeCommandLine,
		"code": "echo \"Hello World\"",
	}

	scriptStepUpdate := map[string]string{
		"name": "build_script",
		"type": api.StepTypeCommandLine,
		"code": "echo \"Hello Foo\"",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccBuildConfigStepsCmdLine,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists(resName, &bc),
					testAccCheckStepExists(&bc.ID, scriptStep),
					resource.TestCheckResourceAttr(resName, "step.0.code", "echo \"Hello World\""),
				),
			},
			{
				Config: TestAccBuildConfigStepsCmdLineUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStepExists(&bc.ID, scriptStepUpdate),
					resource.TestCheckResourceAttr(resName, "step.0.code", "echo \"Hello Foo\""),
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
			{
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

func TestAccBuildConfig_UpdateParameters(t *testing.T) {
	var bc api.BuildType
	resName := "teamcity_build_config.build_configuration_test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccBuildConfigParams,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists(resName, &bc),
					resource.TestCheckResourceAttr(resName, "config_params.%", "1"),
					resource.TestCheckResourceAttr(resName, "env_params.%", "2"),
					resource.TestCheckResourceAttr(resName, "sys_params.%", "1"),
					resource.TestCheckResourceAttr(resName, "env_params.DEPLOY_SERVER", "server.com"),
					resource.TestCheckResourceAttr(resName, "env_params.some_variable", "hello"),
					resource.TestCheckResourceAttr(resName, "config_params.github.repository", "nocode"),
					resource.TestCheckResourceAttr(resName, "sys_params.system_param", "system_value"),
				),
			},
			{
				Config: TestAccBuildConfigParamsUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists(resName, &bc),
					resource.TestCheckResourceAttr(resName, "config_params.%", "1"),
					resource.TestCheckResourceAttr(resName, "env_params.%", "2"),
					resource.TestCheckResourceAttr(resName, "sys_params.%", "0"),
					resource.TestCheckResourceAttr(resName, "env_params.DEPLOY_SERVER", "server.com"),
					resource.TestCheckResourceAttr(resName, "env_params.some_variable", "hello"),
					resource.TestCheckResourceAttr(resName, "config_params.github.repository", "updated_repo"),
				),
			},
		},
	})
}

func TestAccBuildConfig_Settings(t *testing.T) {
	var bc api.BuildType
	resName := "teamcity_build_config.build_configuration_test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccBuildConfigSettings,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists(resName, &bc),
					resource.TestCheckResourceAttr(resName, "settings.3961488219.configuration_type", "DEPLOYMENT"),
					resource.TestCheckResourceAttr(resName, "settings.3961488219.allow_personal_builds", "true"),
					resource.TestCheckResourceAttr(resName, "settings.3961488219.artifact_paths.#", "1"),
					resource.TestCheckResourceAttr(resName, "settings.3961488219.artifact_paths.0", "+:*.json => /config/*.json"),
					resource.TestCheckResourceAttr(resName, "settings.3961488219.build_counter", "20"),
					resource.TestCheckResourceAttr(resName, "settings.3961488219.build_number_format", "1.0.%build.counter%"),
					resource.TestCheckResourceAttr(resName, "settings.3961488219.concurrent_limit", "10"),
					resource.TestCheckResourceAttr(resName, "settings.3961488219.detect_hanging", "true"),
					resource.TestCheckResourceAttr(resName, "settings.3961488219.status_widget", "false"),
				),
			},
		},
	})
}

func TestAccBuildConfig_SettingsUpdate(t *testing.T) {
	var bc api.BuildType
	resName := "teamcity_build_config.build_configuration_test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccBuildConfigSettings,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists(resName, &bc),
					resource.TestCheckResourceAttr(resName, "settings.3961488219.configuration_type", "DEPLOYMENT"),
					resource.TestCheckResourceAttr(resName, "settings.3961488219.allow_personal_builds", "true"),
					resource.TestCheckResourceAttr(resName, "settings.3961488219.artifact_paths.#", "1"),
					resource.TestCheckResourceAttr(resName, "settings.3961488219.artifact_paths.0", "+:*.json => /config/*.json"),
					resource.TestCheckResourceAttr(resName, "settings.3961488219.build_counter", "20"),
					resource.TestCheckResourceAttr(resName, "settings.3961488219.build_number_format", "1.0.%build.counter%"),
					resource.TestCheckResourceAttr(resName, "settings.3961488219.concurrent_limit", "10"),
					resource.TestCheckResourceAttr(resName, "settings.3961488219.detect_hanging", "true"),
					resource.TestCheckResourceAttr(resName, "settings.3961488219.status_widget", "false"),
				),
			},
			{
				Config: TestAccBuildConfigSettingsUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists(resName, &bc),
					resource.TestCheckResourceAttr(resName, "settings.3684435053.configuration_type", "REGULAR"),
					resource.TestCheckResourceAttr(resName, "settings.3684435053.allow_personal_builds", "false"),
					resource.TestCheckResourceAttr(resName, "settings.3684435053.artifact_paths.#", "1"),
					resource.TestCheckResourceAttr(resName, "settings.3684435053.artifact_paths.0", "+:*.json => /artifacts/*.json"),
					resource.TestCheckResourceAttr(resName, "settings.3684435053.build_counter", "25"),
					resource.TestCheckResourceAttr(resName, "settings.3684435053.build_number_format", "2.0.%build.counter%"),
					resource.TestCheckResourceAttr(resName, "settings.3684435053.concurrent_limit", "0"),
					resource.TestCheckResourceAttr(resName, "settings.3684435053.detect_hanging", "false"),
					resource.TestCheckResourceAttr(resName, "settings.3684435053.status_widget", "true"),
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
			{
				Config: TestAccBuildConfigVcsRoot,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists(resName, &bc),
					testAccCheckVcsRootAttached(&bc.VcsRootEntries, "application", "+:*\\n-:README.MD"),
				),
			},
		},
	})
}

func TestAccBuildConfig_AttachTemplates(t *testing.T) {
	var bc, t1, t2 api.BuildType
	var t3 api.BuildType
	resName := "teamcity_build_config.build_configuration_test"
	template1 := "teamcity_build_config.build_configuration_template1"
	template2 := "teamcity_build_config.build_configuration_template2"
	template3 := "teamcity_build_config.build_configuration_template3"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccBuildConfigAttachTemplates,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists(resName, &bc),
					testAccCheckBuildConfigExists(template1, &t1),
					testAccCheckBuildConfigExists(template2, &t2),
					resource.TestCheckResourceAttr(resName, "templates.#", "2"),
					resource.TestCheckResourceAttrPtr(resName, "templates.0", &t1.ID),
					resource.TestCheckResourceAttrPtr(resName, "templates.1", &t2.ID),
				),
			},
			{
				Config: TestAccBuildConfigUpdatedTemplates,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists(template3, &t3),
					resource.TestCheckResourceAttr(resName, "templates.#", "2"),
					resource.TestCheckResourceAttrPtr(resName, "templates.0", &t1.ID),
					resource.TestCheckResourceAttrPtr(resName, "templates.1", &t3.ID),
				),
			},
		},
	})
}

func testAccCheckStepRemoved(buildTypeID *string, stepRemoved map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		exists, _ := testStepExists(client, *buildTypeID, stepRemoved)
		if exists {
			return fmt.Errorf("expected step %s to be removed, but still exists", stepRemoved["name"])
		}
		return nil
	}
}

func testAccCheckStepExists(buildTypeID *string, stepExpected map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		_, err := testStepExists(client, *buildTypeID, stepExpected)
		return err
	}
}

func testStepExists(client *api.Client, buildTypeID string, stepExpected map[string]string) (bool, error) {
	steps, err := client.BuildTypes.GetSteps(buildTypeID)
	if err != nil {
		return false, fmt.Errorf("error when checking steps: %s", err)
	}

	for _, v := range steps {
		if v.GetName() == stepExpected["name"] {
			err := assertStepProperties(v, stepExpected)
			return err != nil, err
		}
	}

	return false, fmt.Errorf("the step named '%s' was not found", stepExpected["name"])
}

func assertStepProperties(actual api.Step, expected map[string]string) error {
	stepType := actual.Type()
	if actual.Type() != expected["type"] {
		return fmt.Errorf("found step %s but types differ, actual: %s, expected: %s", expected["name"], actual.Type(), expected["type"])
	}

	if stepType == api.StepTypePowershell {
		dt := actual.(*api.StepPowershell)
		if p, ok := expected["file"]; ok {
			if p != dt.ScriptFile {
				return fmt.Errorf("property 'file' differs, actual: %s, expected: %s", dt.ScriptFile, p)
			}
		}

		if p, ok := expected["args"]; ok {
			if p != dt.ScriptArgs {
				return fmt.Errorf("property 'args' differs, actual: %s, expected: %s", dt.ScriptArgs, p)
			}
		}

		if p, ok := expected["code"]; ok {
			if p != dt.Code {
				return fmt.Errorf("property 'code' differs, actual: %s, expected: %s", dt.Code, p)
			}
		}
		return nil
	}

	if stepType == api.StepTypeCommandLine {
		dt := actual.(*api.StepCommandLine)
		if p, ok := expected["file"]; ok {
			if p != dt.CommandExecutable {
				return fmt.Errorf("property 'file' differs, actual: %s, expected: %s", dt.CommandExecutable, p)
			}
		}

		if p, ok := expected["args"]; ok {
			if p != dt.CommandParameters {
				return fmt.Errorf("property 'args' differs, actual: %s, expected: %s", dt.CommandParameters, p)
			}
		}

		if p, ok := expected["code"]; ok {
			if p != dt.CustomScript {
				return fmt.Errorf("property 'code' differs, actual: %s, expected: %s", dt.CustomScript, p)
			}
		}
		return nil
	}

	return fmt.Errorf("unexpected step type found: %s", stepType)
}

func testAccCheckVcsRootAttached(vcs *[]*api.VcsRootEntry, n string, co string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *vcs == nil {
			return errors.New("variable VcsRootEntries must not be nil")
		}

		for _, v := range *vcs {
			if v.VcsRoot.Name == n {
				if v.CheckoutRules == co {
					return nil
				}
			}
		}

		return fmt.Errorf("the VCS Root with name '%s' and checkout rules '%s' was not found", n, co)
	}
}

func testAccCheckBuildConfigExists(n string, out *api.BuildType) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		return buildConfigExistsHelper(n, s, client, out)
	}
}

func updateBuildCounter(buildType *api.BuildType, counter int) {
	client := testAccProvider.Meta().(*api.Client)
	id := buildType.ID

	bt, err := client.BuildTypes.GetByID(id)
	if err != nil {
		panic(err)
	}
	bt.Options.BuildCounter = counter

	_, err = client.BuildTypes.Update(bt)
	if err != nil {
		panic(err)
	}
}

func buildConfigExistsHelper(n string, s *terraform.State, client *api.Client, out *api.BuildType) error {
	rs, ok := s.RootModule().Resources[n]
	if !ok {
		return fmt.Errorf("resource not found: %s", n)
	}

	if rs.Primary.ID == "" {
		return fmt.Errorf("no id for %s is set", n)
	}

	resp, err := client.BuildTypes.GetByID(rs.Primary.ID)

	if err != nil {
		return fmt.Errorf("received an error retrieving Build Configurationt: %s", err)
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
			return fmt.Errorf("received an error retrieving the Build Configuration: %s", err)
		}

		return errors.New("build configuration still exists")
	}
	return nil
}

const TestAccBuildConfigBasic = `
resource "teamcity_project" "build_config_project_test" {
  name = "build_config_project_test"
}

resource "teamcity_build_config" "build_configuration_test" {
	name = "build config test"
	project_id = "${teamcity_project.build_config_project_test.id}"
	description = "build config test desc"
	settings {
		build_number_format = "2.0.%build.counter%"
	}
}
`

const TestAccBuildConfigTemplateWithDescription = `
resource "teamcity_project" "build_config_project_test" {
  name = "build_config_project_test"
}

resource "teamcity_build_config" "build_configuration_test" {
	name = "build config template"
	is_template = "true"
	project_id = "${teamcity_project.build_config_project_test.id}"
	description = "build config template desc"
}
`

const TestAccBuildConfigTemplate = `
resource "teamcity_project" "build_config_project_test" {
  name = "build_config_project_test"
}

resource "teamcity_build_config" "build_configuration_test" {
	name = "build config template"
	is_template = "true"
	project_id = "${teamcity_project.build_config_project_test.id}"
}
`

const TestAccBuildConfigBasicUpdated = `
resource "teamcity_project" "build_config_project_test" {
  name = "build_config_project_test"
}

resource "teamcity_build_config" "build_configuration_test" {
	name = "build config updated"
	project_id = "${teamcity_project.build_config_project_test.id}"
	description = "build config test desc updated"
	settings {
		build_number_format = "3.0.%build.counter%"
	}
}
`
const TestAccBuildConfigBuildCounter = `
resource "teamcity_project" "project" {
  name = "project"
}

resource "teamcity_build_config" "build_config" {
	name = "build config test"
	project_id = "${teamcity_project.project.id}"
	description = "build config test desc"
	settings {
		build_counter = 2
	}
}
`

const TestAccBuildConfigBuildCounterUpdated = `
resource "teamcity_project" "project" {
	name = "project"
  }

  resource "teamcity_build_config" "build_config" {
	  name = "build config test"
	  project_id = "${teamcity_project.project.id}"
	  description = "build config test desc"
	  settings {
		build_counter = 10
	  }
  }
`
const TestAccBuildConfigParams = `
resource "teamcity_project" "build_config_project_test" {
  name = "build_config_project_test"
}

resource "teamcity_build_config" "build_configuration_test" {
	name = "build config test"
	project_id = "${teamcity_project.build_config_project_test.id}"

	env_params = {
		DEPLOY_SERVER = "server.com"
		some_variable = "hello"
	}

	config_params = {
		"github.repository" = "nocode"
	}

	sys_params = {
		system_param = "system_value"
	}
}
`

const TestAccBuildConfigParamsUpdated = `
resource "teamcity_project" "build_config_project_test" {
  name = "build_config_project_test"
}

resource "teamcity_build_config" "build_configuration_test" {
	name = "build config test"
	project_id = "${teamcity_project.build_config_project_test.id}"

	env_params = {
		DEPLOY_SERVER = "server.com"
		some_variable = "hello"
	}

	config_params = {
		"github.repository" = "updated_repo"
	}
}
`

const TestAccBuildConfigSettings = `
resource "teamcity_project" "build_config_project_test" {
  name = "build_config_project_test"
}

resource "teamcity_build_config" "build_configuration_test" {
  name = "build config test"
  project_id = "${teamcity_project.build_config_project_test.id}"

  settings {
	configuration_type = "DEPLOYMENT"
    build_number_format = "1.0.%build.counter%"
    build_counter = 20
    allow_personal_builds = true
    artifact_paths = ["+:*.json => /config/*.json"]
    detect_hanging = true
    status_widget = false
    concurrent_limit = 10
  }
}
`

const TestAccBuildConfigSettingsUpdated = `
resource "teamcity_project" "build_config_project_test" {
  name = "build_config_project_test"
}

resource "teamcity_build_config" "build_configuration_test" {
  name = "build config test"
  project_id = "${teamcity_project.build_config_project_test.id}"

  settings {
	configuration_type = "REGULAR"
    build_number_format = "2.0.%build.counter%"
    build_counter = 25
    allow_personal_builds = false
    artifact_paths = ["+:*.json => /artifacts/*.json"]
    detect_hanging = false
    status_widget = true
    concurrent_limit = 0
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

const TestAccBuildConfigAttachTemplates = `
resource "teamcity_project" "build_config_project_test" {
  name = "build_config_project_test"
}

resource "teamcity_build_config" "build_configuration_test" {
	name = "build config test"
	project_id = teamcity_project.build_config_project_test.id

	templates = [ teamcity_build_config.build_configuration_template1.id, teamcity_build_config.build_configuration_template2.id ]

	depends_on = [
		teamcity_build_config.build_configuration_template1,
		teamcity_build_config.build_configuration_template2,
	]
}

resource "teamcity_build_config" "build_configuration_template1" {
	name = "build template 1"
	is_template = true
	project_id = teamcity_project.build_config_project_test.id
}

resource "teamcity_build_config" "build_configuration_template2" {
	name = "build template 2"
	is_template = true
	project_id = teamcity_project.build_config_project_test.id
}
`

const TestAccBuildConfigUpdatedTemplates = `
resource "teamcity_project" "build_config_project_test" {
  name = "build_config_project_test"
}

resource "teamcity_build_config" "build_configuration_test" {
	name = "build config test"
	project_id = teamcity_project.build_config_project_test.id

	templates = [ teamcity_build_config.build_configuration_template1.id, teamcity_build_config.build_configuration_template3.id ]
}

resource "teamcity_build_config" "build_configuration_template1" {
	name = "build template 1"
	is_template = true
	project_id = teamcity_project.build_config_project_test.id
}

# Need to keep this around here or Terraform is trying to destroy template2 before updating the dependant resource
resource "teamcity_build_config" "build_configuration_template2" {
	name = "build template 2"
	is_template = true
	project_id = teamcity_project.build_config_project_test.id
}

resource "teamcity_build_config" "build_configuration_template3" {
	name = "build template 3"
	is_template = true
	project_id = teamcity_project.build_config_project_test.id
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

const TestAccBuildConfigStepsPowershellUpdated = `
resource "teamcity_project" "build_config_project_test" {
  name = "build_config_project_test"
}

resource "teamcity_build_config" "build_configuration_test" {
	name = "build config test"
	project_id = "${teamcity_project.build_config_project_test.id}"

	step {
		type = "powershell"
		name = "updated_script"
		file = "updated.ps1"
		args = "-Target pullrequest"
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

const TestAccBuildConfigStepsCmdLineUpdated = `
resource "teamcity_project" "build_config_project_test" {
  name = "build_config_project_test"
}

resource "teamcity_build_config" "build_configuration_test" {
	name = "build config test"
	project_id = "${teamcity_project.build_config_project_test.id}"

	step {
		type = "cmd_line"
		name = "build_script"
		code = "echo \"Hello Foo\""
	}

	step {
		type = "cmd_line"
		name = "build_executable"
		file = "./build.sh"
		args = "default_target --verbose"
	}
}
`

const TestAccBuildConfigurationIdWithParent = `
resource "teamcity_project" "parent" {
	name = "parent"
}

resource "teamcity_project" "child" {
	name = "child"
	parent_id = "${teamcity_project.parent.id}"
}

resource teamcity_build_config "build_config" {
	name = "build config"
	project_id = "${teamcity_project.child.id}"
}
`

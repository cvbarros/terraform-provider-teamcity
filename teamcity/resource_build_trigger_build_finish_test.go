package teamcity_test

import (
	"testing"

	api "github.com/cvbarros/go-teamcity-sdk/teamcity"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTeamcityBuildTriggerBuildFinish_Basic(t *testing.T) {
	resName := "teamcity_build_trigger_build_finish.test"
	var out api.Trigger
	var bc, sc api.BuildType

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamcityBuildTriggerDestroy(&bc.ID, "teamcity_build_trigger_build_finish"),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccBuildTriggerBuildFinishBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists("teamcity_build_config.config", &bc),
					testAccCheckBuildConfigExists("teamcity_build_config.source", &sc),
					testAccCheckTeamcityBuildTriggerExists(resName, &bc.ID, &out, true),
					resource.TestCheckResourceAttrPtr(resName, "build_config_id", &bc.ID),
					resource.TestCheckResourceAttrPtr(resName, "source_build_config_id", &sc.ID),
					resource.TestCheckResourceAttr(resName, "after_successful_only", "true"),
					resource.TestCheckResourceAttr(resName, "branch_filter.0", "master"),
					resource.TestCheckResourceAttr(resName, "branch_filter.1", "feature"),
				),
			},
		},
	})
}

func TestAccTeamcityBuildTriggerBuildFinish_Update(t *testing.T) {
	resName := "teamcity_build_trigger_build_finish.test"
	var out api.Trigger
	var bc api.BuildType

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamcityBuildTriggerDestroy(&bc.ID, "teamcity_build_trigger_build_finish"),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccBuildTriggerBuildFinishBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists("teamcity_build_config.config", &bc),
					testAccCheckTeamcityBuildTriggerExists(resName, &bc.ID, &out, true),
					resource.TestCheckResourceAttrPtr(resName, "build_config_id", &bc.ID),
					resource.TestCheckResourceAttr(resName, "after_successful_only", "true"),
					resource.TestCheckResourceAttr(resName, "branch_filter.0", "master"),
					resource.TestCheckResourceAttr(resName, "branch_filter.1", "feature"),
				),
			},
			resource.TestStep{
				Config: TestAccBuildTriggerBuildFinishUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists("teamcity_build_config.config", &bc),
					testAccCheckTeamcityBuildTriggerExists(resName, &bc.ID, &out, true),
					resource.TestCheckResourceAttrPtr(resName, "build_config_id", &bc.ID),
					resource.TestCheckResourceAttr(resName, "after_successful_only", "false"),
					resource.TestCheckResourceAttr(resName, "branch_filter.0", "tag1"),
					resource.TestCheckResourceAttr(resName, "branch_filter.1", "tag2"),
				),
			},
		},
	})
}

const TestAccBuildTriggerBuildFinishBasic = `
resource "teamcity_project" "trigger_project_test" {
  name = "Trigger Build Finish Project"
}

resource "teamcity_build_config" "config" {
	name = "BuildConfig"
	project_id = "${teamcity_project.trigger_project_test.id}"
}

resource "teamcity_build_config" "source" {
	name = "SourceConfig"
	project_id = "${teamcity_project.trigger_project_test.id}"
}

resource "teamcity_build_trigger_build_finish" "test" {
	build_config_id = "${teamcity_build_config.config.id}"
	source_build_config_id = "${teamcity_build_config.source.id}"

	after_successful_only = true
	branch_filter = ["master", "feature"]
}
`

const TestAccBuildTriggerBuildFinishUpdated = `
resource "teamcity_project" "trigger_project_test" {
  name = "Trigger Build Finish Project"
}

resource "teamcity_build_config" "config" {
	name = "BuildConfig"
	project_id = "${teamcity_project.trigger_project_test.id}"
}

resource "teamcity_build_config" "source" {
	name = "SourceConfig"
	project_id = "${teamcity_project.trigger_project_test.id}"
}

resource "teamcity_build_trigger_build_finish" "test" {
	build_config_id = "${teamcity_build_config.config.id}"
	source_build_config_id = "${teamcity_build_config.source.id}"

	after_successful_only = false
	branch_filter = ["tag1", "tag2"]
}
`

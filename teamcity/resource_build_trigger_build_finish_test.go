package teamcity_test

import (
	"testing"

	api "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTeamcityBuildTriggerBuildFinish_Basic(t *testing.T) {
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

resource "teamcity_build_trigger_build_finish" "test" {
	build_config_id = "${teamcity_build_config.config.id}"

	after_successful_only = true
	branch_filter = ["master", "feature"]
}
`

package teamcity_test

import (
	"testing"

	api "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTeamcityBuildTriggerSchedule_Daily(t *testing.T) {
	resName := "teamcity_build_trigger_schedule.test"
	var out api.Trigger
	var bc api.BuildType

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamcityBuildTriggerDestroy(&bc.ID, "teamcity_build_trigger_schedule"),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccBuildTriggerScheduleDaily,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists("teamcity_build_config.config", &bc),
					testAccCheckTeamcityBuildTriggerExists(resName, &bc.ID, &out),
					resource.TestCheckResourceAttr(resName, "schedule", "daily"),
					resource.TestCheckNoResourceAttr(resName, "weekday"),
				),
			},
		},
	})
}

func TestAccTeamcityBuildTriggerSchedule_Weekly(t *testing.T) {
	resName := "teamcity_build_trigger_schedule.test"
	var out api.Trigger
	var bc api.BuildType

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamcityBuildTriggerDestroy(&bc.ID, "teamcity_build_trigger_schedule"),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccBuildTriggerScheduleWeekly,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists("teamcity_build_config.config", &bc),
					testAccCheckTeamcityBuildTriggerExists(resName, &bc.ID, &out),
					resource.TestCheckResourceAttr(resName, "schedule", "weekly"),
					resource.TestCheckResourceAttr(resName, "weekday", "Saturday"),
				),
			},
		},
	})
}

const TestAccBuildTriggerScheduleDaily = `
resource "teamcity_project" "trigger_project_test" {
  name = "Trigger Build Finish Project"
}

resource "teamcity_build_config" "config" {
	name = "BuildConfig"
	project_id = "${teamcity_project.trigger_project_test.id}"
}

resource "teamcity_build_trigger_schedule" "test" {
    build_config_id = "${teamcity_build_config.config.id}"

    schedule = "daily"
    timezone = "America/Sao Paulo"
    hour = 12
    minute = 37
    rules = ["+:*", "-:*.md"]
}
`

const TestAccBuildTriggerScheduleWeekly = `
resource "teamcity_project" "trigger_project_test" {
  name = "Trigger Build Finish Project"
}

resource "teamcity_build_config" "config" {
	name = "BuildConfig"
	project_id = "${teamcity_project.trigger_project_test.id}"
}

resource "teamcity_build_trigger_schedule" "test" {
    build_config_id = "${teamcity_build_config.config.id}"

    schedule = "weekly"
    timezone = "America/Sao Paulo"
    hour = 12
    minute = 37
    weekday = "Saturday"
    rules = ["+:*", "-:*.md"]
}
`

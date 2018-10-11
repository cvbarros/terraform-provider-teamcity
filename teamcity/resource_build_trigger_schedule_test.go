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
					testAccCheckTeamcityBuildTriggerExists(resName, &bc.ID, &out, true),
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
					testAccCheckTeamcityBuildTriggerExists(resName, &bc.ID, &out, true),
					resource.TestCheckResourceAttr(resName, "schedule", "weekly"),
					resource.TestCheckResourceAttr(resName, "weekday", "Saturday"),
				),
			},
		},
	})
}

func TestAccTeamcityBuildTriggerSchedule_Options(t *testing.T) {
	resName := "teamcity_build_trigger_schedule.test"
	var out api.Trigger
	var bc api.BuildType
	var watched api.BuildType

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamcityBuildTriggerDestroy(&bc.ID, "teamcity_build_trigger_schedule"),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccBuildTriggerScheduleOptions,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists("teamcity_build_config.config", &bc),
					testAccCheckBuildConfigExists("teamcity_build_config.watched", &watched),
					testAccCheckTeamcityBuildTriggerExists(resName, &bc.ID, &out, true),
					resource.TestCheckResourceAttr(resName, "schedule", "daily"),
					resource.TestCheckResourceAttr(resName, "queue_optimization", "true"),
					resource.TestCheckResourceAttr(resName, "on_all_compatible_agents", "true"),
					resource.TestCheckResourceAttr(resName, "with_pending_changes_only", "true"),
					resource.TestCheckResourceAttr(resName, "promote_watched_build", "false"),
					resource.TestCheckResourceAttr(resName, "enforce_clean_checkout", "true"),
					resource.TestCheckResourceAttr(resName, "enforce_clean_checkout_dependencies", "true"),
					resource.TestCheckResourceAttr(resName, "only_if_watched_changes", "true"),
					resource.TestCheckResourceAttrPtr(resName, "watched_build_config_id", &watched.ID),
					resource.TestCheckResourceAttr(resName, "revision", "lastFinished"),
					resource.TestCheckResourceAttr(resName, "only_if_watched_changes", "true"),
					resource.TestCheckResourceAttr(resName, "watched_branch", "unstable"),
				),
			},
		},
	})
}

func TestAccTeamcityBuildTriggerSchedule_DefaultOptions(t *testing.T) {
	resName := "teamcity_build_trigger_schedule.test"
	var out api.Trigger
	var bc api.BuildType
	var watched api.BuildType

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamcityBuildTriggerDestroy(&bc.ID, "teamcity_build_trigger_schedule"),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccBuildTriggerScheduleDefaultOptions,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists("teamcity_build_config.config", &bc),
					testAccCheckBuildConfigExists("teamcity_build_config.watched", &watched),
					testAccCheckTeamcityBuildTriggerExists(resName, &bc.ID, &out, true),
					resource.TestCheckResourceAttr(resName, "schedule", "daily"),
					resource.TestCheckResourceAttr(resName, "queue_optimization", "true"),
					resource.TestCheckResourceAttr(resName, "with_pending_changes_only", "true"),
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

const TestAccBuildTriggerScheduleOptions = `
resource "teamcity_project" "trigger_project_test" {
  name = "Trigger Build Finish Project"
}

resource "teamcity_build_config" "config" {
	name = "BuildConfig"
	project_id = "${teamcity_project.trigger_project_test.id}"
}

resource "teamcity_build_config" "watched" {
	name = "WatchedBuild"
	project_id = "${teamcity_project.trigger_project_test.id}"
}

resource "teamcity_build_trigger_schedule" "test" {
    build_config_id = "${teamcity_build_config.config.id}"

    schedule = "daily"
    timezone = "America/Sao Paulo"
    hour = 12
    minute = 37
	rules = ["+:*", "-:*.md"]

	queue_optimization = true
	on_all_compatible_agents = true
	with_pending_changes_only = true
	promote_watched_build = false

	enforce_clean_checkout = true
	enforce_clean_checkout_dependencies = true

	only_if_watched_changes = true
	watched_build_config_id = "${teamcity_build_config.watched.id}"
	revision = "lastFinished"
	watched_branch = "unstable"
}
`

const TestAccBuildTriggerScheduleDefaultOptions = `
resource "teamcity_project" "trigger_project_test" {
  name = "Trigger Build Finish Project"
}

resource "teamcity_build_config" "config" {
	name = "BuildConfig"
	project_id = "${teamcity_project.trigger_project_test.id}"
}

resource "teamcity_build_config" "watched" {
	name = "WatchedBuild"
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

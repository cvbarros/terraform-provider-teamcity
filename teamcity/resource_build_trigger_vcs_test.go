package teamcity_test

import (
	"fmt"
	"strings"
	"testing"

	api "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTeamcityBuildTriggerVcs_Basic(t *testing.T) {
	resName := "teamcity_build_trigger_vcs.test"
	var out api.Trigger
	var bc api.BuildType

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamcityBuildTriggerDestroy(&bc.ID, "teamcity_build_trigger_vcs"),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccBuildTriggerVcsBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists("teamcity_build_config.config", &bc),
					testAccCheckTeamcityBuildTriggerExists(resName, &bc.ID, &out),
				),
			},
		},
	})
}

func testAccCheckTeamcityBuildTriggerDestroy(bt *string, resourceType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		return buildTriggerDestroyHelper(s, bt, client, resourceType)
	}
}

func buildTriggerDestroyHelper(s *terraform.State, bt *string, client *api.Client, resourceType string) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != resourceType {
			continue
		}

		ts := client.TriggerService(*bt)
		_, err := ts.GetByID(r.Primary.ID)

		if err != nil {
			if strings.Contains(err.Error(), "404") {
				continue
			}
			return fmt.Errorf("Received an error retrieving the Trigger: %s", err)
		}

		return fmt.Errorf("Trigger still exists")
	}
	return nil
}

func testAccCheckTeamcityBuildTriggerExists(n string, bt *string, snap *api.Trigger) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		return teamcityBuildTriggerExistsHelper(n, bt, s, client, snap)
	}
}

func teamcityBuildTriggerExistsHelper(n string, bt *string, s *terraform.State, client *api.Client, snap *api.Trigger) error {
	rs, ok := s.RootModule().Resources[n]
	if !ok {
		return fmt.Errorf("Not found: %s", n)
	}

	if rs.Primary.ID == "" {
		return fmt.Errorf("No ID is set")
	}

	ts := client.TriggerService(*bt)
	out, err := ts.GetByID(rs.Primary.ID)
	if err != nil {
		return fmt.Errorf("Received an error retrieving Trigger: %s", err)
	}

	*snap = out
	return nil
}

const TestAccBuildTriggerVcsBasic = `
resource "teamcity_project" "trigger_project_test" {
  name = "Trigger Project"
}

resource "teamcity_build_config" "config" {
	name = "BuildConfig"
	project_id = "${teamcity_project.trigger_project_test.id}"
}

resource "teamcity_build_trigger_vcs" "test" {
	build_config_id = "${teamcity_build_config.config.id}"
	rules = "+:*"
	branch_filter = "+:pull/*"
}
`

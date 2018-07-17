package teamcity_test

import (
	"fmt"
	"strings"
	"testing"

	api "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTeamcityTrigger_Basic(t *testing.T) {
	resName := "teamcity_trigger.test"
	var out api.Trigger
	var bc api.BuildType

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckteamcityTriggerDestroy(&out.BuildTypeID),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccTriggerBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists("teamcity_build_config.config", &bc),
					testAccCheckteamcityTriggerExists(resName, &bc.ID, &out),
					resource.TestCheckResourceAttrPtr(resName, "build_config_id", &out.BuildTypeID),
				),
			},
		},
	})
}

func testAccCheckteamcityTriggerDestroy(bt *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		return triggerDestroyHelper(s, bt, client)
	}
}

func triggerDestroyHelper(s *terraform.State, bt *string, client *api.Client) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != "teamcity_trigger" {
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

func testAccCheckteamcityTriggerExists(n string, bt *string, snap *api.Trigger) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		return teamcityTriggerExistsHelper(n, bt, s, client, snap)
	}
}

func teamcityTriggerExistsHelper(n string, bt *string, s *terraform.State, client *api.Client, snap *api.Trigger) error {
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

	*snap = *out
	return nil
}

const TestAccTriggerBasic = `
resource "teamcity_project" "trigger_project_test" {
  name = "Snapshot"
}

resource "teamcity_build_config" "config" {
	name = "BuildConfig"
	project_id = "${teamcity_project.trigger_project_test.id}"
}

resource "teamcity_trigger" "test" {
	build_config_id = "${teamcity_build_config.config.id}"
	rules = "+:*"
	branch_filter = "+:pull/*"
}
`

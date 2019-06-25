package teamcity_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	api "github.com/cvbarros/go-teamcity-sdk/teamcity"
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
					testAccCheckTeamcityBuildTriggerExists(resName, &bc.ID, &out, true),
					resource.TestCheckResourceAttr(resName, "rules.0", "+:*"),
					resource.TestCheckResourceAttr(resName, "branch_filter.0", "+:pull/*"),
				),
			},
		},
	})
}

func TestAccTeamcityBuildTriggerVcs_Update(t *testing.T) {
	resName := "teamcity_build_trigger_vcs.test"
	var before, after api.Trigger
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
					testAccCheckTeamcityBuildTriggerExists(resName, &bc.ID, &before, true),
				),
			},
			resource.TestStep{
				Config: TestAccBuildTriggerVcsUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists("teamcity_build_config.config", &bc),
					testAccCheckTeamcityBuildTriggerExists(resName, &bc.ID, &after, true),
					resource.TestCheckResourceAttr(resName, "rules.0", "updated_rules"),
					resource.TestCheckResourceAttr(resName, "branch_filter.0", "+:refs/head/master"),
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

func testAccCheckTeamcityBuildTriggerRemoved(buildTypeId *string, t *api.Trigger) resource.TestCheckFunc {
	return func(S *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)

		_, err := client.TriggerService(*buildTypeId).GetByID((*t).ID())
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				return nil //It's ok, removed
			}
		}

		return fmt.Errorf("expected resource with id: %s to be removed, but it wasn't", (*t).ID())
	}
}

func testAccCheckTeamcityBuildTriggerExists(n string, bt *string, t *api.Trigger, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)

		found, err := teamcityBuildTriggerExistsHelper(n, bt, s, client, t)
		if !exists {
			if found {
				return errors.New("expected trigger to be removed, but still exists")
			}
		}
		return err
	}
}

func teamcityBuildTriggerExistsHelper(n string, bt *string, s *terraform.State, client *api.Client, t *api.Trigger) (bool, error) {
	rs, ok := s.RootModule().Resources[n]
	if !ok {
		return false, fmt.Errorf("Not found: %s", n)
	}

	if rs.Primary.ID == "" {
		return false, fmt.Errorf("No ID is set")
	}

	ts := client.TriggerService(*bt)
	out, err := ts.GetByID(rs.Primary.ID)
	if err != nil {
		return false, fmt.Errorf("Received an error retrieving Trigger: %s", err)
	}

	*t = out
	return true, nil
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
	rules = ["+:*"]
	branch_filter = ["+:pull/*"]
}
`

const TestAccBuildTriggerVcsUpdated = `
resource "teamcity_project" "trigger_project_test" {
  name = "Trigger Project"
}

resource "teamcity_build_config" "config" {
	name = "BuildConfig"
	project_id = "${teamcity_project.trigger_project_test.id}"
}

resource "teamcity_build_trigger_vcs" "test" {
	build_config_id = "${teamcity_build_config.config.id}"
	rules = ["updated_rules"]
	branch_filter = ["+:refs/head/master"]
}
`

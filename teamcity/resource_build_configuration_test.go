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
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildConfigDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccBuildConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists("teamcity_build_config.build_configuration_test"),
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

// func TestAccBuildConfig_Delete(t *testing.T) {
// 	resName := "teamcity_build_config.build_configuration_test"

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckBuildConfigDestroy,
// 		Steps: []resource.TestStep{
// 			resource.TestStep{
// 				Config: TestAccBuildConfigConfig("_Root"),
// 				Check:  resource.TestCheckResourceAttr(resName, "project_id", "_Root"),
// 			},
// 		},
// 	})
// }

func testAccCheckBuildConfigExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		return buildConfigExistsHelper(s, client)
	}
}

func buildConfigExistsHelper(s *terraform.State, client *api.Client) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != "teamcity_build_config" {
			continue
		}

		if _, err := client.BuildTypes.GetById(r.Primary.ID); err != nil {
			return fmt.Errorf("Received an error retrieving Build Configurationt: %s", err)
		}
	}

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

		_, err := client.BuildTypes.GetById(r.Primary.ID)

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

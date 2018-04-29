package teamcity_test

import (
	"fmt"
	"strings"
	"testing"

	api "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTeamcityProject_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamcityProjectDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccTeamcityProjectConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamcityProjectExists("teamcity_project.testproj"),
					resource.TestCheckResourceAttr(
						"teamcity_project.testproj", "name", "testproj",
					),
				),
			},
		},
	})
}

func TestAccTeamcityProject_UpdateName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamcityProjectDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccTeamcityProjectConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamcityProjectExists("teamcity_project.testproj"),
					resource.TestCheckResourceAttr(
						"teamcity_project.testproj", "name", "testproj",
					),
				),
			},
			resource.TestStep{
				Config: testAccTeamcityProjectConfigUpdatedName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamcityProjectExists("teamcity_project.testproj"),
					resource.TestCheckResourceAttr(
						"teamcity_project.testproj", "name", "testproj_updated",
					),
				),
			},
		},
	})
}

func testAccCheckTeamcityProjectExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		return teamcityProjectExistsHelper(s, client)
	}
}

func teamcityProjectExistsHelper(s *terraform.State, client *api.Client) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != "teamcity_project" {
			continue
		}

		if _, err := client.Projects.GetById(r.Primary.ID); err != nil {
			return fmt.Errorf("Received an error retrieving project: %s", err)
		}
	}

	return nil
}

func testAccCheckTeamcityProjectDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	return teamcityProjectDestroyHelper(s, client)
}

func teamcityProjectDestroyHelper(s *terraform.State, client *api.Client) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != "teamcity_project" {
			continue
		}

		_, err := client.Projects.GetById(r.Primary.ID)

		if err != nil {
			if strings.Contains(err.Error(), "404") {
				continue
			}
			return fmt.Errorf("Received an error retrieving the project: %s", err)
		}

		return fmt.Errorf("Project still exists")
	}
	return nil
}

const testAccTeamcityProjectConfig = `
resource "teamcity_project" "testproj" {
  name = "testproj"
}
`

const testAccTeamcityProjectConfigUpdatedName = `
resource "teamcity_project" "testproj" {
	name = "testproj_updated"
}
`

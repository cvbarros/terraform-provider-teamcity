package teamcity

import (
	"testing"

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
					resource.TestCheckResourceAttr(
						"teamcity_project.canary", "name", "canary",
					),
				),
			},
		},
	})
}

func testAccCheckTeamcityProjectDestroy(s *terraform.State) error {
	return nil
}

const testAccTeamcityProjectConfig = `
resource "teamcity_project" "canary" {
  name = "canary"
}
`

package teamcity_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceAgentPool_Basic(t *testing.T) {
	resName := "data.teamcity_agent_pool.project"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAgentPool,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", "Default"),
					resource.TestCheckResourceAttr(resName, "max_agents", "-1"),
				),
			},
		},
	})
}

const testAccDataSourceAgentPool = `
data "teamcity_agent_pool" "test" {
  name = "Default"
}
`

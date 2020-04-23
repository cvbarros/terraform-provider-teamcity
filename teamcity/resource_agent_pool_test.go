package teamcity_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	api "github.com/cvbarros/go-teamcity/teamcity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTeamCityAgentPool_Basic(t *testing.T) {
	resName := "teamcity_agent_pool.test"
	name := fmt.Sprintf("Pool %d", time.Now().Unix())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamCityAgentPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamCityAgentPoolBasicConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamCityAgentPoolExists(resName),
					resource.TestCheckResourceAttr(resName, "max_agents", "-1"),
					resource.TestCheckResourceAttr(resName, "name", name),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTeamCityAgentPool_Complete(t *testing.T) {
	resName := "teamcity_agent_pool.test"
	name := fmt.Sprintf("Pool %d", time.Now().Unix())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamCityAgentPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamCityAgentPoolCompleteConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamCityAgentPoolExists(resName),
					resource.TestCheckResourceAttr(resName, "max_agents", "5"),
					resource.TestCheckResourceAttr(resName, "name", name),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTeamCityAgentPoolExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = client.AgentPools.GetByID(id)
		if err != nil {
			return fmt.Errorf("Received an error retrieving agent pool: %s", err)
		}

		return nil
	}
}

func testAccCheckTeamCityAgentPoolDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "teamcity_agent_pool" {
			continue
		}

		id, err := strconv.Atoi(r.Primary.ID)
		if err != nil {
			return err
		}

		_, err = client.AgentPools.GetByID(id)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				continue
			}
			return fmt.Errorf("Received an error retrieving the agent pool: %s", err)
		}

		return fmt.Errorf("Agent Pool still exists")
	}
	return nil
}

func testAccTeamCityAgentPoolBasicConfig(name string) string {
	return fmt.Sprintf(`
resource "teamcity_agent_pool" "test" {
  name = %s
}
`, name)
}

func testAccTeamCityAgentPoolCompleteConfig(name string) string {
	return fmt.Sprintf(`
resource "teamcity_agent_pool" "test" {
  name       = %s
  max_agents = 5
}
`, name)
}

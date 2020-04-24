package teamcity_test

import (
	"fmt"
	"strings"
	"testing"

	api "github.com/cvbarros/go-teamcity/teamcity"
	"github.com/cvbarros/terraform-provider-teamcity/teamcity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTeamcityFeatureGolang_Basic(t *testing.T) {
	resName := "teamcity_feature_golang.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildFeatureGolangDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccBuildFeatureGolang_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildFeatureGolangExists(resName),
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

func testAccCheckBuildFeatureGolangDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "teamcity_feature_github" {
			continue
		}

		id, err := teamcity.ParseFeatureGolangID(rs.Primary.ID)
		if err != nil {
			return err
		}

		srv := client.BuildFeatureService(id.BuildConfigID)
		if _, err := srv.GetByID(id.FeatureID); err != nil {
			if strings.Contains(err.Error(), "404") {
				continue
			}

			return fmt.Errorf("Received an error retrieving the Golang Build Feature: %s", err)
		}

		return fmt.Errorf("Golang Build Feature still exists")
	}
	return nil
}

func testAccCheckBuildFeatureGolangExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := teamcity.ParseFeatureGolangID(rs.Primary.ID)
		if err != nil {
			return err
		}

		srv := client.BuildFeatureService(id.BuildConfigID)
		if _, err := srv.GetByID(id.FeatureID); err != nil {
			return fmt.Errorf("Received an error retrieving Golang Build Feature: %s", err)
		}

		return nil
	}
}

const TestAccBuildFeatureGolang_basic = `
resource "teamcity_project" "test" {
  name = "Build Feature"
}

resource "teamcity_build_config" "test" {
  name = "BuildConfig"
  project_id = teamcity_project.test.id
}

resource "teamcity_feature_golang" "test" {
  build_config_id = teamcity_build_config.test.id
}
`

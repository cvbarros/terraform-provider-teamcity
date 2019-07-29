package teamcity_test

import (
	"fmt"
	"strings"
	"testing"

	api "github.com/cvbarros/go-teamcity-sdk/teamcity"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTeamcityFeatureCommitStatusPublisher_Github(t *testing.T) {
	resName := "teamcity_feature_commit_status_publisher.test"
	var out api.BuildFeature
	var bc api.BuildType

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildFeatureDestroy(&bc.ID),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccBuildFeatureCommitStatusPublisher_GithubPassword,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists("teamcity_build_config.config", &bc),
					testAccCheckBuildFeatureExists(resName, &bc.ID, &out),
					resource.TestCheckResourceAttr(resName, "publisher", "github"),
					resource.TestCheckResourceAttr(resName, "github.3735060251.auth_type", "password"),
					resource.TestCheckResourceAttr(resName, "github.3735060251.host", "https://api.github.com"),
					resource.TestCheckResourceAttr(resName, "github.3735060251.username", "bob"),
				),
			},
		},
	})
}

func TestAccTeamcityFeatureCommitStatusPublisher_GithubUpdate(t *testing.T) {
	resName := "teamcity_feature_commit_status_publisher.test"
	var out api.BuildFeature
	var bc api.BuildType

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildFeatureDestroy(&bc.ID),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccBuildFeatureCommitStatusPublisher_GithubPassword,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists("teamcity_build_config.config", &bc),
					testAccCheckBuildFeatureExists(resName, &bc.ID, &out),
					resource.TestCheckResourceAttr(resName, "publisher", "github"),
					resource.TestCheckResourceAttr(resName, "github.3735060251.auth_type", "password"),
					resource.TestCheckResourceAttr(resName, "github.3735060251.host", "https://api.github.com"),
					resource.TestCheckResourceAttr(resName, "github.3735060251.username", "bob"),
				),
			},
			resource.TestStep{
				Config: TestAccBuildFeatureCommitStatusPublisher_GithubPasswordUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists("teamcity_build_config.config", &bc),
					testAccCheckBuildFeatureExists(resName, &bc.ID, &out),
					resource.TestCheckResourceAttr(resName, "publisher", "github"),
					resource.TestCheckResourceAttr(resName, "github.3764292600.host", "https://api.github.com/v3"),
					resource.TestCheckResourceAttr(resName, "github.3764292600.username", "bob_updated"),
				),
			},
		},
	})
}

func testAccCheckBuildFeatureDestroy(bt *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		return buildFeatureDestroyHelper(s, bt, client)
	}
}

func buildFeatureDestroyHelper(s *terraform.State, bt *string, client *api.Client) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != "teamcity_feature_commit_status_publisher" {
			continue
		}

		srv := client.BuildFeatureService(*bt)
		_, err := srv.GetByID(r.Primary.ID)

		if err != nil {
			if strings.Contains(err.Error(), "404") {
				continue
			}
			return fmt.Errorf("Received an error retrieving the BuildFeature: %s", err)
		}

		return fmt.Errorf("BuildFeature still exists")
	}
	return nil
}

func testAccCheckBuildFeatureExists(n string, bt *string, out *api.BuildFeature) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		return teamcityBuildFeatureExistsHelper(n, bt, s, client, out)
	}
}

func teamcityBuildFeatureExistsHelper(n string, bt *string, s *terraform.State, client *api.Client, out *api.BuildFeature) error {
	rs, ok := s.RootModule().Resources[n]
	if !ok {
		return fmt.Errorf("Not found: %s", n)
	}

	if rs.Primary.ID == "" {
		return fmt.Errorf("No ID is set")
	}

	srv := client.BuildFeatureService(*bt)
	dt, err := srv.GetByID(rs.Primary.ID)
	if err != nil {
		return fmt.Errorf("Received an error retrieving BuildFeature: %s", err)
	}

	*out = dt
	return nil
}

const TestAccBuildFeatureCommitStatusPublisher_GithubPassword = `
resource "teamcity_project" "build_feature_project_test" {
  name = "Build Feature"
}

resource "teamcity_build_config" "config" {
	name = "BuildConfig"
	project_id = "${teamcity_project.build_feature_project_test.id}"
}

resource "teamcity_feature_commit_status_publisher" "test" {
	build_config_id = "${teamcity_build_config.config.id}"
	publisher = "github"
	github {
		auth_type = "password"
		host = "https://api.github.com"
		username = "bob"
		password = "1234"
	}
}
`

const TestAccBuildFeatureCommitStatusPublisher_GithubPasswordUpdated = `
resource "teamcity_project" "build_feature_project_test" {
  name = "Build Feature"
}

resource "teamcity_build_config" "config" {
	name = "BuildConfig"
	project_id = "${teamcity_project.build_feature_project_test.id}"
}

resource "teamcity_feature_commit_status_publisher" "test" {
	build_config_id = "${teamcity_build_config.config.id}"
	publisher = "github"
	github {
		auth_type = "password"
		host = "https://api.github.com/v3"
		username = "bob_updated"
		password = "1234"
	}
}
`

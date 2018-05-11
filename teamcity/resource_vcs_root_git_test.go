package teamcity_test

import (
	"fmt"
	"strings"
	"testing"

	api "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccVcsRootGit_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVcsRootGitDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVcsRootGitBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVcsRootGitExists("teamcity_vcs_root_git.vcs_root_git_test"),
					resource.TestCheckResourceAttr(
						"teamcity_vcs_root_git.vcs_root_git_test", "name", "application",
					),
				),
			},
		},
	})
}

func testAccCheckVcsRootGitExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		return vcsRootGitExistsHelper(s, client)
	}
}

func vcsRootGitExistsHelper(s *terraform.State, client *api.Client) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != "teamcity_vcs_root_git" {
			continue
		}

		if _, err := client.VcsRoots.GetById(r.Primary.ID); err != nil {
			return fmt.Errorf("Received an error retrieving VCS Root: %s", err)
		}
	}

	return nil
}

func testAccCheckVcsRootGitDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	return vcsRootGitDestroyHelper(s, client)
}

func vcsRootGitDestroyHelper(s *terraform.State, client *api.Client) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != "teamcity_vcs_root_git" {
			continue
		}

		_, err := client.VcsRoots.GetById(r.Primary.ID)

		if err != nil {
			if strings.Contains(err.Error(), "404") {
				continue
			}
			return fmt.Errorf("Received an error retrieving the VCS Root: %s", err)
		}

		return fmt.Errorf("VCS Root still exists")
	}
	return nil
}

const testAccVcsRootGitBasic = `
resource "teamcity_project" "vcs_root_project" {
  name = "vcs_root_project"
}

resource "teamcity_vcs_root_git" "vcs_root_git_test" {
	name = "application"
	project_id = "${teamcity_project.vcs_root_project.id}"
	repo_url = "https://github.com/kelseyhightower/nocode"
	default_branch = "refs/head/master"
}
`

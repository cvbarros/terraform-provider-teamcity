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
	var vcs api.GitVcsRoot
	resName := "teamcity_vcs_root_git.git_test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVcsRootGitDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVcsRootGitBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVcsRootGitExists(resName, &vcs),
					resource.TestCheckResourceAttr(resName, "name", "application"),
					resource.TestCheckResourceAttr(resName, "fetch_url", "https://github.com/kelseyhightower/nocode"),
					resource.TestCheckResourceAttr(resName, "default_branch", "refs/head/master"),
					resource.TestCheckResourceAttr(resName, "project_id", "VcsRootProject"),
					resource.TestCheckResourceAttr(resName, "enable_branch_spec_tags", "true"),
					resource.TestCheckResourceAttr(resName, "submodule_checkout", "CHECKOUT"),
					resource.TestCheckResourceAttr(resName, "username_style", "userid"),
				),
			},
		},
	})
}

func TestAccVcsRootGit_UserpassAuth(t *testing.T) {
	var vcs api.GitVcsRoot
	resourceName := "teamcity_vcs_root_git.git_test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVcsRootGitDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVcsRootGitUserpass,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVcsRootGitExists(resourceName, &vcs),
					resource.TestCheckResourceAttr(resourceName, "auth.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "auth.2360613679.type", "userpass"),
					resource.TestCheckResourceAttr(resourceName, "auth.2360613679.username", "admin"),
					resource.TestCheckResourceAttrSet(resourceName, "auth.2360613679.password"),
				),
			},
		},
	})
}

func TestAccVcsRootGit_SshUploadedKeyAuth(t *testing.T) {
	var vcs api.GitVcsRoot
	resourceName := "teamcity_vcs_root_git.git_test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVcsRootGitDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVcsRootGitSSHUploadedKey,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVcsRootGitExists(resourceName, &vcs),
					resource.TestCheckResourceAttr(resourceName, "auth.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "auth.3867327421.type", "ssh"),
					resource.TestCheckResourceAttr(resourceName, "auth.3867327421.username", "admin"),
					resource.TestCheckResourceAttr(resourceName, "auth.3867327421.ssh_type", "uploadedKey"),
					resource.TestCheckResourceAttr(resourceName, "auth.3867327421.key_spec", "myKey"),
				),
			},
		},
	})
}

func TestAccVcsRootGit_Delete(t *testing.T) {
	resName := "teamcity_vcs_root_git.git_test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVcsRootGitDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVcsRootGitConfig("_Root"),
				Check:  resource.TestCheckResourceAttr(resName, "project_id", "_Root"),
			},
		},
	})
}

func testAccCheckVcsRootGitExists(name string, out *api.GitVcsRoot) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		return vcsRootGitExistsHelper(s, client, out)
	}
}

func vcsRootGitExistsHelper(s *terraform.State, client *api.Client, out *api.GitVcsRoot) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != "teamcity_vcs_root_git" {
			continue
		}

		resp, err := client.VcsRoots.GetByID(r.Primary.ID)
		if err != nil {
			return fmt.Errorf("Received an error retrieving VCS Root: %s", err)
		}

		*out = *resp.(*api.GitVcsRoot)
		return nil
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

		_, err := client.VcsRoots.GetByID(r.Primary.ID)

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

func testAccVcsRootGitConfig(projectId string) string {
	return fmt.Sprintf(`
resource "teamcity_vcs_root_git" "git_test" {
	name = "application"
	project_id = "%s"
	fetch_url = "https://github.com/kelseyhightower/nocode"
	default_branch = "refs/head/master"
}
`, projectId)
}

const testAccVcsRootGitBasic = `
resource "teamcity_project" "vcs_root_project" {
  name = "vcs_root_project"
}

resource "teamcity_vcs_root_git" "git_test" {
	name = "application"
	project_id = "${teamcity_project.vcs_root_project.id}"
	fetch_url = "https://github.com/kelseyhightower/nocode"
	default_branch = "refs/head/master"
	username_style = "userid"
	submodule_checkout = "checkout"
	enable_branch_spec_tags = true
}
`

const testAccVcsRootGitUserpass = `
resource "teamcity_project" "vcs_root_project" {
  name = "vcs_root_project"
}

resource "teamcity_vcs_root_git" "git_test" {
	name = "application"
	project_id = "${teamcity_project.vcs_root_project.id}"
	fetch_url = "https://github.com/cvbarros/terraform-provider-teamcity"
	default_branch = "refs/head/master"

	auth {
		type = "userpass"
		username = "admin"
		password = "admin"
	}
}
`

const testAccVcsRootGitSSHUploadedKey = `
resource "teamcity_project" "vcs_root_project" {
  name = "vcs_root_project"
}

resource "teamcity_vcs_root_git" "git_test" {
	name = "application"
	project_id = "${teamcity_project.vcs_root_project.id}"
	fetch_url = "https://github.com/cvbarros/terraform-provider-teamcity"
	default_branch = "refs/head/master"

	auth {
		type = "ssh"
		username = "admin"
		ssh_type = "uploadedKey"
		key_spec = "myKey"
		password = "key_password"
	}
}
`

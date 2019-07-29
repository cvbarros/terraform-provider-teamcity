package teamcity_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	api "github.com/cvbarros/go-teamcity-sdk/teamcity"
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
					resource.TestCheckResourceAttr(resName, "fetch_url", "https://github.com/cvbarros/terraform-provider-teamcity"),
					resource.TestCheckResourceAttr(resName, "default_branch", "refs/head/master"),
					resource.TestCheckResourceAttr(resName, "branches.#", "2"),
					resource.TestCheckResourceAttr(resName, "branches.0", "+:refs/(pull/*)/head"),
					resource.TestCheckResourceAttr(resName, "branches.1", "+:refs/heads/develop"),
					resource.TestCheckResourceAttr(resName, "project_id", "VcsRootProject"),
					resource.TestCheckResourceAttr(resName, "enable_branch_spec_tags", "true"),
					resource.TestCheckResourceAttr(resName, "submodule_checkout", "CHECKOUT"),
					resource.TestCheckResourceAttr(resName, "username_style", "userid"),
					resource.TestCheckResourceAttr(resName, "modification_check_interval", "60"),
				),
			},
		},
	})
}

func TestAccVcsRootGit_UpdateBasic(t *testing.T) {
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
					resource.TestCheckResourceAttr(resName, "fetch_url", "https://github.com/cvbarros/terraform-provider-teamcity"),
					resource.TestCheckResourceAttr(resName, "default_branch", "refs/head/master"),
					resource.TestCheckResourceAttr(resName, "branches.#", "2"),
					resource.TestCheckResourceAttr(resName, "branches.0", "+:refs/(pull/*)/head"),
					resource.TestCheckResourceAttr(resName, "branches.1", "+:refs/heads/develop"),
					resource.TestCheckResourceAttr(resName, "project_id", "VcsRootProject"),
					resource.TestCheckResourceAttr(resName, "enable_branch_spec_tags", "true"),
					resource.TestCheckResourceAttr(resName, "submodule_checkout", "CHECKOUT"),
					resource.TestCheckResourceAttr(resName, "username_style", "userid"),
					resource.TestCheckResourceAttr(resName, "modification_check_interval", "60"),
				),
			},
			resource.TestStep{
				Config: testAccVcsRootGitUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVcsRootGitExists(resName, &vcs),
					resource.TestCheckResourceAttr(resName, "name", "application_updated"),
					resource.TestCheckResourceAttr(resName, "fetch_url", "https://github.com/cvbarros/go-teamcity-sdk"),
					resource.TestCheckResourceAttr(resName, "default_branch", "refs/head/develop"),
					resource.TestCheckResourceAttr(resName, "branches.#", "1"),
					resource.TestCheckResourceAttr(resName, "branches.0", "+:refs/heads/master"),
					resource.TestCheckResourceAttr(resName, "project_id", "VcsRootProjectNew"),
					resource.TestCheckResourceAttr(resName, "enable_branch_spec_tags", "false"),
					resource.TestCheckResourceAttr(resName, "submodule_checkout", "IGNORE"),
					resource.TestCheckResourceAttr(resName, "username_style", "author_name"),
					resource.TestCheckResourceAttr(resName, "modification_check_interval", "180"),
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

func TestAccVcsRootGit_AgentSettings(t *testing.T) {
	var vcs api.GitVcsRoot
	resourceName := "teamcity_vcs_root_git.git_test"
	expected := map[string]string{
		"clean_policy":       "ON_BRANCH_CHANGE",
		"git_path":           "/usr/bin/git",
		"clean_files_policy": "IGNORED_ONLY",
		"use_mirrors":        "true",
	}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVcsRootGitDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVcsRootGitAgentSettings,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVcsRootGitExists(resourceName, &vcs),
					testAccCheckVcsRootGitAgentSettings(&vcs, expected),
					resource.TestCheckResourceAttr(resourceName, "agent.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "agent.2618137949.git_path", "/usr/bin/git"),
					resource.TestCheckResourceAttr(resourceName, "agent.2618137949.clean_policy", "branch_change"),
					resource.TestCheckResourceAttr(resourceName, "agent.2618137949.clean_files_policy", "ignored_only"),
					resource.TestCheckResourceAttr(resourceName, "agent.2618137949.use_mirrors", "true"),
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

func testAccCheckVcsRootGitAgentSettings(vcs *api.GitVcsRoot, expected map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		dt, err := client.VcsRoots.GetByID((*vcs).ID)
		if err != nil {
			return err
		}
		actual := dt.(*api.GitVcsRoot)
		as := actual.Options.AgentSettings

		if string(as.CleanPolicy) != expected["clean_policy"] {
			return fmt.Errorf("agent setting %s: got '%s', expected '%s'", "clean_policy", as.CleanPolicy, expected["clean_policy"])
		}
		if string(as.CleanFilesPolicy) != expected["clean_files_policy"] {
			return fmt.Errorf("agent setting %s: got '%s', expected '%s'", "clean_files_policy", as.CleanFilesPolicy, expected["clean_files_policy"])
		}
		if as.GitPath != expected["git_path"] {
			return fmt.Errorf("agent setting %s: got '%s', expected '%s'", "git_path", as.GitPath, expected["git_path"])
		}
		if strconv.FormatBool(as.UseMirrors) != expected["use_mirrors"] {
			return fmt.Errorf("agent setting %s: got '%v', expected '%s'", "use_mirrors", as.UseMirrors, expected["use_mirrors"])
		}
		return nil
	}
}

func testAccVcsRootGitConfig(projectId string) string {
	return fmt.Sprintf(`
resource "teamcity_vcs_root_git" "git_test" {
	name = "application"
	project_id = "%s"
	fetch_url = "https://github.com/cvbarros/terraform-provider-teamcity"
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
	fetch_url = "https://github.com/cvbarros/terraform-provider-teamcity"
	default_branch = "refs/head/master"
	branches = [
    "+:refs/(pull/*)/head",
    "+:refs/heads/develop",
  	]
	username_style = "userid"
	submodule_checkout = "checkout"
	enable_branch_spec_tags = true
	modification_check_interval = 60
}
`

const testAccVcsRootGitUpdated = `
resource "teamcity_project" "vcs_root_project" {
  name = "vcs_root_project"
}

resource "teamcity_project" "vcs_root_project_new" {
	name = "vcs_root_project_new"
  }

resource "teamcity_vcs_root_git" "git_test" {
	name = "application_updated"
	project_id = "${teamcity_project.vcs_root_project_new.id}"
	fetch_url = "https://github.com/cvbarros/go-teamcity-sdk"
	default_branch = "refs/head/develop"
	branches = [
    	"+:refs/heads/master",
  	]
	username_style = "author_name"
	submodule_checkout = "ignore"
	enable_branch_spec_tags = false
	modification_check_interval = 180
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

const testAccVcsRootGitAgentSettings = `
resource "teamcity_project" "vcs_root_project" {
  name = "vcs_root_project"
}

resource "teamcity_vcs_root_git" "git_test" {
	name = "application"
	project_id = "${teamcity_project.vcs_root_project.id}"
	fetch_url = "https://github.com/cvbarros/terraform-provider-teamcity"
	default_branch = "refs/head/master"

	agent {
		git_path = "/usr/bin/git"
		clean_policy = "branch_change"
		clean_files_policy = "ignored_only"
		use_mirrors = true
	}
}
`

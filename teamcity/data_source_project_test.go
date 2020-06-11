package teamcity_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceProject_Root(t *testing.T) {
	resName := "data.teamcity_project.root"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProjectRoot,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", "<Root project>"),
					resource.TestCheckResourceAttr(resName, "project_id", "_Root"),
					resource.TestCheckResourceAttr(resName, "parent_project_id", ""),
					resource.TestCheckResourceAttr(resName, "url", "http://127.0.0.1:8112/project.html?projectId=_Root"),
				),
			},
		},
	})
}

func TestAccDataSourceProject_Basic(t *testing.T) {
	resName := "data.teamcity_project.project"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProject,
			},
			{
				Config: testAccDataSourceProject,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", "Test Project"),
					resource.TestCheckResourceAttr(resName, "project_id", "TestProject"),
					resource.TestCheckResourceAttr(resName, "parent_project_id", ""),
					resource.TestCheckResourceAttr(resName, "url", "http://127.0.0.1:8112/project.html?projectId=TestProject"),
				),
			},
		},
	})
}

const testAccDataSourceProjectRoot = `
data "teamcity_project" "root" {
  name = "<Root project>"
}
`

const testAccDataSourceProject = `
resource "teamcity_project" "project" {
	name = "Test Project"
}

data "teamcity_project" "project" {
  name = teamcity_project.project.name
}
`

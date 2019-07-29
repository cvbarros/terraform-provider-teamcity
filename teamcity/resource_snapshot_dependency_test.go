package teamcity_test

import (
	"fmt"
	"strings"
	"testing"

	api "github.com/cvbarros/go-teamcity-sdk/teamcity"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTeamcitySnapshotDependency_Basic(t *testing.T) {
	resName := "teamcity_snapshot_dependency.test"
	sd := api.SnapshotDependency{SourceBuildType: &api.BuildTypeReference{}}
	var bc api.BuildType

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamcitySnapshotDependencyDestroy(&sd.BuildTypeID),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccSnapshotDependencyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists("teamcity_build_config.config", &bc),
					testAccCheckTeamcitySnapshotDependencyExists(resName, &bc.ID, &sd),
					resource.TestCheckResourceAttrPtr(resName, "build_config_id", &sd.BuildTypeID),
					testAccCheckSnapshotSourceBuildType(resName, &sd),
				),
			},
		},
	})
}

func TestAccTeamcitySnapshotDependency_Updated(t *testing.T) {
	resName := "teamcity_snapshot_dependency.test"
	sd := api.SnapshotDependency{SourceBuildType: &api.BuildTypeReference{}}
	var bc api.BuildType

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamcitySnapshotDependencyDestroy(&sd.BuildTypeID),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccSnapshotDependencyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists("teamcity_build_config.config", &bc),
					testAccCheckTeamcitySnapshotDependencyExists(resName, &bc.ID, &sd),
					resource.TestCheckResourceAttrPtr(resName, "build_config_id", &sd.BuildTypeID),
					testAccCheckSnapshotSourceBuildType(resName, &sd),
				),
			},
			resource.TestStep{
				Config: TestAccSnapshotDependencyBasicUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists("teamcity_build_config.config", &bc),
					testAccCheckTeamcitySnapshotDependencyExists(resName, &bc.ID, &sd),
					resource.TestCheckResourceAttr(resName, "source_build_config_id", "Snapshot_Dependency2"),
				),
			},
		},
	})
}

func testAccCheckSnapshotSourceBuildType(n string, sd *api.SnapshotDependency) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		key := "source_build_config_id"
		value := (*sd).SourceBuildType.ID

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if v, ok := rs.Primary.Attributes[key]; !ok || v != value {
			if !ok {
				return fmt.Errorf("%s: Attribute '%s' not found", n, key)
			}

			return fmt.Errorf(
				"%s: Attribute '%s' expected %#v, got %#v",
				n,
				key,
				value,
				v)
		}
		return nil
	}
}

func testAccCheckTeamcitySnapshotDependencyDestroy(bt *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		return snapshotDependencyDestroyHelper(s, bt, client)
	}
}

func snapshotDependencyDestroyHelper(s *terraform.State, bt *string, client *api.Client) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != "teamcity_snapshot_dependency" {
			continue
		}

		dep := client.DependencyService(*bt)
		_, err := dep.GetSnapshotByID(r.Primary.ID)

		if err != nil {
			if strings.Contains(err.Error(), "404") {
				continue
			}
			return fmt.Errorf("Received an error retrieving the Snapshot dependency: %s", err)
		}

		return fmt.Errorf("Snapshot Dependency still exists")
	}
	return nil
}

func testAccCheckTeamcitySnapshotDependencyExists(n string, bt *string, snap *api.SnapshotDependency) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		return teamcitySnapshotDependencyExistsHelper(n, bt, s, client, snap)
	}
}

func teamcitySnapshotDependencyExistsHelper(n string, bt *string, s *terraform.State, client *api.Client, snap *api.SnapshotDependency) error {
	rs, ok := s.RootModule().Resources[n]
	if !ok {
		return fmt.Errorf("Not found: %s", n)
	}

	if rs.Primary.ID == "" {
		return fmt.Errorf("No ID is set")
	}

	dep := client.DependencyService(*bt)
	out, err := dep.GetSnapshotByID(rs.Primary.ID)
	if err != nil {
		return fmt.Errorf("Received an error retrieving snapshot dependency: %s", err)
	}

	*snap = *out
	return nil
}

const TestAccSnapshotDependencyBasic = `
resource "teamcity_project" "snapshop_dependency_project_test" {
  name = "Snapshot"
}

resource "teamcity_build_config" "dependency" {
	name = "Dependency"
	project_id = "${teamcity_project.snapshop_dependency_project_test.id}"
}

resource "teamcity_build_config" "config" {
	name = "BuildConfig"
	project_id = "${teamcity_project.snapshop_dependency_project_test.id}"
}

resource "teamcity_snapshot_dependency" "test" {
	source_build_config_id = "${teamcity_build_config.dependency.id}"
	build_config_id = "${teamcity_build_config.config.id}"
}
`

const TestAccSnapshotDependencyBasicUpdated = `
resource "teamcity_project" "snapshop_dependency_project_test" {
  name = "Snapshot"
}

resource "teamcity_build_config" "dependency" {
	name = "Dependency"
	project_id = "${teamcity_project.snapshop_dependency_project_test.id}"
}

resource "teamcity_build_config" "dependency2" {
	name = "Dependency 2"
	project_id = "${teamcity_project.snapshop_dependency_project_test.id}"
}

resource "teamcity_build_config" "config" {
	name = "BuildConfig"
	project_id = "${teamcity_project.snapshop_dependency_project_test.id}"
}

resource "teamcity_snapshot_dependency" "test" {
	source_build_config_id = "${teamcity_build_config.dependency2.id}"
	build_config_id = "${teamcity_build_config.config.id}"
}
`

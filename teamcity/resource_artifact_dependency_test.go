package teamcity_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	api "github.com/cvbarros/go-teamcity-sdk/teamcity"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTeamcityArtifactDependency_Basic(t *testing.T) {
	resName := "teamcity_artifact_dependency.test"
	var dep api.ArtifactDependency
	var bc api.BuildType
	var sb api.BuildType

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamcityArtifactDependencyDestroy(&bc.ID, "teamcity_artifact_dependency"),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccArtifactDependencyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists("teamcity_build_config.config", &bc),
					testAccCheckBuildConfigExists("teamcity_build_config.dependency", &sb),
					testAccCheckTeamcityArtifactDependencyExists(resName, &bc.ID, &dep),
					resource.TestCheckResourceAttrPtr(resName, "build_config_id", &bc.ID),
					resource.TestCheckResourceAttrPtr(resName, "source_build_config_id", &sb.ID),
					resource.TestCheckResourceAttr(resName, "path_rules.0", "+:*"),
				),
			},
		},
	})
}

func TestAccTeamcityArtifactDependency_Update(t *testing.T) {
	resName := "teamcity_artifact_dependency.test"
	var dep api.ArtifactDependency
	var bc api.BuildType
	var sb api.BuildType

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamcityArtifactDependencyDestroy(&bc.ID, "teamcity_artifact_dependency"),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccArtifactDependencyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists("teamcity_build_config.config", &bc),
					testAccCheckBuildConfigExists("teamcity_build_config.dependency", &sb),
					testAccCheckTeamcityArtifactDependencyExists(resName, &bc.ID, &dep),
					resource.TestCheckResourceAttrPtr(resName, "build_config_id", &bc.ID),
					resource.TestCheckResourceAttrPtr(resName, "source_build_config_id", &sb.ID),
					resource.TestCheckResourceAttr(resName, "path_rules.0", "+:*"),
				),
			},
			resource.TestStep{
				Config: TestAccArtifactDependencyBasicUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists("teamcity_build_config.config", &bc),
					testAccCheckBuildConfigExists("teamcity_build_config.dependency", &sb),
					testAccCheckTeamcityArtifactDependencyExists(resName, &bc.ID, &dep),
					resource.TestCheckResourceAttrPtr(resName, "build_config_id", &bc.ID),
					resource.TestCheckResourceAttr(resName, "source_build_config_id", "ArtifactDependency_Dependency2"),
					resource.TestCheckResourceAttr(resName, "path_rules.0", "-:*.md"),
				),
			},
		},
	})
}

func TestAccTeamcityArtifactDependency_DependencyRevisionUpdate(t *testing.T) {
	resName := "teamcity_artifact_dependency.test"
	var dep api.ArtifactDependency
	var bc api.BuildType
	var sb api.BuildType

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamcityArtifactDependencyDestroy(&bc.ID, "teamcity_artifact_dependency"),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccArtifactDependencyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists("teamcity_build_config.config", &bc),
					testAccCheckBuildConfigExists("teamcity_build_config.dependency", &sb),
					testAccCheckTeamcityArtifactDependencyExists(resName, &bc.ID, &dep),
					resource.TestCheckResourceAttrPtr(resName, "build_config_id", &bc.ID),
					resource.TestCheckResourceAttrPtr(resName, "source_build_config_id", &sb.ID),
					resource.TestCheckResourceAttr(resName, "path_rules.0", "+:*"),
				),
			},
			resource.TestStep{
				Config:             TestAccArtifactDependencyDependencyRevisionUpdated,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTeamcityArtifactDependency_OptionsNoRevision(t *testing.T) {
	resName := "teamcity_artifact_dependency.test"
	var dep api.ArtifactDependency
	var bc api.BuildType
	var sb api.BuildType

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamcityArtifactDependencyDestroy(&bc.ID, "teamcity_artifact_dependency"),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccArtifactDependencyOptionsNoRevision,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists("teamcity_build_config.config", &bc),
					testAccCheckBuildConfigExists("teamcity_build_config.dependency", &sb),
					testAccCheckTeamcityArtifactDependencyExists(resName, &bc.ID, &dep),
					resource.TestCheckResourceAttrPtr(resName, "build_config_id", &bc.ID),
					resource.TestCheckResourceAttrPtr(resName, "source_build_config_id", &sb.ID),
					resource.TestCheckResourceAttr(resName, "path_rules.0", "+:*"),
					resource.TestCheckResourceAttr(resName, "dependency_revision", "lastFinished"),
					resource.TestCheckResourceAttr(resName, "clean_destination", "true"),
				),
			},
		},
	})
}

func TestAccTeamcityArtifactDependency_OptionsRevision(t *testing.T) {
	resName := "teamcity_artifact_dependency.test"
	var dep api.ArtifactDependency
	var bc api.BuildType
	var sb api.BuildType

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamcityArtifactDependencyDestroy(&bc.ID, "teamcity_artifact_dependency"),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccArtifactDependencyOptionsRevision,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildConfigExists("teamcity_build_config.config", &bc),
					testAccCheckBuildConfigExists("teamcity_build_config.dependency", &sb),
					testAccCheckTeamcityArtifactDependencyExists(resName, &bc.ID, &dep),
					resource.TestCheckResourceAttrPtr(resName, "build_config_id", &bc.ID),
					resource.TestCheckResourceAttrPtr(resName, "source_build_config_id", &sb.ID),
					resource.TestCheckResourceAttr(resName, "path_rules.0", "+:*"),
					resource.TestCheckResourceAttr(resName, "dependency_revision", "buildTag"),
					resource.TestCheckResourceAttr(resName, "revision", "stable"),
					resource.TestCheckResourceAttr(resName, "clean_destination", "false"),
				),
			},
		},
	})
}

func TestAccTeamcityArtifactDependency_ConfigErrorForRevision(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      TestAccArtifactDependencyConfigErrorRevision,
				ExpectError: regexp.MustCompile("'revision' property is required if using 'buildTag' or 'buildNumber' for 'dependency_revision'"),
			},
		},
	})
}

func testAccCheckDependencySourceBuildType(n string, dep *api.ArtifactDependency) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		key := "source_build_config_id"
		value := dep.SourceBuildTypeID

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

func testAccCheckTeamcityArtifactDependencyDestroy(bt *string, resourceType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		return ArtifactDependencyDestroyHelper(s, bt, client, resourceType)
	}
}

func ArtifactDependencyDestroyHelper(s *terraform.State, bt *string, client *api.Client, resourceType string) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != resourceType {
			continue
		}

		dep := client.DependencyService(*bt)
		_, err := dep.GetSnapshotByID(r.Primary.ID)

		if err != nil {
			if strings.Contains(err.Error(), "404") {
				continue
			}
			return fmt.Errorf("Received an error retrieving the dependency: %s", err)
		}

		return fmt.Errorf("Dependency still exists")
	}
	return nil
}

func testAccCheckTeamcityArtifactDependencyExists(n string, bt *string, snap *api.ArtifactDependency) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		return teamcityArtifactDependencyExistsHelper(n, bt, s, client, snap)
	}
}

func teamcityArtifactDependencyExistsHelper(n string, bt *string, s *terraform.State, client *api.Client, dep *api.ArtifactDependency) error {
	rs, ok := s.RootModule().Resources[n]
	if !ok {
		return fmt.Errorf("Not found: %s", n)
	}

	if rs.Primary.ID == "" {
		return fmt.Errorf("No ID is set")
	}

	depSrv := client.DependencyService(*bt)
	out, err := depSrv.GetArtifactByID(rs.Primary.ID)
	if err != nil {
		return fmt.Errorf("Received an error retrieving dependency: %s", err)
	}

	*dep = *out
	return nil
}

const TestAccArtifactDependencyBasic = `
resource "teamcity_project" "artifact_dependency_project_test" {
  name = "Artifact Dependency"
}

resource "teamcity_build_config" "dependency" {
	name = "Dependency"
	project_id = "${teamcity_project.artifact_dependency_project_test.id}"
}

resource "teamcity_build_config" "config" {
	name = "BuildConfig"
	project_id = "${teamcity_project.artifact_dependency_project_test.id}"
}

resource "teamcity_artifact_dependency" "test" {
	source_build_config_id = "${teamcity_build_config.dependency.id}"
	build_config_id = "${teamcity_build_config.config.id}"

	path_rules = ["+:*"]
}
`

const TestAccArtifactDependencyBasicUpdated = `
resource "teamcity_project" "artifact_dependency_project_test" {
  name = "Artifact Dependency"
}

resource "teamcity_build_config" "dependency" {
	name = "Dependency"
	project_id = "${teamcity_project.artifact_dependency_project_test.id}"
}

resource "teamcity_build_config" "dependency2" {
	name = "Dependency 2"
	project_id = "${teamcity_project.artifact_dependency_project_test.id}"
}

resource "teamcity_build_config" "config" {
	name = "BuildConfig"
	project_id = "${teamcity_project.artifact_dependency_project_test.id}"
}

resource "teamcity_artifact_dependency" "test" {
	source_build_config_id = "${teamcity_build_config.dependency2.id}"
	build_config_id = "${teamcity_build_config.config.id}"

	path_rules = ["-:*.md"]
}
`

const TestAccArtifactDependencyOptionsNoRevision = `
resource "teamcity_project" "artifact_dependency_project_test" {
  name = "Artifact Dependency"
}

resource "teamcity_build_config" "dependency" {
	name = "Dependency"
	project_id = "${teamcity_project.artifact_dependency_project_test.id}"
}

resource "teamcity_build_config" "config" {
	name = "BuildConfig"
	project_id = "${teamcity_project.artifact_dependency_project_test.id}"
}

resource "teamcity_artifact_dependency" "test" {
	source_build_config_id = "${teamcity_build_config.dependency.id}"
	build_config_id = "${teamcity_build_config.config.id}"

	path_rules = ["+:*"]

	clean_destination = true
	dependency_revision = "lastFinished"
}
`

const TestAccArtifactDependencyOptionsRevision = `
resource "teamcity_project" "artifact_dependency_project_test" {
  name = "Artifact Dependency"
}

resource "teamcity_build_config" "dependency" {
	name = "Dependency"
	project_id = "${teamcity_project.artifact_dependency_project_test.id}"
}

resource "teamcity_build_config" "config" {
	name = "BuildConfig"
	project_id = "${teamcity_project.artifact_dependency_project_test.id}"
}

resource "teamcity_artifact_dependency" "test" {
	source_build_config_id = "${teamcity_build_config.dependency.id}"
	build_config_id = "${teamcity_build_config.config.id}"

	path_rules = ["+:*"]

	dependency_revision = "buildTag"
	revision = "stable"
}
`

const TestAccArtifactDependencyConfigErrorRevision = `
resource "teamcity_project" "artifact_dependency_project_test" {
  name = "Artifact Dependency"
}

resource "teamcity_build_config" "dependency" {
	name = "Dependency"
	project_id = "${teamcity_project.artifact_dependency_project_test.id}"
}

resource "teamcity_build_config" "config" {
	name = "BuildConfig"
	project_id = "${teamcity_project.artifact_dependency_project_test.id}"
}

resource "teamcity_artifact_dependency" "test" {
	source_build_config_id = "${teamcity_build_config.dependency.id}"
	build_config_id = "${teamcity_build_config.config.id}"

	path_rules = ["+:*"]

	dependency_revision = "buildNumber"
	#Missing revision required property
}
`

const TestAccArtifactDependencyDependencyRevisionUpdated = `
resource "teamcity_project" "artifact_dependency_project_test" {
  name = "Artifact Dependency"
}

resource "teamcity_build_config" "dependency" {
	name = "Dependency"
	project_id = "${teamcity_project.artifact_dependency_project_test.id}"
}

resource "teamcity_build_config" "config" {
	name = "BuildConfig"
	project_id = "${teamcity_project.artifact_dependency_project_test.id}"
}

resource "teamcity_artifact_dependency" "test" {
	source_build_config_id = "${teamcity_build_config.dependency.id}"
	build_config_id = "${teamcity_build_config.config.id}"

	path_rules = ["+:*"]
	dependency_revision = "lastFinished" #Added
}
`

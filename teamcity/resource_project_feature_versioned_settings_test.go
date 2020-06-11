package teamcity_test

import (
	"fmt"
	"strings"
	"testing"

	api "github.com/cvbarros/go-teamcity/teamcity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTeamcityProjectVersionedSettings_Basic(t *testing.T) {
	resName := "teamcity_project_feature_versioned_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamCityProjectVersionedSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamCityProjectVersionedSettingsBasicConfig("PREFER_VCS", "kotlin"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamCityProjectVersionedSettingsExists(resName),
					resource.TestCheckResourceAttr(resName, "build_settings", "PREFER_VCS"),
					resource.TestCheckResourceAttr(resName, "format", "kotlin"),
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

func TestAccTeamcityProjectVersionedSettings_Update(t *testing.T) {
	resName := "teamcity_project_feature_versioned_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamCityProjectVersionedSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamCityProjectVersionedSettingsBasicConfig("PREFER_VCS", "kotlin"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamCityProjectVersionedSettingsExists(resName),
					resource.TestCheckResourceAttr(resName, "build_settings", "PREFER_VCS"),
					resource.TestCheckResourceAttr(resName, "format", "kotlin"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeamCityProjectVersionedSettingsBasicConfig("ALWAYS_USE_CURRENT", "xml"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamCityProjectVersionedSettingsExists(resName),
					resource.TestCheckResourceAttr(resName, "build_settings", "ALWAYS_USE_CURRENT"),
					resource.TestCheckResourceAttr(resName, "format", "xml"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeamCityProjectVersionedSettingsBasicConfig("PREFER_VCS", "xml"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamCityProjectVersionedSettingsExists(resName),
					resource.TestCheckResourceAttr(resName, "build_settings", "PREFER_VCS"),
					resource.TestCheckResourceAttr(resName, "format", "xml"),
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

func TestAccTeamcityProjectVersionedSettings_ContextParameters(t *testing.T) {
	resName := "teamcity_project_feature_versioned_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamCityProjectVersionedSettingsDestroy,
		Steps: []resource.TestStep{
			{
				// single
				Config: testAccTeamCityProjectVersionedSettingsContextParametersConfig("PREFER_VCS", "kotlin"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamCityProjectVersionedSettingsExists(resName),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// multiple
				Config: testAccTeamCityProjectVersionedSettingsContextParametersUpdatedConfig("PREFER_VCS", "kotlin"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamCityProjectVersionedSettingsExists(resName),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// removed
				Config: testAccTeamCityProjectVersionedSettingsBasicConfig("PREFER_VCS", "kotlin"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamCityProjectVersionedSettingsExists(resName),
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

func TestAccTeamcityProjectVersionedSettings_CredentialsStorageTypeSettings(t *testing.T) {
	resName := "teamcity_project_feature_versioned_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamCityProjectVersionedSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamCityProjectVersionedSettingsCredentialsStorageTypeConfig("scrambled"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamCityProjectVersionedSettingsExists(resName),
					resource.TestCheckResourceAttr(resName, "credentials_storage_type", "scrambled"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeamCityProjectVersionedSettingsCredentialsStorageTypeConfig("credentialsJSON"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamCityProjectVersionedSettingsExists(resName),
					resource.TestCheckResourceAttr(resName, "credentials_storage_type", "credentialsJSON"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeamCityProjectVersionedSettingsCredentialsStorageTypeConfig("scrambled"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamCityProjectVersionedSettingsExists(resName),
					resource.TestCheckResourceAttr(resName, "credentials_storage_type", "scrambled"),
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

func testAccCheckTeamCityProjectVersionedSettingsExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		service := client.ProjectFeatureService(rs.Primary.ID)
		feature, err := service.GetByType("versionedSettings")
		if err != nil {
			return fmt.Errorf("Received an error retrieving project versioned settings: %s", err)
		}

		if _, ok := feature.(*api.ProjectFeatureVersionedSettings); !ok {
			return fmt.Errorf("Expected a Versioned Setting but it wasn't!")
		}

		return nil
	}
}

func testAccCheckTeamCityProjectVersionedSettingsDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "teamcity_project_feature_versioned_settings" {
			continue
		}
		service := client.ProjectFeatureService(r.Primary.ID)
		if _, err := service.GetByType("versionedSettings"); err != nil {
			if strings.Contains(err.Error(), "404") {
				// expected, since it's gone
				continue
			}

			return fmt.Errorf("Received an error retrieving project versioned settings: %s", err)
		}

		return fmt.Errorf("Project Versioned Settings still exists")
	}
	return nil
}

func testAccTeamCityProjectVersionedSettingsBasicConfig(buildSettings string, format string) string {
	return fmt.Sprintf(`
%s

resource "teamcity_project_feature_versioned_settings" "test" {
  project_id     = teamcity_project.test.id
  vcs_root_id    = teamcity_vcs_root_git.test.id
  build_settings = "%s"
  format         = "%s"
}
`, testAccTeamCityProjectVersionedSettingsTemplate, buildSettings, format)
}

func testAccTeamCityProjectVersionedSettingsContextParametersConfig(buildSettings string, format string) string {
	return fmt.Sprintf(`
%s

resource "teamcity_project_feature_versioned_settings" "test" {
  project_id     = teamcity_project.test.id
  vcs_root_id    = teamcity_vcs_root_git.test.id
  build_settings = "%s"
  format         = "%s"

  context_parameters = {
    Hello = "World"
  }
}
`, testAccTeamCityProjectVersionedSettingsTemplate, buildSettings, format)
}

func testAccTeamCityProjectVersionedSettingsContextParametersUpdatedConfig(buildSettings string, format string) string {
	return fmt.Sprintf(`
%s

resource "teamcity_project_feature_versioned_settings" "test" {
  project_id     = teamcity_project.test.id
  vcs_root_id    = teamcity_vcs_root_git.test.id
  build_settings = "%s"
  format         = "%s"

  context_parameters = {
    Hello = "World"
    abc   = 123
  }
}
`, testAccTeamCityProjectVersionedSettingsTemplate, buildSettings, format)
}

func testAccTeamCityProjectVersionedSettingsCredentialsStorageTypeConfig(credentialsStorageType string) string {
	return fmt.Sprintf(`
%s

resource "teamcity_project_feature_versioned_settings" "test" {
  project_id               = teamcity_project.test.id
  vcs_root_id              = teamcity_vcs_root_git.test.id
  build_settings           = "PREFER_VCS"
  format                   = "kotlin"
  credentials_storage_type = "%s"
}
`, testAccTeamCityProjectVersionedSettingsTemplate, credentialsStorageType)
}

const testAccTeamCityProjectVersionedSettingsTemplate = `
resource "teamcity_project" "test" {
  name = "Test Project"
}

resource "teamcity_vcs_root_git" "test" {
  name          = "application"
  project_id     = teamcity_project.test.id
  fetch_url      = "https://github.com/cvbarros/terraform-provider-teamcity"
  default_branch = "refs/head/master"
}
`

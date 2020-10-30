package teamcity_test

import (
	"fmt"
	"strings"
	"testing"

	api "github.com/cvbarros/go-teamcity/teamcity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTeamcityProjectVaultConnection_Basic(t *testing.T) {
	resName := "teamcity_project_feature_vault_connection.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamCityProjectVaultConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamCityProjectVaultConnectionBasicConfig("http://vault.service:8200"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamCityProjectVaultConnectionExists(resName),
					resource.TestCheckResourceAttr(resName, "auth_method", "approle"),
					resource.TestCheckResourceAttr(resName, "url", "http://vault.service:8200"),
				),
			},
			// {
			// 	ResourceName:      resName,
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
		},
	})
}

// func TestAccTeamcityProjectVersionedSettings_Update(t *testing.T) {
// 	resName := "teamcity_project_feature_versioned_settings.test"

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckTeamCityProjectVersionedSettingsDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccTeamCityProjectVersionedSettingsBasicConfig("PREFER_VCS", "kotlin"),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckTeamCityProjectVersionedSettingsExists(resName),
// 					resource.TestCheckResourceAttr(resName, "build_settings", "PREFER_VCS"),
// 					resource.TestCheckResourceAttr(resName, "format", "kotlin"),
// 				),
// 			},
// 			{
// 				ResourceName:      resName,
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 			{
// 				Config: testAccTeamCityProjectVersionedSettingsBasicConfig("ALWAYS_USE_CURRENT", "xml"),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckTeamCityProjectVersionedSettingsExists(resName),
// 					resource.TestCheckResourceAttr(resName, "build_settings", "ALWAYS_USE_CURRENT"),
// 					resource.TestCheckResourceAttr(resName, "format", "xml"),
// 				),
// 			},
// 			{
// 				ResourceName:      resName,
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 			{
// 				Config: testAccTeamCityProjectVersionedSettingsBasicConfig("PREFER_VCS", "xml"),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckTeamCityProjectVersionedSettingsExists(resName),
// 					resource.TestCheckResourceAttr(resName, "build_settings", "PREFER_VCS"),
// 					resource.TestCheckResourceAttr(resName, "format", "xml"),
// 				),
// 			},
// 			{
// 				ResourceName:      resName,
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 		},
// 	})
// }

// func TestAccTeamcityProjectVersionedSettings_ContextParameters(t *testing.T) {
// 	resName := "teamcity_project_feature_versioned_settings.test"

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckTeamCityProjectVersionedSettingsDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				// single
// 				Config: testAccTeamCityProjectVersionedSettingsContextParametersConfig("PREFER_VCS", "kotlin"),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckTeamCityProjectVersionedSettingsExists(resName),
// 				),
// 			},
// 			{
// 				ResourceName:      resName,
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 			{
// 				// multiple
// 				Config: testAccTeamCityProjectVersionedSettingsContextParametersUpdatedConfig("PREFER_VCS", "kotlin"),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckTeamCityProjectVersionedSettingsExists(resName),
// 				),
// 			},
// 			{
// 				ResourceName:      resName,
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 			{
// 				// removed
// 				Config: testAccTeamCityProjectVersionedSettingsBasicConfig("PREFER_VCS", "kotlin"),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckTeamCityProjectVersionedSettingsExists(resName),
// 				),
// 			},
// 			{
// 				ResourceName:      resName,
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 		},
// 	})
// }

// func TestAccTeamcityProjectVersionedSettings_CredentialsStorageTypeSettings(t *testing.T) {
// 	resName := "teamcity_project_feature_versioned_settings.test"

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckTeamCityProjectVersionedSettingsDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccTeamCityProjectVersionedSettingsCredentialsStorageTypeConfig("scrambled"),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckTeamCityProjectVersionedSettingsExists(resName),
// 					resource.TestCheckResourceAttr(resName, "credentials_storage_type", "scrambled"),
// 				),
// 			},
// 			{
// 				ResourceName:      resName,
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 			{
// 				Config: testAccTeamCityProjectVersionedSettingsCredentialsStorageTypeConfig("credentialsJSON"),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckTeamCityProjectVersionedSettingsExists(resName),
// 					resource.TestCheckResourceAttr(resName, "credentials_storage_type", "credentialsJSON"),
// 				),
// 			},
// 			{
// 				ResourceName:      resName,
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 			{
// 				Config: testAccTeamCityProjectVersionedSettingsCredentialsStorageTypeConfig("scrambled"),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckTeamCityProjectVersionedSettingsExists(resName),
// 					resource.TestCheckResourceAttr(resName, "credentials_storage_type", "scrambled"),
// 				),
// 			},
// 			{
// 				ResourceName:      resName,
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 		},
// 	})
// }

func testAccCheckTeamCityProjectVaultConnectionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		service := client.ProjectFeatureService(rs.Primary.ID)
		feature, err := service.GetByTypeAndProvider("OAuthProvider", "teamcity-vault")
		if err != nil {
			return fmt.Errorf("Received an error retrieving project Vault connection: %s", err)
		}

		if _, ok := feature.(*api.ConnectionProviderVault); !ok {
			return fmt.Errorf("Expected a Vault connection provider but it wasn't!")
		}

		return nil
	}
}

func testAccCheckTeamCityProjectVaultConnectionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "teamcity_project_feature_vault_connection" {
			continue
		}
		service := client.ProjectFeatureService(r.Primary.ID)
		if _, err := service.GetByTypeAndProvider("OAuthProvider", "teamcity-vault"); err != nil {
			if strings.Contains(err.Error(), "404") {
				// expected, since it's gone
				continue
			}

			return fmt.Errorf("Received an error retrieving project Vault connection: %s", err)
		}

		return fmt.Errorf("Project Vault connection still exists")
	}
	return nil
}

func testAccTeamCityProjectVaultConnectionBasicConfig(vault_url string) string {
	return fmt.Sprintf(`
%s

resource "teamcity_project_feature_vault_connection" "test" {
	project_id = teamcity_project.test.id
	role_id    = "123456"
	secret_id  = "abcdef"
  url        = "%s"
}
`, testAccTeamCityProjectVaultConnectionTemplate, vault_url)
}

// func testAccTeamCityProjectVersionedSettingsContextParametersConfig(buildSettings string, format string) string {
// 	return fmt.Sprintf(`
// %s

// resource "teamcity_project_feature_versioned_settings" "test" {
//   project_id     = teamcity_project.test.id
//   vcs_root_id    = teamcity_vcs_root_git.test.id
//   build_settings = "%s"
//   format         = "%s"

//   context_parameters = {
//     Hello = "World"
//   }
// }
// `, testAccTeamCityProjectVersionedSettingsTemplate, buildSettings, format)
// }

// func testAccTeamCityProjectVersionedSettingsContextParametersUpdatedConfig(buildSettings string, format string) string {
// 	return fmt.Sprintf(`
// %s

// resource "teamcity_project_feature_versioned_settings" "test" {
//   project_id     = teamcity_project.test.id
//   vcs_root_id    = teamcity_vcs_root_git.test.id
//   build_settings = "%s"
//   format         = "%s"

//   context_parameters = {
//     Hello = "World"
//     abc   = 123
//   }
// }
// `, testAccTeamCityProjectVersionedSettingsTemplate, buildSettings, format)
// }

// func testAccTeamCityProjectVersionedSettingsCredentialsStorageTypeConfig(credentialsStorageType string) string {
// 	return fmt.Sprintf(`
// %s

// resource "teamcity_project_feature_versioned_settings" "test" {
//   project_id               = teamcity_project.test.id
//   vcs_root_id              = teamcity_vcs_root_git.test.id
//   build_settings           = "PREFER_VCS"
//   format                   = "kotlin"
//   credentials_storage_type = "%s"
// }
// `, testAccTeamCityProjectVersionedSettingsTemplate, credentialsStorageType)
// }

const testAccTeamCityProjectVaultConnectionTemplate = `
resource "teamcity_project" "test" {
  name = "Test Project"
}
`

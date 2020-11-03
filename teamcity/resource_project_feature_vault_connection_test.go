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
					resource.TestCheckResourceAttr(resName, "approle_role_id", "123456"),
					resource.TestCheckResourceAttr(resName, "approle_secret_id", "abcdef"),
					resource.TestCheckResourceAttr(resName, "url", "http://vault.service:8200"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"approle_secret_id",
				},
			},
		},
	})
}

func TestAccTeamcityProjectVaultConnection_Update(t *testing.T) {
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
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"approle_secret_id",
				},
			},
			{
				Config: testAccTeamCityProjectVaultConnectionBasicConfig("http://vault.anotherservice:8200"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamCityProjectVaultConnectionExists(resName),
					resource.TestCheckResourceAttr(resName, "auth_method", "approle"),
					resource.TestCheckResourceAttr(resName, "url", "http://vault.anotherservice:8200"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"approle_secret_id",
				},
			},
		},
	})
}

func TestAccTeamcityProjectVaultConnection_UpdateApproleAuth(t *testing.T) {
	resName := "teamcity_project_feature_vault_connection.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamCityProjectVaultConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamCityProjectVaultConnectionApproleConfig("abcdef", "123456"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamCityProjectVaultConnectionExists(resName),
					resource.TestCheckResourceAttr(resName, "approle_role_id", "abcdef"),
					resource.TestCheckResourceAttr(resName, "approle_secret_id", "123456"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"approle_secret_id",
				},
			},
			{
				Config: testAccTeamCityProjectVaultConnectionApproleConfig("zxywvu", "123456"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamCityProjectVaultConnectionExists(resName),
					resource.TestCheckResourceAttr(resName, "approle_role_id", "zxywvu"),
					resource.TestCheckResourceAttr(resName, "approle_secret_id", "123456"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"approle_secret_id",
				},
			},
			{
				Config: testAccTeamCityProjectVaultConnectionApproleConfig("zxywvu", "987654"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamCityProjectVaultConnectionExists(resName),
					resource.TestCheckResourceAttr(resName, "approle_role_id", "zxywvu"),
					resource.TestCheckResourceAttr(resName, "approle_secret_id", "987654"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"approle_secret_id",
				},
			},
		},
	})
}

func TestAccTeamcityProjectVaultConnection_IAMAuth(t *testing.T) {
	resName := "teamcity_project_feature_vault_connection.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamCityProjectVaultConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamCityProjectVaultConnectionIAMAuthConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamCityProjectVaultConnectionExists(resName),
					resource.TestCheckResourceAttr(resName, "auth_method", "iam"),
					resource.TestCheckNoResourceAttr(resName, "approle_role_id"),
					resource.TestCheckNoResourceAttr(resName, "approle_secret_id"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"approle_auth_path",
				},
			},
		},
	})
}

func TestAccTeamcityProjectVaultConnection_UpdateFailOnError(t *testing.T) {
	resName := "teamcity_project_feature_vault_connection.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamCityProjectVaultConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamCityProjectVaultConnectionFailOnErrorConfig(true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamCityProjectVaultConnectionExists(resName),
					resource.TestCheckResourceAttr(resName, "fail_on_error", "true"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"approle_secret_id",
				},
			},
			{
				Config: testAccTeamCityProjectVaultConnectionFailOnErrorConfig(false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamCityProjectVaultConnectionExists(resName),
					resource.TestCheckResourceAttr(resName, "fail_on_error", "false"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"approle_secret_id",
				},
			},
		},
	})
}

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
	project_id         = teamcity_project.test.id
	approle_role_id    = "123456"
	approle_secret_id  = "abcdef"
  url                = "%s"
}
`, testAccTeamCityProjectVaultConnectionTemplate, vault_url)
}

func testAccTeamCityProjectVaultConnectionApproleConfig(approle_role_id, approle_secret_id string) string {
	return fmt.Sprintf(`
%s

resource "teamcity_project_feature_vault_connection" "test" {
	project_id         = teamcity_project.test.id
	approle_role_id    = "%s"
	approle_secret_id  = "%s"
  url                = "https://vault.service:8200"
}
`, testAccTeamCityProjectVaultConnectionTemplate, approle_role_id, approle_secret_id)
}

func testAccTeamCityProjectVaultConnectionIAMAuthConfig() string {
	return fmt.Sprintf(`
%s

resource "teamcity_project_feature_vault_connection" "test" {
	project_id         = teamcity_project.test.id
	auth_method				 = "iam"
  url                = "https://vault.service:8200"
}
`, testAccTeamCityProjectVaultConnectionTemplate)
}

func testAccTeamCityProjectVaultConnectionFailOnErrorConfig(fail_on_error bool) string {
	return fmt.Sprintf(`
%s

resource "teamcity_project_feature_vault_connection" "test" {
	project_id         = teamcity_project.test.id
	approle_role_id    = "123456"
	approle_secret_id  = "abcdef"
	fail_on_error      = "%t"
	url                = "https://vault.service:8200"
}
`, testAccTeamCityProjectVaultConnectionTemplate, fail_on_error)
}

const testAccTeamCityProjectVaultConnectionTemplate = `
resource "teamcity_project" "test" {
  name = "Test Project"
}
`

package teamcity_test

import (
	"os"
	"testing"

	teamcity "github.com/cvbarros/terraform-provider-teamcity/teamcity"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = teamcity.Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"teamcity": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := teamcity.Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("Error: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = teamcity.Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("TEAMCITY_ADDR"); v == "" {
		t.Fatal("TEAMCITY_ADDR must be set for acceptance tests")
	}
	if v := os.Getenv("TEAMCITY_USER"); v == "" {
		t.Fatal("TEAMCITY_USER must be set for acceptance tests")
	}
	if v := os.Getenv("TEAMCITY_PASSWORD"); v == "" {
		t.Fatal("TEAMCITY_PASSWORD must be set for acceptance tests")
	}
}

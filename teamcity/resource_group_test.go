package teamcity_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"strings"
	"testing"

	api "github.com/cvbarros/go-teamcity/teamcity"
	"hash/crc32"
	"regexp"
)

func TestAccGroupCreate_Basic(t *testing.T) {
	var g api.Group
	resName := "teamcity_group.test_group"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGroupDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccGroupConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupExists(resName, &g),
					resource.TestCheckResourceAttr(resName, "key", generateKey("test-group")),
					resource.TestCheckResourceAttr(resName, "name", "test-group"),
					resource.TestCheckResourceAttr(resName, "description", "Description of test group"),
				),
			},
		},
	})
}

func TestAccGroupCreate_BasicUpdate(t *testing.T) {
	var g api.Group
	resName := "teamcity_group.test_group"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGroupDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccGroupConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupExists(resName, &g),
					resource.TestCheckResourceAttr(resName, "key", generateKey("test-group")),
					resource.TestCheckResourceAttr(resName, "name", "test-group"),
					resource.TestCheckResourceAttr(resName, "description", "Description of test group"),
				),
			},

			resource.TestStep{
				Config: TestAccGroupConfigBasicUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupExists(resName, &g),
					resource.TestCheckResourceAttr(resName, "key", generateKey("test-group-updated")),
					resource.TestCheckResourceAttr(resName, "name", "test-group-updated"),
					resource.TestCheckResourceAttr(resName, "description", "Updated description of test group"),
				),
			},
		},
	})
}

func generateKey(name string) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	processedName := reg.ReplaceAllString(strings.ToUpper(name), "")
	return fmt.Sprintf("%0.7s_%X", processedName, crc32.ChecksumIEEE([]byte(name)))
}

func testAccCheckGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	return buildGroupDestroyHelper(s, client)
}

func buildGroupDestroyHelper(s *terraform.State, client *api.Client) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != "teamcity_group" {
			continue
		}

		_, err := client.Groups.GetByKey(r.Primary.ID)

		if err != nil {
			if strings.Contains(err.Error(), "404") {
				continue
			}
			return fmt.Errorf("Received an error retrieving the Groups: %s", err)
		}

		return fmt.Errorf("Group still exists")
	}
	return nil
}

func testAccCheckGroupExists(n string, out *api.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		return groupExistsHelper(n, s, client, out)
	}
}

func groupExistsHelper(n string, s *terraform.State, client *api.Client, out *api.Group) error {
	rs, ok := s.RootModule().Resources[n]
	if !ok {
		return fmt.Errorf("Not found: %s", n)
	}

	if rs.Primary.ID == "" {
		return fmt.Errorf("No id for %s is set", n)
	}

	resp, err := client.Groups.GetByKey(rs.Primary.ID)

	if err != nil {
		return fmt.Errorf("Received an error retrieving Groups: %s", err)
	}

	*out = *resp

	return nil
}

const TestAccGroupConfigBasic = `
resource "teamcity_group" "test_group" {
  name = "test-group"
  description = "Description of test group"
}
`
const TestAccGroupConfigBasicUpdate = `
resource "teamcity_group" "test_group" {
  name = "test-group-updated"
  description = "Updated description of test group"
}
`

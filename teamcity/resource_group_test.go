package teamcity_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"hash/crc32"
	"regexp"

	api "github.com/cvbarros/go-teamcity/teamcity"
)

func TestAccGroup_Create(t *testing.T) {
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

func TestAccGroup_Update(t *testing.T) {
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

func TestAccGroup_DisallowImportOnCreate(t *testing.T) {
	resName := "teamcity_group.test_group"
	groupName := "test-group"
	testGroupKey := generateKey(groupName)
	groupDescription := "Description of test group"

	createGroup(testGroupKey, groupName, groupDescription)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGroupDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccGroupConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "key", testGroupKey),
					resource.TestCheckResourceAttr(resName, "name", groupName),
					resource.TestCheckResourceAttr(resName, "description", groupDescription),
				),
				ExpectError: regexp.MustCompile(".*group with the same key already exists.*"),
			},
		},
	})
}

func TestAccGroup_AllowImportOnCreate(t *testing.T) {
	resName := "teamcity_group.test_group"
	groupName := "test-group"
	testGroupKey := generateKey(groupName)
	groupDescription := "Description of test group"
	createGroup(testGroupKey, groupName, groupDescription)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGroupDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccGroupConfigBasicImport,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "key", testGroupKey),
					resource.TestCheckResourceAttr(resName, "name", groupName),
					resource.TestCheckResourceAttr(resName, "description", groupDescription),
				),
			},
		},
	})
}

func createGroup(testGroupKey string, groupName string, groupDescription string) {
	client, err := api.NewClient(api.BasicAuth("admin", "admin"), http.DefaultClient)
	if err == nil {
		newGroup, _ := api.NewGroup(testGroupKey, groupName, groupDescription)
		client.Groups.Create(newGroup)
		fmt.Sprintf("Made new group: %s", newGroup.Key)
	}
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
const TestAccGroupConfigBasicImport = `
resource "teamcity_group" "test_group" {
  name = "test-group"
  description = "Description of test group"
  import_if_exists = true
}
`

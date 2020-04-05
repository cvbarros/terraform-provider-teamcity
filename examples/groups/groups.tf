# This configuration sample shows how to manage projects
provider "teamcity" {
  address = var.teamcity_url
  username = var.teamcity_username
  password = var.teamcity_password
}

resource "teamcity_group" "test_group" {
  name = "test-group"
  description = "Description of test group"
}

resource "teamcity_group" "short_group" {
  name = "grp"
}

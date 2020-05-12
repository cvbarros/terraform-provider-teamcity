# This configuration sample shows how to manage projects
provider "teamcity" {
  address  = var.teamcity_url
  username = var.teamcity_username
  password = var.teamcity_password
}

# Assigning System Admin at the global level
resource "teamcity_group" "test_group_1" {
  name = "Test Group #1"
}

resource "teamcity_group_role_assignment" "test_group_1_sys_admin_global" {
  group_key  = teamcity_group.test_group_1.id
  role_id    = "SYSTEM_ADMIN"
  project_id = "g"
}

# Assigning Project Developer to the Test Project
resource "teamcity_group" "test_group_2" {
  name = "Test Group #2"
}

resource "teamcity_project" "test_project" {
  name = "Test Project"
}

resource "teamcity_group_role_assignment" "test_group_2_project_dev_test_project" {
  group_key  = teamcity_group.test_group_2.id
  role_id    = "PROJECT_DEVELOPER"
  project_id = "p:${teamcity_project.test_project.id}"
}

# This configuration sample shows how to manage projects
provider "teamcity" {
  address  = var.teamcity_url
  username = var.teamcity_username
  password = var.teamcity_password
}

resource "teamcity_project" "parent" {
  name        = "Parent"
  description = "Parent Project, will be created under the 'Root' project"

  env_params = {
    variable1 = "env_value1"
    variable2 = "env_value2"
  }
}

resource "teamcity_project" "child" {
  name        = "Child"
  description = "Child Project, will be created under 'Parent' project"
  parent_id   = teamcity_project.parent.id

  config_params = {
    variable1 = "config_value1"
    variable2 = "config_value2"
  }

  sys_params = {
    variable1 = "system_value1"
    variable2 = "system_value2"
  }
}

resource "teamcity_root_project" "root" {
  config_params = {
      variable1 = "config_value1"
      variable2 = "config_value2"
  }

  sys_params = {
      variable4 = "system_value1"
  }

  env_params = {
    foo = "env_value1"
    bar = "env_value2"
  }
}

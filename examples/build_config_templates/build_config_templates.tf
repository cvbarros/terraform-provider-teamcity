# This configuration sample shows how to manage Build Configuration Templates
provider "teamcity" {
  address  = var.teamcity_url
  username = var.teamcity_username
  password = var.teamcity_password
}

resource "teamcity_project" "nocode" {
  name = "Samples - Build Config Templates Project"
}

resource "teamcity_build_config" "build_config" {
  project_id  = teamcity_project.nocode.id
  name        = "Main Configuration"
  description = "Main Configuration Description"

  templates = [teamcity_build_config.template1.id, teamcity_build_config.template2.id]
}

resource "teamcity_build_config" "template1" {
  project_id = teamcity_project.nocode.id
  name       = "Template 1"
  # Description is not supported with build config templates! See https://youtrack.jetbrains.com/issue/TW-63617
  is_template = true
}

resource "teamcity_build_config" "template2" {
  project_id = teamcity_project.nocode.id
  name       = "Template 2"
  # Description is not supported with build config templates! See https://youtrack.jetbrains.com/issue/TW-63617
  is_template = true
}

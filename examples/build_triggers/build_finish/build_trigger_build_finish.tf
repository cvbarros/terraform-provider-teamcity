# This configuration sample shows how to setup a build configuration
# to be triggered after another given build finishes.
provider "teamcity" {
  address  = var.teamcity_url
  username = var.teamcity_username
  password = var.teamcity_password
}

resource "teamcity_project" "project" {
  name = "Samples - Build Trigger - Build Finish Project"
}

resource "teamcity_vcs_root_git" "vcs" {
  name       = "Application"
  project_id = teamcity_project.project.id

  fetch_url      = "https://github.com/cvbarros/go-teamcity"
  default_branch = "refs/head/master"
}

resource "teamcity_build_config" "source" {
  project_id = teamcity_project.project.id
  name       = "First Build Config"

  step {
    type = "cmd_line"
    file = "build.sh"
    args = "-t build"
  }

  vcs_root {
    id             = teamcity_vcs_root_git.vcs.id
    checkout_rules = ["+:*"]
  }
}

resource "teamcity_build_config" "triggered" {
  project_id  = teamcity_project.project.id
  name        = "Triggered Build"
  description = "Build triggered when 'First Build Config' is finished"

  step {
    type = "cmd_line"
    file = "build.sh"
    args = "-t release"
  }
}

resource "teamcity_build_trigger_build_finish" "finish_trigger" {
  build_config_id        = teamcity_build_config.triggered.id
  source_build_config_id = teamcity_build_config.source.id

  #Optional, defaults to false
  after_successful_only = true

  #Optional
  branch_filter = ["master", "feature"]
}

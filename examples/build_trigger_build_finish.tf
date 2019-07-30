resource "teamcity_project" "project" {
  name = "Go TeamCity SDK"
}

resource "teamcity_vcs_root_git" "vcs" {
  name       = "Application"
  project_id = "${teamcity_project.project.id}"

  url    = "https://github.com/cvbarros/go-teamcity"
  branch = "refs/head/master"
}

resource "teamcity_buildconfiguration" "build_release" {
  project_id          = "${teamcity_project.project.id}"
  name                = "Build Release"

  step {
    type = "command_line"
    file = "build.sh"
    args = "-t buildrelease"
  }

  vcs_root {
    id             = "${teamcity_vcs_root_git.vcs.id}"
    checkout_rules = ["+:*"]
  }
}

resource "teamcity_buildconfiguration" "triggered_build" {
  project_id          = "${teamcity_project.project.id}"
  name                = "Triggered Build"
  description         = "Build triggered when 'Build Release' is finished"

  step {
    type = "command_line"
    file = "build.sh"
    args = "-t release"
  }
}

resource "teamcity_build_trigger_build_finish" "buildrelease_finish_trigger" {
    build_config_id = "${teamcity_buildconfiguration.build_release.id}"

    #Optional, defaults to false
    after_successful_only = true

    #Optional
    branch_filter = ["master", "feature"]
}

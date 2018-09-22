resource "teamcity_project" "nocode" {
  name = "No Code"
}

resource "teamcity_vcs_root_git" "nocode_vcs" {
  name       = "Application"
  project_id = "${teamcity_project.nocode.id}"

  url    = "https://github.com/kelseyhightower/nocode"
  branch = "refs/head/master"
}

resource "teamcity_buildconfiguration" "nocode_build_release" {
  project_id          = "${teamcity_project.nocode.id}"
  name                = "Build Release"
  description         = "Master build with \"BuildRelease\" hook"
  build_number_format = "0.0.%build.counter%"
  artifact_paths      = [""]

  options {
    status_widget         = false
    detect_hanging        = true
    allow_personal_builds = true
  }

  step {
    type = "powershell"
    file = "build.ps1"
    args = "-Target buildrelease -Verbosity %verbosity%"
  }

  vcs_root {
    id             = "${teamcity_vcs_root_git.nocode_vcs}"
    checkout_rules = ["+:*"]
  }
}

resource "teamcity_buildconfiguration" "nocode_triggered_build" {
  project_id          = "${teamcity_project.nocode.id}"
  name                = "Triggered Build"
  description         = "Build triggered when 'Build Release' is finished"
  build_number_format = "0.0.%build.counter%"
  artifact_paths      = [""]

  options {
    status_widget         = false
    detect_hanging        = true
    allow_personal_builds = true
  }

  step {
    type = "powershell"
    file = "build.ps1"
    args = "-Target buildrelease -Verbosity %verbosity%"
  }

  vcs_root {
    id             = "${teamcity_vcs_root_git.nocode_vcs}"
    checkout_rules = ["+:*"]
  }
}

resource "teamcity_build_trigger_build_finish" "buildrelease_finish_trigger" {
    build_config_id = "${teamcity_buildconfiguration.nocode_build_release.id}"

    #Optional, defaults to false
    after_successful_only = true

    #Optional
    branch_filter = ["master", "feature"]
}
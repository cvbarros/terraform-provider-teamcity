#This sample illustrates a simple complete deployment pipeline using Powershell
provider "teamcity" {
  address  = "http://localhost:8112"
  username = "admin"
  password = "admin"
}

resource "teamcity_project" "nocode" {
  name = "No Code"
}

resource "teamcity_vcs_root_git" "nocode_vcs" {
  name       = "Application"
  project_id = "${teamcity_project.nocode.id}"

  url    = "https://github.com/kelseyhightower/nocode"
  branch = "refs/head/master"
}

resource "teamcity_buildconfiguration" "nocode_pullrequest" {
  project_id          = "${teamcity_project.nocode.id}"
  name                = "Pull Request"
  description         = "Inspection Build with \"Pull-Request\" hook"
  build_number_format = "0.0.%build.counter%"
  artifact_paths      = [""]

  options {
    status_widget         = false
    detect_hanging        = true
    allow_personal_builds = true
  }

  step {
    type = "powershell"
  }

  vcs_root {
    id             = "${teamcity_vcs_root_git.nocode_vcs}"
    checkout_rules = ["+:*"]
  }

  env_params    = {}
  sys_params    = {}
  config_params = {}
  feature       = {}
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

resource "teamcity_buildconfiguration" "nocode_release_testing" {
  project_id          = "${teamcity_project.nocode.id}"
  name                = "Deploy To Testing"
  description         = "Perform a deployment to Testing environment"
  build_number_format = "0.0.%build.counter%"

  step {
    type = "powershell"
    file = "build.ps1"
    args = "-Target release -Verbosity %verbosity%"
  }
}

resource "teamcity_build_trigger" "nocode_vcs_trigger" {
  build_config_id = "${teamcity_buildconfiguration.nocode_pullrequest}"

  type = "vcs"

  rules = ["+:*"]
}

resource "teamcity_snapshot_dependency" "nocode_release_testing" {
  build_config_id        = "${teamcity_buildconfiguration.nocode_release_testing}"
  source_build_config_id = "${teamcity_buildconfiguration.nocode_build_release}"
}

resource "teamcity_agent_requirement" "env_testing" {
  build_config_id = "${teamcity_buildconfiguration.nocode_release_testing.id}"
  condition       = "equals"
  name            = "environment"
  value           = "testing"
}

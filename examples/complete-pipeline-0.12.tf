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

  default_branch = "refs/head/master"
  fetch_url = "https://github.com/kelseyhightower/nocode"
}
resource "teamcity_build_config" "nocode_pullrequest" {
  project_id          = "${teamcity_project.nocode.id}"
  name                = "Pull Request"
  description         = "Inspection Build with \"Pull-Request\" hook"
  
  settings {
    build_number_format = "0.0.%build.counter%"
    artifact_paths      = [""]
    status_widget         = false
    detect_hanging        = true
    allow_personal_builds = true
  }

  step {
    type = "powershell"
    code = "Write-Host Hello"
  }

  vcs_root {
    id             = "${teamcity_vcs_root_git.nocode_vcs.id}"
    checkout_rules = ["+:*"]
  }

  env_params    = {}
  sys_params    = {}
  config_params = {}
}

resource "teamcity_build_config" "nocode_build_release" {
  project_id          = "${teamcity_project.nocode.id}"
  name                = "Build Release"
  description         = "Master build with \"BuildRelease\" hook"

  settings {
    build_number_format = "0.0.%build.counter%"
    artifact_paths      = [""]
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
    id             = "${teamcity_vcs_root_git.nocode_vcs.id}"
    checkout_rules = ["+:*"]
  }
}

resource "teamcity_build_config" "nocode_release_testing" {
  project_id          = "${teamcity_project.nocode.id}"
  name                = "Deploy To Testing"
  description         = "Perform a deployment to Testing environment"

  settings {
    build_number_format = "0.0.%build.counter%"
    artifact_paths      = [""]
  }

  step {
    type = "powershell"
    file = "build.ps1"
    args = "-Target release -Verbosity %verbosity%"
  }
}

resource "teamcity_build_trigger_vcs" "nocode_vcs_trigger" {
  build_config_id = "${teamcity_build_config.nocode_pullrequest.id}"

  rules = ["+:*"]
}

resource "teamcity_snapshot_dependency" "nocode_release_testing" {
  build_config_id        = "${teamcity_build_config.nocode_release_testing.id}"
  source_build_config_id = "${teamcity_build_config.nocode_build_release.id}"
}

resource "teamcity_agent_requirement" "env_testing" {
  build_config_id = "${teamcity_build_config.nocode_release_testing.id}"
  condition       = "equals"
  name            = "environment"
  value           = "testing"
}
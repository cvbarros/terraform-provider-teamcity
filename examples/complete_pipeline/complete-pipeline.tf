#This sample illustrates a simple complete deployment pipeline using Powershell steps
provider "teamcity" {
  address  = var.teamcity_url
  username = var.teamcity_username
  password = var.teamcity_password
}

resource "teamcity_project" "project" {
  name = "Sample - Pipeline Project"
}

resource "teamcity_vcs_root_git" "vcs_git" {
  name       = "Application"
  project_id = teamcity_project.project.id

  fetch_url      = "https://github.com/cvbarros/go-teamcity"
  default_branch = "refs/head/master"

  // BE CAREFUL - THIS IS STORED IN PLAIN TEXT ON STATE FILE!
  // Refer to: https://www.terraform.io/docs/state/sensitive-data.html
  auth {
    type     = "userpass"
    username = var.github_username
    password = var.github_password
  }
}

resource "teamcity_build_config" "pullrequest" {
  project_id  = teamcity_project.project.id
  name        = "Pull Request"
  description = "Inspection Build with \"Pull-Request\" hook"

  settings {
    build_number_format   = "0.0.%build.counter%"
    status_widget         = false
    detect_hanging        = true
    allow_personal_builds = true
  }

  step {
    type = "powershell"
    file = "build.ps1"
    args = "-Target pullrequest"
  }

  vcs_root {
    id             = teamcity_vcs_root_git.vcs_git.id
    checkout_rules = ["+:*"]
  }
}

resource "teamcity_feature_commit_status_publisher" "github_status" {
  build_config_id = teamcity_build_config.pullrequest.id
  publisher       = "github"

  // BE CAREFUL - THIS IS STORED IN PLAIN TEXT ON STATE FILE!
  // Refer to: https://www.terraform.io/docs/state/sensitive-data.html
  github {
    auth_type = "password"
    username  = var.github_username
    password  = var.github_password
  }
}

resource "teamcity_build_config" "build_release" {
  project_id  = teamcity_project.project.id
  name        = "Build Release"
  description = "Master build with \"BuildRelease\" hook"

  settings {
    build_number_format   = "0.0.%build.counter%"
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
    id             = teamcity_vcs_root_git.vcs_git.id
    checkout_rules = ["+:*"]
  }
}

resource "teamcity_build_config" "deploy_testing" {
  project_id  = teamcity_project.project.id
  name        = "Deploy To Testing"
  description = "Perform a deployment to Testing environment"

  step {
    type = "powershell"
    file = "build.ps1"
    args = "-Target release -Verbosity %verbosity%"
  }
}

#Trigger source build on VCS changes
resource "teamcity_build_trigger_vcs" "vcs_trigger" {
  build_config_id = teamcity_build_config.pullrequest.id

  rules = ["+:*"]
}

# Snapshot dependency
resource "teamcity_snapshot_dependency" "dep_release_testing" {
  build_config_id        = teamcity_build_config.deploy_testing.id
  source_build_config_id = teamcity_build_config.build_release.id
}

# Ensures this build only runs on agents that have environment = testing
resource "teamcity_agent_requirement" "agent_req_testing" {
  build_config_id = teamcity_build_config.deploy_testing.id
  condition       = "equals"
  name            = "environment"
  value           = "testing"
}

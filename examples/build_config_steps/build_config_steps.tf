# This configuration sample shows how to manage Build Configuration Settings
provider "teamcity" {
  address  = var.teamcity_url
  username = var.teamcity_username
  password = var.teamcity_password

  version = "0.6.0"
}

resource "teamcity_project" "project" {
  name = "Samples - Build Config Steps Project"
}

resource "teamcity_build_config" "nocode_config" {
  name        = "SimpleConfig"
  description = "Configuration to showcase build configuration steps"
  project_id  = teamcity_project.project.id

  step {
    name = "First Step"
    type = "cmd_line"
    code = "echo"
  }

  step {
    name = "Second Step - File"
    type = "powershell"
    file = "build.ps1"
    args = "-Target inspection"
  }

  step {
    name = "Third Step - Code"
    type = "powershell"
    code = "Get-Date"
  }
}


resource "teamcity_project" "nocode" {
  name = "No Code"
}

resource "teamcity_vcs_root_git" "nocode_vcs" {
  name       = "Application"
  project_id = "${teamcity_project.nocode.id}"

  url    = "https://github.com/kelseyhightower/nocode"
  branch = "refs/head/master"
}
resource "teamcity_buildconfiguration" "nocode_triggered_build" {
  project_id          = "${teamcity_project.nocode.id}"
  name                = "Triggered Build"
  description         = "Build triggered on schedules"
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

resource "teamcity_build_trigger_schedule" "buildrelease_schedule_trigger" {
    build_config_id = "${teamcity_buildconfiguration.nocode_build_release.id}"

    #daily, weekly (cron yet not supported)
    schedule = "daily"

    #Use values from TeamCity UI without abbreviation or (UTC-+X), like: America/Sao Paulo. Defaults to "SERVER", which uses the SERVER's timezone
    timezone = "America/Sao Paulo"

    #Hour of the day that trigger will fire [0..23]
    hour = 12

    #Minute that the trigger will fire [0..59]. Optional, defaults to 0
    minute = 37

    #Weekday for weekly triggers. Use only if schedule = "weekly". Possible values are the weekday names in english, like: "Monday", "Tuesday"...
    weekday = "Saturday"

    #Optional, only trigger on changes that match the rules. If none used, will trigger for any change.
    rules = ["+:*", "-:*.md"]
}
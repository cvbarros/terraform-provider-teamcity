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

  #Defaults to true on TeamCity UI - Queued build can be replaced with an already started build or a more recent queued build
  queue_optimization = false

  #Triggers on all compatible agents - Default: false
  on_all_compatible_agents = true

  #Trigger only if watched build changes - Default: false
  with_pending_changes_only = true

  #Promote watched build if there is a dependency (snapshot or artifact) on its build configuration. Default: true
  promote_watched_build = false

  #Delete all files in checkout directory before the build - Default: false
  enforce_clean_checkout = true
  
  #Delete all files in checkout directory before the build also for snapshot dependencies. Default: false
  enforce_clean_checkout_dependencies = true

  #Trigger only if a given build configuration has pending changes. Default: false. If set to true, watched_build_config_id must be set
  only_if_watched_changes = true

  #Configures the watched build for this trigger
  watched_build_config_id = "${teamcity_build_config.watched.id}"
  
  #Specify which version of the watched build should be considered. "lastFinished", "lastPinned", "lastSuccessful", "buildTag"
  revision = "lastFinished"

  #Used with revision = "buildTag", to specify the tag/branch for the watched build to be considered
  watched_branch = "unstable"
}
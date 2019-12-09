# This configuration sample shows how to setup a build configuration
# to be triggered on a schedule
provider "teamcity" {
  address  = var.teamcity_url
  username = var.teamcity_username
  password = var.teamcity_password
}

resource "teamcity_project" "project" {
  name = "Samples - Build Trigger - Schedule Project"
}

resource "teamcity_build_config" "triggered" {
  project_id  = teamcity_project.project.id
  name        = "Triggered Build"
  description = "Build triggered on schedule"
}

resource "teamcity_build_config" "watched" {
  project_id  = teamcity_project.project.id
  name        = "Watched Build"
  description = "Trigger will only fire if watched build has pending changes"
}

resource "teamcity_build_trigger_schedule" "schedule_trigger" {
  build_config_id = teamcity_build_config.triggered.id

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
  promote_watched_build = true

  #Delete all files in checkout directory before the build - Default: false
  enforce_clean_checkout = true

  #Delete all files in checkout directory before the build also for snapshot dependencies. Default: false
  enforce_clean_checkout_dependencies = true

  #Trigger only if a given build configuration has pending changes. Default: false. If set to true, watched_build_config_id must be set
  only_if_watched_changes = true

  #Configures the watched build for this trigger
  watched_build_config_id = teamcity_build_config.watched.id

  #Specify which version of the watched build should be considered. "lastFinished", "lastPinned", "lastSuccessful", "buildTag"
  revision = "lastFinished"

  #Used with revision = "buildTag", to specify the tag/branch for the watched build to be considered
  watched_branch = "unstable"
}

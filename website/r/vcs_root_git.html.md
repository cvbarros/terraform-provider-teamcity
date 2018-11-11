---
layout: teamcity
page_title: "TeamCity: teamcity_vcs_root_git resource"
sidebar_current: "docs-teamcity-vcs-root-git"
description: |-
  Manages TeamCity Git VCS Roots
---

# teamcity\_vcs\_root\_git

The Git VCS Root resource allows managing VCS Roots with type `Git`.

~> **WARNING:** When using `userpass` or `ssh` with `customKey` authentication, credentials will be persisted in plain-text to the state file. Seek using other forms of authentication to private Git repositories or treat state files treated as sensitive and protected accordingly.

## Example Usage

```hcl
resource "teamcity_project" "project" {
  name = "Project"
}

resource "teamcity_vcs_root_git" "vcsroot" {
  name       = "Application"
  project_id = "${teamcity_project.project.id}"

  fetch_url = "https://github.com/cvbarros/terraform-provider-teamcity"

  default_branch = "refs/head/master"

  branches = [
    "+:refs/(pull/*)/head",
    "refs/heads/develop",
  ]
  enable_branch_spec_tags = false
  usernameStyle = "userid"
  submoduleCheckout = true

  # Auth block configures the authentication to Git VCS
  auth {
    type     = "userpass"
    username = "admin"

    # Sensitive -> always updated on apply because TeamCity doesn't return passwords
    password = "<<<secret>>>"
  }

  # Configure agent settings
  agent {
    git_path = "/usr/bin/git"
    clean_policy = "branch_change"
    clean_files_policy = "untracked"
    use_mirrors = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `name`: (Required) Specifies the name which the VCS Root will be have. TeamCity [automatically generates](https://confluence.jetbrains.com/display/TCD18/Identifier) a friendly `ID`  based on name. If duplicate names are found within a same project, TeamCity will append a number to the end of the `ID`. It is better to avoid duplicating VCS Root names in the scope of the same project.

* `project_id`: (Required) ID of the project under which this VCS Root will be created. Use `_Root` to create a top-level VCS Root.

* `fetch_url`: (Required) URL used to pull source code for this VCS. For HTTP, prefix with http(s)://. For SSH, use user@server.com.

* `push_url`: (Optional) URL used to push source code for this VCS. Assumes the same as `fetch_url`if not specified.

* `default_branch`: (Required) Branch specification for the default branch to pull/push from/to and watch for changes. Ex: `refs/head/master`.

* `branches`: (Optional) A list of branches to monitor besides the default with a set of rules in the form of +|-:branch_name (with the optional * placeholder).
  
* `enable_branch_spec_tags`: (Optional) If true, tags can be used in the branch specification.

* `submodule_checkout`: (Optional) If true, submodules will be checkout out along with the main repository. Defaults to `true`.

* `username_style`: (Optional) Defines a way TeamCity binds VCS changes to the user. Changing username style will affect only newly collected changes. Old changes will continue to be stored with the style that was active at the time of collecting changes. Allowed values: `userid`, `author_name`, `author_email`, `author_full`.

The `auth` block is used to manage authentication configuration. If not specified, defaults to anonymous auth.

* `type`: (Required) Authentication type to use. Can be `userpass`, `ssh` or `anonymous`.

* `username`: (Optional) Username to connect if using `userpass`.

* `password`: (Optional) Password if using 'userpass' auth. Private key passphrase if using `uploadedKey` or `customKey`. Required if not using `anonymous` auth.

* `ssh_type`: (Optional) IF using `ssh` auth, this field specifies how the SSH key will be sourced. `uploadedKey` refers to a previously uploaded SSH Key to a project in the hierarchy. `customKey` is a key already provisioned on the server. `defaultKey` uses the keys available on the file system in the default locations used by common ssh tools.

* `key_spec`: (Optional) For `customKey` refers to the path on the server to a private key. For `uploadedKey`, corresponds to the name of the SSH Key uploaded into the project. Required if using `customKey` or `uploadedKey`.

The `agent` block is used to tweak agent settings.

* `git_path`: (Optional) The path to a git executable on the agent. If blank, the location set up in TEAMCITY_GIT_PATH environment variable is used by the server.

* `clean_policy`: (Optional) This option specifies when the "git clean" command is run on the agent. Allowed values are `branch_change`, `always`, `never`.

* `clean_files_policy`: (Optional) This option specifies which files will be removed when "git clean" command is run on agent. Allowed values are `untracked`, `ignored_only`, `non_ignored_only`.

* `use_mirrors`: (Optional) If true, TeamCity creates a separate clone of the repository on each agent and uses it in the checkout directory via git alternates.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:

* `id`: The auto-generated ID of the VCS Root.

## Import
Git VCS Roots can be imported using their ID, e.g.

```
$ terraform import teamcity_vcs_root_git.vcsroot Project_Application
```

# Data Source: teamcity_project

Retrieves information about an existing TeamCity Project

## Example Usage

```hcl
data "teamcity_project" "by-id" {
  project_id = "ABC123"
}

data "teamcity_project" "by-name" {
  name = "Project Name"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The Name of the Project.

* `project_id` - (Optional) The [Identifier](https://confluence.jetbrains.com/display/TCD18/Identifier) assigned to this Project

~> **Note:** At least one of `name` or `project_id` must be specified. 


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `description` - A description assigned to this Project.

* `parent_project_id` - The ID of the Parent Project this Project is nested under.

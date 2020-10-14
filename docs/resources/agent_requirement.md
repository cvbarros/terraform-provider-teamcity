# teamcity_agent_requirement

The Agent Requirement resource allows managing Agent Requirements for build configurations.

## Example Usage

```hcl
resource "teamcity_project" "project" {
  name = "Project"
}

resource "teamcity_build_config" "build" {
  name       = "source"
  project_id = teamcity_project.project.id
}

resource "teamcity_agent_requirement" "requirement" {
  build_config_id = teamcity_build_config.build.id
  name            = "teamcity.agent.env"
  condition       = "equals"
  value           = "testing"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of parameter provided by the agent that represents the requirement.

* `build_config_id` - (Required) ID of the build configuration which this requirement will be configured.

* `condition` - (Required) A string operator to match the `name` variable to the `value` which will satistfy the requirement. Possible [values are documented](https://godoc.org/github.com/cvbarros/go-teamcity/pkg/teamcity#pkg-variables) in the TeamCity API.

* `value` - (Required) Right-side operand of the condition to be checked against the parameter.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The auto-generated ID of the agent requirement.

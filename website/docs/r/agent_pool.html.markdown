---
subcategory: "Agent Pools"
layout: teamcity
page_title: "TeamCity: Resource - teamcity_agent_pool"
description: |-
  Manages an Agent Pool
---

# teamcity_agent_pool

Manages an Agent Pool

## Example Usage

```hcl
resource "teamcity_agent_pool" "example" {
  name = "Example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Agent Pool.

* `max_agents` - (Optional) Specifies the maximum number of Build Agents which can be associated with this Agent Pool. Defaults to `-1` which means `unlimited`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The auto-generated ID of the Agent Pool.

## Import

Agent Pools can be imported using their ID, e.g.

```
$ terraform import teamcity_agent_pool.example 28
```

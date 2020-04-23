---
subcategory: "Agent Pools"
layout: "teamcity"
page_title: "TeamCity: Data Source - teamcity_agent_pool"
description: |-
  Retrieves information about an existing TeamCity Agent Pool
---

# Data Source: teamcity_agent_pool

Retrieves information about an existing TeamCity Agent Pool

## Example Usage

```hcl
data "teamcity_project" "by-name" {
  name = "Default"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Name of the Agent Pool.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `max_agents` - The maximum number of agents in this pool. If set to `unlimited` this'll be `-1`.

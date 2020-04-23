---
subcategory: "Agent Pools"
layout: teamcity
page_title: "TeamCity: Resource - teamcity_agent_pool_project_assignment"
description: |-
  Manages an Assignment between a Agent Pool and a Project
---

# teamcity_agent_pool_project_assignment

Manages an Assignment between a Agent Pool and a Project

## Example Usage

```hcl
resource "teamcity_project" "example" {
  name = "Example"
}

resource "teamcity_agent_pool" "example" {
  name = "Example"
}

resource "teamcity_agent_pool_project_assignment" "test" {
  agent_pool_id                 = teamcity_agent_pool.example.id
  project_id                    = teamcity_project.example.id
  disassociate_from_other_pools = true
}
```

## Argument Reference

The following arguments are supported:

* `agent_pool_id` - (Required) Specifies the ID of the Agent Pool. Changing this forces a new resource to be created.

* `project_id` - (Required) Specifies the ID of the Project. Changing this forces a new resource to be created.

---

* `disassociate_from_other_pools` - (Optional) When creating this assignment, should all other assignments be removed?

~> **Note:** TeamCity requires that a Project always have at least one Agent Pool associated with it - as such when destroying an assignment Terraform will re-assign back to the "Default" node pool if required.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The auto-generated ID of the Agent Pool - Project Assignment.

## Import

Agent Pools - Project associations can be imported using their ID, e.g.

```
$ terraform import teamcity_agent_pool.example 28|Project12
```

-> **Note:** This is a Terraform specific ID in the format `AgentPoolID|ProjectID`

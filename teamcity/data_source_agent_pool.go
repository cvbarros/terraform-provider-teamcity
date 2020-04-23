package teamcity

import (
	"fmt"

	api "github.com/cvbarros/go-teamcity/teamcity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAgentPool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAgentPoolRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"max_agents": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceAgentPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	name := d.Get("name").(string)
	agentPool, err := client.AgentPools.GetByName(name)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", agentPool.Id))
	d.Set("name", agentPool.Name)
	maxAgents := -1
	if agentPool.MaxAgents != nil {
		maxAgents = *agentPool.MaxAgents
	}
	d.Set("max_agents", maxAgents)

	return nil
}

package teamcity

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	api "github.com/cvbarros/go-teamcity/teamcity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAgentPool() *schema.Resource {
	return &schema.Resource{
		// the TC Docs say this supports update - but the documented API doesn't work
		// (returns 405 Method Not Allowed) so for the moment this can't support update
		Create: resourceAgentPoolCreate,
		Read:   resourceAgentPoolRead,
		Delete: resourceAgentPoolDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"max_agents": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  -1,
			},
		},
	}
}

func resourceAgentPoolCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	agentPool := api.CreateAgentPool{
		Name: d.Get("name").(string),
	}
	if v := d.Get("max_agents").(int); v >= 0 {
		agentPool.MaxAgents = &v
	}

	createdAgentPool, err := client.AgentPools.Create(agentPool)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", createdAgentPool.Id))

	return resourceAgentPoolRead(d, client)
}

func resourceAgentPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	agentPool, err := client.AgentPools.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Printf("[DEBUG] Agent Pool not found - removing from state!")
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("name", agentPool.Name)
	maxAgents := -1
	if agentPool.MaxAgents != nil {
		maxAgents = *agentPool.MaxAgents
	}
	d.Set("max_agents", maxAgents)

	return nil
}

func resourceAgentPoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	return client.AgentPools.Delete(id)
}

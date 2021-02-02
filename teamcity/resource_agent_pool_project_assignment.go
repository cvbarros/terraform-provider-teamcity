package teamcity

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	api "github.com/cvbarros/go-teamcity/teamcity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAgentPoolProjectAssignment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAgentPoolProjectAssignmentCreate,
		Read:   resourceAgentPoolProjectAssignmentRead,
		Delete: resourceAgentPoolProjectAssignmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"agent_pool_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"disassociate_from_other_pools": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
		},
	}
}

func resourceAgentPoolProjectAssignmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	agentPoolId := d.Get("agent_pool_id").(int)
	projectId := d.Get("project_id").(string)

	if err := client.AgentPools.AssignProject(agentPoolId, projectId); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d|%s", agentPoolId, projectId))

	disassociateFromOtherPools := d.Get("disassociate_from_other_pools").(bool)
	if disassociateFromOtherPools {
		agentPoolsForProject, err := client.AgentPools.ListForProject(projectId)
		if err != nil {
			return fmt.Errorf("error listing agent pools for project %q: %s", projectId, err)
		}

		for _, pool := range agentPoolsForProject.AgentPools {
			if pool.Id == agentPoolId {
				continue
			}

			log.Printf("[DEBUG] Removing association between Project %q and Agent Pool %d", projectId, pool.Id)
			if err := client.AgentPools.UnassignProject(pool.Id, projectId); err != nil {
				return fmt.Errorf("Error removing association between Project %q and Agent Pool %d: %+v", projectId, pool.Id, err)
			}
		}
	}

	return resourceAgentPoolProjectAssignmentRead(d, client)
}

func resourceAgentPoolProjectAssignmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	id, err := ParseAgentPoolProjectAssignmentID(d.Id())
	if err != nil {
		return err
	}

	agentPool, err := client.AgentPools.GetByID(id.AgentPoolId)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Printf("[DEBUG] Agent Pool not found, so assignment can't - removing from state!")
			d.SetId("")
			return nil
		}

		return err
	}

	projectExists := false
	if agentPool.Projects != nil {
		for _, project := range agentPool.Projects.Project {
			if project.ID == id.ProjectId {
				projectExists = true
				break
			}
		}
	}

	if !projectExists {
		log.Printf("[DEBUG] Agent Pool <-> Project Assignment not found - removing from state!")
		d.SetId("")
		return nil
	}

	d.Set("agent_pool_id", agentPool.Id)
	d.Set("project_id", id.ProjectId)

	return nil
}

func resourceAgentPoolProjectAssignmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	id, err := ParseAgentPoolProjectAssignmentID(d.Id())
	if err != nil {
		return err
	}

	// TeamCity requires that a Project is in at least one Agent Pool
	// as such, if this is the only Agent Pool Assignment, force-move it back to the "_Root" project
	// since that's guaranteed to be around
	agentPoolsForProject, err := client.AgentPools.ListForProject(id.ProjectId)
	if err != nil {
		return fmt.Errorf("error retrieving agent pools for project: %s", err)
	}

	needsReassigningToRootProject := false
	if len(agentPoolsForProject.AgentPools) == 1 && agentPoolsForProject.AgentPools[0].Id == id.AgentPoolId {
		needsReassigningToRootProject = true
	}

	if needsReassigningToRootProject {
		log.Printf("[DEBUG] TeamCity requires that a Build Configuration exists in at least one Agent Pool")
		log.Printf("[DEBUG] Since this has no other assignments, adding one to the Default Agent Pool")
		defaultAgentPoolId := 0
		if err := client.AgentPools.AssignProject(defaultAgentPoolId, id.ProjectId); err != nil {
			return fmt.Errorf("error assigning project to Default Agent Pool: %s", err)
		}
	}

	if err := client.AgentPools.UnassignProject(id.AgentPoolId, id.ProjectId); err != nil {
		return fmt.Errorf("error unassigning project %q from pool %d: %s", id.ProjectId, id.AgentPoolId, err)
	}

	return nil
}

type agentPoolProjectAssignmentId struct {
	AgentPoolId int
	ProjectId   string
}

func ParseAgentPoolProjectAssignmentID(input string) (*agentPoolProjectAssignmentId, error) {
	// format: "AgentPoolID|ProjectID"
	segments := strings.Split(input, "|")
	if len(segments) != 2 {
		return nil, fmt.Errorf("Expected 2 segments but got %d", len(segments))
	}

	agentPoolId, err := strconv.Atoi(segments[0])
	if err != nil {
		return nil, err
	}

	id := agentPoolProjectAssignmentId{
		AgentPoolId: agentPoolId,
		ProjectId:   segments[1],
	}
	return &id, nil
}

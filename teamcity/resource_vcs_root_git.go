package teamcity

import (
	api "github.com/cvbarros/go-teamcity-sdk/pkg/teamcity"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVcsRootGit() *schema.Resource {
	return &schema.Resource{
		Create: resourceVcsRootGitCreate,
		Read:   resourceVcsRootGitRead,
		Update: resourceVcsRootGitUpdate,
		Delete: resourceVcsRootGitDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"default_branch": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceVcsRootGitCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	projectID := d.Get("project_id").(string)

	newVcsRoot := &api.VcsRoot{
		Name:    d.Get("name").(string),
		VcsName: api.VcsNames.Git,
		Project: &api.ProjectReference{
			ID: projectID,
		},
		Properties: api.NewProperties(
			&api.Property{
				Name:  "url",
				Value: d.Get("repo_url").(string),
			},
			&api.Property{
				Name:  "branch",
				Value: d.Get("default_branch").(string),
			},
		),
	}

	created, err := client.VcsRoots.Create(projectID, newVcsRoot)
	if err != nil {
		return err
	}

	d.MarkNewResource()
	d.SetId(created.ID)
	return nil
}

func resourceVcsRootGitRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVcsRootGitUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVcsRootGitDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	return client.VcsRoots.Delete(d.Id())
}

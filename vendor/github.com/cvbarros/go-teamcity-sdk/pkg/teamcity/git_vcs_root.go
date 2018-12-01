package teamcity

import (
	"encoding/json"
	"errors"
	"fmt"
)

//GitVcsRoot is a VCS Root of type Git, strongly-typed model.
type GitVcsRoot struct {
	Options *GitVcsRootOptions

	// id
	ID string `json:"id,omitempty" xml:"id"`

	// project
	Project *ProjectReference `json:"project,omitempty"`

	modificationCheckInterval *int32
	name                      string
	vcsRootJSON               *vcsRootJSON
	properties                *Properties
}

//NewGitVcsRoot returns a VCS Root instance that connects to Git VCS.
func NewGitVcsRoot(projectID string, name string, opts *GitVcsRootOptions) (*GitVcsRoot, error) {
	if projectID == "" {
		return nil, errors.New("projectID is required")
	}
	if name == "" {
		return nil, errors.New("name is required")
	}
	if opts == nil {
		return nil, errors.New("opts is required")
	}
	return &GitVcsRoot{
		name: name,
		Project: &ProjectReference{
			ID: projectID,
		},
		Options: opts,
		vcsRootJSON: &vcsRootJSON{
			Project: &ProjectReference{
				ID: projectID,
			},
			Name:       name,
			VcsName:    VcsNames.Git,
			Properties: opts.properties(),
		},
	}, nil
}

//GetID returns the ID of this VCS Root.
func (d *GitVcsRoot) GetID() string {
	return d.ID
}

//VcsName returns the type of VCS Root. See VcsNames
func (d *GitVcsRoot) VcsName() string {
	return VcsNames.Git
}

//Name returns the name of VCS Root.
func (d *GitVcsRoot) Name() string {
	return d.name
}

//SetName changes the name of VCS Root.
func (d *GitVcsRoot) SetName(name string) {
	d.name = name
}

//ModificationCheckInterval returns how often TeamCity polls the VCS repository for VCS changes (in seconds).
func (d *GitVcsRoot) ModificationCheckInterval() *int32 {
	return d.modificationCheckInterval
}

//SetModificationCheckInterval specifies how often TeamCity polls the VCS repository for VCS changes (in seconds).
func (d *GitVcsRoot) SetModificationCheckInterval(seconds int32) {
	d.modificationCheckInterval = &seconds
}

//ProjectID returns the projectID where this VCS Root is defined
func (d *GitVcsRoot) ProjectID() string {
	return d.Project.ID
}

//SetProjectID specifies the project for this VCS Root. When moving VCS Roots between projects, it must not be in use by any other build configurations or sub-projects.
func (d *GitVcsRoot) SetProjectID(id string) {
	d.Project.ID = id
}

//Properties returns the properties for this VCS Root
func (d *GitVcsRoot) Properties() *Properties {
	return d.Options.properties()
}

//MarshalJSON implements JSON serialization for GitVcsRoot
func (d *GitVcsRoot) MarshalJSON() ([]byte, error) {
	out := &vcsRootJSON{
		ID:         d.ID,
		Name:       d.name,
		Project:    d.Project,
		VcsName:    d.VcsName(),
		Properties: d.Options.properties(),
	}

	if d.modificationCheckInterval != nil {
		out.ModificationCheckInterval = *d.modificationCheckInterval
	}

	return json.Marshal(out)
}

//UnmarshalJSON implements JSON deserialization for GitVcsRoot
func (d *GitVcsRoot) UnmarshalJSON(data []byte) error {
	var aux vcsRootJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.VcsName != VcsNames.Git {
		return fmt.Errorf("invalid VcsName %s trying to deserialize into GitVcsRoot entity", aux.VcsName)
	}

	d.name = aux.Name
	d.Project = aux.Project
	if aux.ModificationCheckInterval != 0 {
		d.modificationCheckInterval = NewInt32(aux.ModificationCheckInterval)
	}
	d.ID = aux.ID
	d.properties = NewProperties(aux.Properties.Items...)
	d.Options = d.properties.gitVcsOptions()

	return nil
}

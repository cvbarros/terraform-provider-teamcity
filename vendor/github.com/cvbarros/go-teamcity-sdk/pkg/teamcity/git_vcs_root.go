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

	// ModificationCheckInterval value in seconds to override the global server setting.
	ModificationCheckInterval int32 `json:"modificationCheckInterval,omitempty" xml:"modificationCheckInterval"`

	// name
	Name string `json:"name,omitempty" xml:"name"`

	// project
	Project *ProjectReference `json:"project,omitempty"`

	vcsRootJSON *vcsRootJSON
	properties  *Properties
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
		Name: name,
		Project: &ProjectReference{
			ID: projectID,
		},
		Options: opts,
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

//Properties returns the properties for this VCS Root
func (d *GitVcsRoot) Properties() *Properties {
	return d.properties
}

//MarshalJSON implements JSON serialization for GitVcsRoot
func (d *GitVcsRoot) MarshalJSON() ([]byte, error) {
	out := &vcsRootJSON{
		ID:   d.ID,
		Name: d.Name,
		ModificationCheckInterval: d.ModificationCheckInterval,
		Project:                   d.Project,
		VcsName:                   d.VcsName(),
		Properties:                d.Options.gitVcsRootProperties(),
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
	d.Name = aux.Name
	d.Project = aux.Project
	d.ModificationCheckInterval = aux.ModificationCheckInterval
	d.ID = aux.ID
	d.properties = NewProperties(aux.Properties.Items...)
	d.Options = d.properties.gitVcsOptions()

	return nil
}

package teamcity

import (
	"encoding/json"
	"errors"
	"fmt"
)

//Dependency is an interface representing a Build dependency, for creating build chains
type Dependency interface {
	ID() string
	Type() string
	SetBuildTypeID(string)
	BuildTypeID() string
	SetSourceBuildTypeID(string)
	SourceBuildTypeID() string
}

type dependencyJSON struct {
	// disabled - Read Only, no effect on post
	Disabled *bool `json:"disabled,omitempty" xml:"disabled"`

	// id
	ID string `json:"id,omitempty" xml:"id"`

	// inherited
	Inherited *bool `json:"inherited,omitempty" xml:"inherited"`

	// Properties are serializable options for this artifact dependency. Do not change this field directly, use the NewArtifactDependency... constructors
	Properties *Properties `json:"properties,omitempty"`

	// source build type
	SourceBuildType *BuildTypeReference `json:"source-buildType,omitempty"`

	// type
	Type string `json:"type,omitempty" xml:"type"`
}

// ArtifactDependency represents a single artifact dependency for a build type
type ArtifactDependency struct {
	dependencyJSON    *dependencyJSON
	buildTypeID       string
	sourceBuildTypeID string

	Options *ArtifactDependencyOptions
}

//ID for this entity
func (s *ArtifactDependency) ID() string {
	return s.dependencyJSON.ID
}

//Type for this entity
func (s *ArtifactDependency) Type() string {
	return "artifact_dependency"
}

//BuildTypeID gets the build type identifier
func (s *ArtifactDependency) BuildTypeID() string {
	return s.buildTypeID
}

//SetBuildTypeID sets the build type identifier
func (s *ArtifactDependency) SetBuildTypeID(id string) {
	s.buildTypeID = id
}

//SourceBuildTypeID gets the source build type identifier
func (s *ArtifactDependency) SourceBuildTypeID() string {
	return s.sourceBuildTypeID
}

//SetSourceBuildTypeID sets the source build type identifier
func (s *ArtifactDependency) SetSourceBuildTypeID(id string) {
	s.sourceBuildTypeID = id
}

//SetDisabled controls whether this dependency is disabled or not
func (s *ArtifactDependency) SetDisabled(disabled bool) {
	s.dependencyJSON.Disabled = NewBool(disabled)
}

//Disabled gets the disabled status for this dependency
func (s *ArtifactDependency) Disabled() bool {
	return *s.dependencyJSON.Disabled
}

// NewArtifactDependency creates a ArtifactDependency with specified options
func NewArtifactDependency(sourceBuildTypeID string, opt *ArtifactDependencyOptions) (*ArtifactDependency, error) {
	if sourceBuildTypeID == "" {
		return nil, errors.New("sourceBuildTypeID is required")
	}

	if opt == nil {
		return nil, errors.New("options must be valid")
	}

	return &ArtifactDependency{
		sourceBuildTypeID: sourceBuildTypeID,
		Options:           opt,
		dependencyJSON: &dependencyJSON{
			SourceBuildType: &BuildTypeReference{ID: sourceBuildTypeID},
			Disabled:        NewFalse(),
			Type:            "artifact_dependency",
			Properties:      opt.artifactDependencyProperties(),
		},
	}, nil
}

//MarshalJSON implements JSON serialization for ArtifactDependency
func (s *ArtifactDependency) MarshalJSON() ([]byte, error) {
	out := &dependencyJSON{
		ID:              s.ID(),
		Type:            s.Type(),
		Disabled:        NewBool(s.Disabled()),
		SourceBuildType: &BuildTypeReference{ID: s.SourceBuildTypeID()},
		Properties:      s.Options.artifactDependencyProperties(),
	}

	return json.Marshal(out)
}

//UnmarshalJSON implements JSON deserialization for ArtifactDependency
func (s *ArtifactDependency) UnmarshalJSON(data []byte) error {
	var aux dependencyJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Type != "artifact_dependency" {
		return fmt.Errorf("invalid type %s trying to deserialize into ArtifactDependency entity", aux.Type)
	}

	if aux.Disabled != nil {
		s.SetDisabled(*aux.Disabled)
	}
	s.dependencyJSON = &aux
	s.SetSourceBuildTypeID(aux.SourceBuildType.ID)
	s.Options = aux.Properties.artifactDepencyOptions()

	return nil
}

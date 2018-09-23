package teamcity

//Dependency is an interface representing a Build dependency, for creating build chains
type Dependency interface {
	ID() string
	Type() string
	SetBuildTypeID(string)
	BuildTypeID() string
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

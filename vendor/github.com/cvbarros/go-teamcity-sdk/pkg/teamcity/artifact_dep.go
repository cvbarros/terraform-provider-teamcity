package teamcity

import "errors"

// ArtifactDependency represents a single artifact dependency for a build type
type ArtifactDependency struct {

	// disabled - Read Only, no effect on post
	Disabled *bool `json:"disabled,omitempty" xml:"disabled"`

	// href
	Href string `json:"href,omitempty" xml:"href"`

	// id
	ID string `json:"id,omitempty" xml:"id"`

	// inherited
	Inherited *bool `json:"inherited,omitempty" xml:"inherited"`

	// Properties are serializable options for this artifact dependency. Do not change this field directly, use the NewArtifactDependency... constructors
	Properties *Properties `json:"properties,omitempty"`

	// source build type
	SourceBuildType *BuildTypeReference `json:"source-buildType,omitempty"`

	// Build type id this dependency belongs to
	BuildTypeID string `json:"-"`

	// type
	Type string `json:"type,omitempty" xml:"type"`
}

//ArtifactDependencies represents a collection of ArtifactDependency
type ArtifactDependencies struct {
	// count
	Count int32 `json:"count,omitempty" xml:"count"`

	// property
	Items []*ArtifactDependency `json:"artifact-dependency"`
}

// NewArtifactDependency creates a ArtifactDependency struct with default options
func NewArtifactDependency(sourceBuildTypeID string, opt *ArtifactDependencyOptions) (*ArtifactDependency, error) {
	if sourceBuildTypeID == "" {
		return nil, errors.New("sourceBuildTypeID is required")
	}

	if opt == nil {
		return nil, errors.New("options must be valid")
	}

	return &ArtifactDependency{
		SourceBuildType: &BuildTypeReference{ID: sourceBuildTypeID},
		Disabled:        NewFalse(),
		Type:            "artifact_dependency",
		Properties:      opt.artifactDependencyProperties(),
	}, nil
}

// func (opt *Arti) properties() *Properties {
// 	var props []*Property

// 	p := NewProperty("run-build-if-dependency-failed", opt.OnFailedDependency)
// 	props = append(props, p)

// 	p = NewProperty("run-build-if-dependency-failed-to-start", opt.OnFailedToStartOrCanceledDependency)
// 	props = append(props, p)

// 	p = NewProperty("run-build-on-the-same-agent", strconv.FormatBool(opt.RunSameAgent))
// 	props = append(props, p)

// 	p = NewProperty("take-started-build-with-same-revisions", strconv.FormatBool(opt.DoNotRunNewBuildIfThereIsASuitable))
// 	props = append(props, p)

// 	p = NewProperty("take-successful-builds-only", strconv.FormatBool(opt.TakeSuccessfulBuildsOnly))
// 	props = append(props, p)

// 	return NewProperties(props...)
// }

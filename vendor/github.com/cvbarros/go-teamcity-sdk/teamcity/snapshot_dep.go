package teamcity

import (
	"strconv"
)

// SnapshotDependency represents a single snapshot dependency for a build type
type SnapshotDependency struct {

	// disabled - Read Only, no effect on post
	Disabled *bool `json:"disabled,omitempty" xml:"disabled"`

	// href
	Href string `json:"href,omitempty" xml:"href"`

	// id
	ID string `json:"id,omitempty" xml:"id"`

	// inherited
	Inherited *bool `json:"inherited,omitempty" xml:"inherited"`

	// Properties are serializable options for this snapshot dependency. Do not change this field directly, use the NewSnapshotDependency... constructors
	Properties *Properties `json:"properties,omitempty"`

	// source build type
	SourceBuildType *BuildTypeReference `json:"source-buildType,omitempty"`

	// Build type id this dependency belongs to
	BuildTypeID string `json:"-"`

	// type
	Type string `json:"type,omitempty" xml:"type"`
}

//SnapshotDependencies represents a collection of SnapshotDependency
type SnapshotDependencies struct {
	// count
	Count int32 `json:"count,omitempty" xml:"count"`

	// property
	Items []*SnapshotDependency `json:"snapshot-dependency"`
}

// SnapshotDependencyOptions represents possible options for when creating a snapshot dependency
type SnapshotDependencyOptions struct { //TODO: Encapsulate those string fields to validate when setting
	OnFailedDependency                  string
	OnFailedToStartOrCanceledDependency string
	RunSameAgent                        bool
	TakeSuccessfulBuildsOnly            bool
	DoNotRunNewBuildIfThereIsASuitable  bool
}

// DefaultSnapshotDependencyOptions are the same options presented by default on Teamcity UI. Do not change this.
var DefaultSnapshotDependencyOptions = &SnapshotDependencyOptions{
	OnFailedDependency:                  "RUN_ADD_PROBLEM",
	OnFailedToStartOrCanceledDependency: "MAKE_FAILED_TO_START",
	RunSameAgent:                        false,
	TakeSuccessfulBuildsOnly:            true,
	DoNotRunNewBuildIfThereIsASuitable:  true,
}

// NewSnapshotDependency created a SnapshotDependency struct with default SnapshotDependencyOptions
func NewSnapshotDependency(sourceBuildTypeID string) *SnapshotDependency {
	return NewSnapshotDependencyWithOptions(sourceBuildTypeID, DefaultSnapshotDependencyOptions)
}

// NewSnapshotDependencyWithOptions creates a SnapshotDependency struct with the provided options
func NewSnapshotDependencyWithOptions(sourceBuildTypeID string, opt *SnapshotDependencyOptions) *SnapshotDependency {
	sourceBuild := &BuildTypeReference{ID: sourceBuildTypeID}
	return &SnapshotDependency{
		SourceBuildType: sourceBuild,
		Properties:      opt.properties(),
		Type:            "snapshot_dependency",
	}
}

func (opt *SnapshotDependencyOptions) properties() *Properties {
	var props []*Property

	p := NewProperty("run-build-if-dependency-failed", opt.OnFailedDependency)
	props = append(props, p)

	p = NewProperty("run-build-if-dependency-failed-to-start", opt.OnFailedToStartOrCanceledDependency)
	props = append(props, p)

	p = NewProperty("run-build-on-the-same-agent", strconv.FormatBool(opt.RunSameAgent))
	props = append(props, p)

	p = NewProperty("take-started-build-with-same-revisions", strconv.FormatBool(opt.DoNotRunNewBuildIfThereIsASuitable))
	props = append(props, p)

	p = NewProperty("take-successful-builds-only", strconv.FormatBool(opt.TakeSuccessfulBuildsOnly))
	props = append(props, p)

	return NewProperties(props...)
}

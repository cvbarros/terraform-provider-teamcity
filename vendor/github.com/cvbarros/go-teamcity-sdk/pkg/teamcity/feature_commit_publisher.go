package teamcity

import (
	"encoding/json"
)

//FeatureCommitStatusPublisherOptions represents options needed to create a commit status publisher build feature
type FeatureCommitStatusPublisherOptions interface {
	Properties() *Properties
}

//FeatureCommitStatusPublisher represents a commit status publisher build feature. Implements BuildFeature interface
type FeatureCommitStatusPublisher struct {
	id          string
	vcsRootID   string
	disabled    bool
	Options     FeatureCommitStatusPublisherOptions
	buildTypeID string

	properties *Properties
}

func (f *FeatureCommitStatusPublisher) ID() string {
	return f.id
}

func (f *FeatureCommitStatusPublisher) SetID(value string) {
	f.id = value
}

func (f *FeatureCommitStatusPublisher) Type() string {
	return "commit-status-publisher"
}

func (f *FeatureCommitStatusPublisher) VcsRootID() string {
	return f.vcsRootID
}

func (f *FeatureCommitStatusPublisher) SetVcsRootID(value string) {
	f.vcsRootID = value
}

func (f *FeatureCommitStatusPublisher) Disabled() bool {
	return f.disabled
}

func (f *FeatureCommitStatusPublisher) SetDisabled(value bool) {
	f.disabled = value
}

func (f *FeatureCommitStatusPublisher) BuildTypeID() string {
	return f.buildTypeID
}

func (f *FeatureCommitStatusPublisher) SetBuildTypeID(value string) {
	f.buildTypeID = value
}

func (f *FeatureCommitStatusPublisher) Properties() *Properties {
	return f.properties
}

func (f *FeatureCommitStatusPublisher) MarshalJSON() ([]byte, error) {
	out := &buildFeatureJSON{
		ID:         f.id,
		Disabled:   NewBool(f.disabled),
		Properties: f.properties,
		Inherited:  NewFalse(),
		Type:       f.Type(),
	}

	return json.Marshal(out)
}

func (f *FeatureCommitStatusPublisher) UnmarshalJSON(data []byte) error {
	var aux buildFeatureJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	f.id = aux.ID

	disabled := aux.Disabled
	if disabled == nil {
		disabled = NewFalse()
	}
	f.disabled = *disabled
	f.properties = NewProperties(aux.Properties.Items...)

	return nil
}

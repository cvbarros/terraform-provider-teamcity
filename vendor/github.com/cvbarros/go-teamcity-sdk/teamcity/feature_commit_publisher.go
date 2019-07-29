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

//ID returns the ID for this instance.
func (f *FeatureCommitStatusPublisher) ID() string {
	return f.id
}

//SetID sets the ID for this instance.
func (f *FeatureCommitStatusPublisher) SetID(value string) {
	f.id = value
}

//Type returns the "commit-status-publisher", the keyed-type for this build feature instance
func (f *FeatureCommitStatusPublisher) Type() string {
	return "commit-status-publisher"
}

//VcsRootID returns the ID that this build feature is associated with.
func (f *FeatureCommitStatusPublisher) VcsRootID() string {
	return f.vcsRootID
}

//SetVcsRootID sets the ID that this build feature is associated with.
func (f *FeatureCommitStatusPublisher) SetVcsRootID(value string) {
	f.vcsRootID = value
}

//Disabled returns whether this build feature is disabled or not.
func (f *FeatureCommitStatusPublisher) Disabled() bool {
	return f.disabled
}

//SetDisabled sets whether this build feature is disabled or not.
func (f *FeatureCommitStatusPublisher) SetDisabled(value bool) {
	f.disabled = value
}

//BuildTypeID is a getter for the Build Type ID associated with this build feature.
func (f *FeatureCommitStatusPublisher) BuildTypeID() string {
	return f.buildTypeID
}

//SetBuildTypeID is a setter for the Build Type ID associated with this build feature.
func (f *FeatureCommitStatusPublisher) SetBuildTypeID(value string) {
	f.buildTypeID = value
}

//Properties returns a *Properties instance representing a serializable collection to be used.
func (f *FeatureCommitStatusPublisher) Properties() *Properties {
	return f.properties
}

//MarshalJSON implements JSON serialization for FeatureCommitStatusPublisher
func (f *FeatureCommitStatusPublisher) MarshalJSON() ([]byte, error) {
	out := &buildFeatureJSON{
		ID:         f.id,
		Disabled:   NewBool(f.disabled),
		Properties: f.properties,
		Inherited:  NewFalse(),
		Type:       f.Type(),
	}

	if f.vcsRootID != "" {
		out.Properties.AddOrReplaceValue("vcsRootId", f.vcsRootID)
	}
	return json.Marshal(out)
}

//UnmarshalJSON implements JSON deserialization for FeatureCommitStatusPublisher
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

	opt, err := CommitStatusPublisherGithubOptionsFromProperties(f.properties)
	if err != nil {
		return err
	}

	if v, ok := f.properties.GetOk("vcsRootId"); ok {
		f.vcsRootID = v
	}
	f.Options = opt

	return nil
}

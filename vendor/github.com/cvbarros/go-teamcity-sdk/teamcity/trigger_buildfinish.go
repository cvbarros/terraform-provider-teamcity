package teamcity

import (
	"encoding/json"
	"fmt"
)

//TriggerBuildFinish represents a build trigger that fires when the given build ID to is finished
type TriggerBuildFinish struct {
	triggerJSON   *triggerJSON
	buildTypeID   string
	SourceBuildID string
	Options       *TriggerBuildFinishOptions
}

//ID for this entity
func (t *TriggerBuildFinish) ID() string {
	return t.triggerJSON.ID
}

//Type returns TriggerTypes.BuildFinish ("buildDependencyTrigger")
func (t *TriggerBuildFinish) Type() string {
	return TriggerTypes.BuildFinish
}

//SetDisabled controls whether this trigger is disabled or not
func (t *TriggerBuildFinish) SetDisabled(disabled bool) {
	t.triggerJSON.Disabled = NewBool(disabled)
}

//Disabled gets the disabled status for this trigger
func (t *TriggerBuildFinish) Disabled() bool {
	return *t.triggerJSON.Disabled
}

//BuildTypeID gets the build type identifier
func (t *TriggerBuildFinish) BuildTypeID() string {
	return t.buildTypeID
}

//SetBuildTypeID sets the build type identifier
func (t *TriggerBuildFinish) SetBuildTypeID(id string) {
	t.buildTypeID = id
}

// NewTriggerBuildFinish returns a finish build trigger type with TriggerVcsOptions.
func NewTriggerBuildFinish(sourceBuildID string, opt *TriggerBuildFinishOptions) (*TriggerBuildFinish, error) {
	if opt == nil {
		return nil, fmt.Errorf("opt parameter must be valid TriggerVcsOptions, not nil")
	}

	newTrigger := &TriggerBuildFinish{
		triggerJSON: &triggerJSON{
			Disabled: NewFalse(),
			Type:     TriggerTypes.Vcs,
		},
		SourceBuildID: sourceBuildID,
		Options:       opt,
	}

	newTrigger.triggerJSON.Properties = newTrigger.properties()
	return newTrigger, nil
}

func (t *TriggerBuildFinish) properties() *Properties {
	props := t.Options.properties()
	props.AddOrReplaceValue("dependsOn", t.SourceBuildID)
	return props
}

//MarshalJSON implements JSON serialization for TriggerVcs
func (t *TriggerBuildFinish) MarshalJSON() ([]byte, error) {
	out := &triggerJSON{
		ID:         t.ID(),
		Type:       t.Type(),
		Disabled:   NewBool(t.Disabled()),
		Properties: t.properties(),
	}

	return json.Marshal(out)
}

//UnmarshalJSON implements JSON deserialization for TriggerBuildFinish
func (t *TriggerBuildFinish) UnmarshalJSON(data []byte) error {
	var aux triggerJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Type != TriggerTypes.BuildFinish {
		return fmt.Errorf("invalid type %s trying to deserialize into TriggerBuildFinish entity", aux.Type)
	}

	if aux.Disabled != nil {
		t.SetDisabled(*aux.Disabled)
	}
	t.triggerJSON = &aux

	opt := aux.Properties.triggerBuildFinishOptions()
	t.Options = opt

	if v, ok := aux.Properties.GetOk("dependsOn"); ok {
		t.SourceBuildID = v
	} else {
		return fmt.Errorf("Invalid 'dependsOn' property. It is mandatory to have a valid BuildID for a finish build trigger")
	}
	return nil
}

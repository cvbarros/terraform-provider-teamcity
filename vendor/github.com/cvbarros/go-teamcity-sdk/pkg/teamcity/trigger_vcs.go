package teamcity

import (
	"encoding/json"
	"fmt"
	"strings"
)

//TriggerVcs represents a build trigger on VCS changes
type TriggerVcs struct {
	triggerJSON  *triggerJSON
	buildTypeID  string
	BranchFilter []string
	Rules        []string
	Options      *TriggerVcsOptions
}

//ID for this entity
func (t *TriggerVcs) ID() string {
	return t.triggerJSON.ID
}

//Type returns TriggerTypes.Vcs ("vcsTrigger")
func (t *TriggerVcs) Type() string {
	return TriggerTypes.Vcs
}

//SetDisabled controls whether this trigger is disabled or not
func (t *TriggerVcs) SetDisabled(disabled bool) {
	t.triggerJSON.Disabled = NewBool(disabled)
}

//Disabled gets the disabled status for this trigger
func (t *TriggerVcs) Disabled() bool {
	return *t.triggerJSON.Disabled
}

//BuildTypeID gets the build type identifier
func (t *TriggerVcs) BuildTypeID() string {
	return t.buildTypeID
}

//SetBuildTypeID sets the build type identifier
func (t *TriggerVcs) SetBuildTypeID(id string) {
	t.buildTypeID = id
}

// NewTriggerVcs returns a VCS trigger type with the triggerRules specified. triggerRules is required, but branchFilter can be optional if the VCS root uses multiple branches.
func NewTriggerVcs(triggerRules []string, branchFilter []string) (*TriggerVcs, error) {
	opt, err := NewTriggerVcsOptions(QuietPeriodDoNotUse, 0)
	if err != nil {
		return nil, err
	}
	i, _ := NewTriggerVcsWithOptions(triggerRules, branchFilter, opt)
	return i, nil
}

// NewTriggerVcsWithOptions returns a VCS trigger type with TriggerVcsOptions. See also NewTriggerVcs for other parameters.
func NewTriggerVcsWithOptions(triggerRules []string, branchFilter []string, opt *TriggerVcsOptions) (*TriggerVcs, error) {
	if opt == nil {
		return nil, fmt.Errorf("opt parameter must be valid TriggerVcsOptions, not nil")
	}

	newTriggerVcs := &TriggerVcs{
		triggerJSON: &triggerJSON{
			Disabled: NewFalse(),
			Type:     TriggerTypes.Vcs,
		},
		BranchFilter: branchFilter,
		Rules:        triggerRules,
		Options:      opt,
	}

	newTriggerVcs.triggerJSON.Properties = newTriggerVcs.properties()
	return newTriggerVcs, nil
}

func (t *TriggerVcs) properties() *Properties {
	props := t.Options.properties()

	if len(t.BranchFilter) > 0 {
		props.AddOrReplaceValue("branchFilter", strings.Join(t.BranchFilter, "\\n"))
	}

	if len(t.Rules) > 0 {
		props.AddOrReplaceValue("triggerRules", strings.Join(t.Rules, "\\n"))
	}

	return props
}

//MarshalJSON implements JSON serialization for TriggerVcs
func (t *TriggerVcs) MarshalJSON() ([]byte, error) {
	out := &triggerJSON{
		ID:         t.ID(),
		Type:       t.Type(),
		Disabled:   NewBool(t.Disabled()),
		Properties: t.properties(),
	}

	return json.Marshal(out)
}

//UnmarshalJSON implements JSON deserialization for TriggerVcs
func (t *TriggerVcs) UnmarshalJSON(data []byte) error {
	var aux triggerJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Type != TriggerTypes.Vcs {
		return fmt.Errorf("invalid type %s trying to deserialize into TriggerVcs entity", aux.Type)
	}

	if aux.Disabled != nil {
		t.SetDisabled(*aux.Disabled)
	}
	t.triggerJSON = &aux

	if v, ok := aux.Properties.GetOk("branchFilter"); ok {
		t.BranchFilter = strings.Split(v, "\\n")
	}

	if v, ok := aux.Properties.GetOk("triggerRules"); ok {
		t.Rules = strings.Split(v, "\\n")
	}

	opt, err := aux.Properties.triggerVcsOptions()
	if err != nil {
		return err
	}
	t.Options = opt

	return nil
}

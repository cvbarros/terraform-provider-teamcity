package teamcity

import (
	"encoding/json"
	"fmt"
)

type triggerType = string

const (
	//BuildTriggerVcs trigger type
	BuildTriggerVcs triggerType = "vcsTrigger"
	//BuildTriggerBuildFinish build trigger type
	BuildTriggerBuildFinish triggerType = "buildDependencyTrigger"
	//BuildTriggerSchedule build trigger tyope
	BuildTriggerSchedule triggerType = "schedulingTrigger"
)

// TriggerTypes represents possible types for build triggers
var TriggerTypes = struct {
	Vcs         triggerType
	BuildFinish triggerType
	Schedule    triggerType
}{
	Vcs:         BuildTriggerVcs,
	BuildFinish: BuildTriggerBuildFinish,
	Schedule:    BuildTriggerSchedule,
}

type triggerJSON struct {
	BuildTypeID string      `json:"-"`
	Disabled    *bool       `json:"disabled,omitempty" xml:"disabled"`
	Href        string      `json:"href,omitempty" xml:"href"`
	ID          string      `json:"id,omitempty" xml:"id"`
	Properties  *Properties `json:"properties,omitempty"`
	Type        string      `json:"type,omitempty" xml:"type"`
}

// Trigger represents a build trigger to be associated with a build configuration. Use the constructor methods to create new instances.
type Trigger interface {
	ID() string
	Type() string
	Disabled() bool
	SetBuildTypeID(buildTypeID string)
	BuildTypeID() string
}

var triggerReadingFunc = func(dt []byte, out interface{}) error {
	var payload triggerJSON
	if err := json.Unmarshal(dt, &payload); err != nil {
		return err
	}

	var obj Trigger
	switch payload.Type {
	case string(TriggerTypes.Vcs):
		var vcs TriggerVcs
		if err := vcs.UnmarshalJSON(dt); err != nil {
			return err
		}
		obj = &vcs
	case string(TriggerTypes.BuildFinish):
		var finish TriggerBuildFinish
		if err := finish.UnmarshalJSON(dt); err != nil {
			return err
		}
		obj = &finish
	case string(TriggerTypes.Schedule):
		var sch TriggerSchedule
		if err := sch.UnmarshalJSON(dt); err != nil {
			return err
		}
		obj = &sch
	default:
		return fmt.Errorf("Unsupported trigger type: '%s' (id:'%s')", payload.Type, payload.ID)
	}

	replaceValue(out, &obj)
	return nil
}

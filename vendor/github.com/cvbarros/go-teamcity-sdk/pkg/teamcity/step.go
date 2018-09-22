package teamcity

import (
	"encoding/json"
	"fmt"
)

// BuildStepType represents most common step types for build steps
type BuildStepType = string

const (
	//StepTypePowershell step type
	StepTypePowershell BuildStepType = "jetbrains_powershell"
	//StepTypeDotnetCli step type
	StepTypeDotnetCli BuildStepType = "dotnet.cli"
	//StepTypeCommandLine (shell/cmd) step type
	StepTypeCommandLine BuildStepType = "simpleRunner"
)

//StepExecuteMode represents how a build configuration step will execute regarding others.
type StepExecuteMode = string

const (
	//StepExecuteModeDefault executes the step only if all previous steps finished successfully.
	StepExecuteModeDefault = "default"
	//StepExecuteModeOnlyIfBuildIsSuccessful executes the step only if the whole build is successful.
	StepExecuteModeOnlyIfBuildIsSuccessful = "execute_if_success"
	//StepExecuteModeEvenWhenFailed executes the step even if previous steps failed.
	StepExecuteModeEvenWhenFailed = "execute_if_failed"
	//StepExecuteAlways executes even if build stop command was issued.
	StepExecuteAlways = "execute_always"
)

// Step interface represents a a build configuration/template build step. To intereact with concrete step types, see the Step* types.
type Step interface {
	ID() string
	Type() string
	Name() string
}

type stepJSON struct {
	Disabled   *bool       `json:"disabled,omitempty" xml:"disabled"`
	Href       string      `json:"href,omitempty" xml:"href"`
	ID         string      `json:"id,omitempty" xml:"id"`
	Inherited  *bool       `json:"inherited,omitempty" xml:"inherited"`
	Name       string      `json:"name,omitempty" xml:"name"`
	Properties *Properties `json:"properties,omitempty"`
	Type       string      `json:"type,omitempty" xml:"type"`
}

type stepsJSON struct {
	Count int32       `json:"count,omitempty" xml:"count"`
	Items []*stepJSON `json:"step"`
}

var stepReadingFunc = func(dt []byte, out interface{}) error {
	var payload stepJSON
	if err := json.Unmarshal(dt, &payload); err != nil {
		return err
	}

	var step Step
	switch payload.Type {
	case string(StepTypePowershell):
		var ps StepPowershell
		if err := ps.UnmarshalJSON(dt); err != nil {
			return err
		}
		step = &ps
	case string(StepTypeCommandLine):
		var cmd StepCommandLine
		if err := cmd.UnmarshalJSON(dt); err != nil {
			return err
		}
		step = &cmd
	default:
		return fmt.Errorf("Unsupported step type: '%s' (id:'%s')", payload.Type, payload.ID)
	}

	replaceValue(out, &step)
	return nil
}

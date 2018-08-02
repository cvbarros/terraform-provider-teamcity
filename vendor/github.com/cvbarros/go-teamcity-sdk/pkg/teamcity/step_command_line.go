package teamcity

import (
	"encoding/json"
	"errors"
	"fmt"
)

//StepCommandLine represents a a build step of type "CommandLine"
type StepCommandLine struct {
	id           string
	name         string
	stepType     string
	stepJSON     *stepJSON
	isExecutable bool
	//CustomScript contains code for platform specific script, like .cmd on windows or shell script on Unix-like environments.
	CustomScript string
	//CommandExecutable is the executable program to be called from this step.
	CommandExecutable string
	//CommandParameters are additional parameters to be passed on to the CommandExecutable.
	CommandParameters string
	//ExecuteMode is the execute mode for the step. See StepExecuteMode for details.
	ExecuteMode StepExecuteMode
}

//NewStepCommandLineScript creates a command line build step that runs an inline platform-specific script.
func NewStepCommandLineScript(name string, script string) (*StepCommandLine, error) {
	if script == "" {
		return nil, errors.New("script is required")
	}

	return &StepCommandLine{
		name:         name,
		isExecutable: false,
		stepType:     StepTypeCommandLine,
		CustomScript: script,
		ExecuteMode:  StepExecuteModeDefault,
	}, nil
}

//NewStepCommandLineExecutable creates a command line that invokes an external executable.
func NewStepCommandLineExecutable(name string, executable string, args string) (*StepCommandLine, error) {
	if executable == "" {
		return nil, errors.New("executable is required")
	}

	return &StepCommandLine{
		name:              name,
		stepType:          StepTypeCommandLine,
		isExecutable:      true,
		CommandExecutable: executable,
		CommandParameters: args,
		ExecuteMode:       StepExecuteModeDefault,
	}, nil
}

//ID for this entity.
func (s *StepCommandLine) ID() string {
	return s.id
}

//Name is a useful description of the step.
func (s *StepCommandLine) Name() string {
	return s.name
}

//Type returns the step type, in this case "StepTypePowershell".
func (s *StepCommandLine) Type() BuildStepType {
	return StepTypeCommandLine
}

func (s *StepCommandLine) properties() *Properties {
	props := NewPropertiesEmpty()
	props.AddOrReplaceValue("teamcity.step.mode", string(s.ExecuteMode))

	if s.isExecutable {
		props.AddOrReplaceValue("command.executable", s.CommandExecutable)

		if s.CommandParameters != "" {
			props.AddOrReplaceValue("command.parameters", s.CommandParameters)
		}
	} else {
		props.AddOrReplaceValue("script.content", s.CustomScript)
		props.AddOrReplaceValue("use.custom.script", "true")
	}

	return props
}

//MarshalJSON implements JSON serialization for StepCommandLine
func (s *StepCommandLine) MarshalJSON() ([]byte, error) {
	out := &stepJSON{
		ID:         s.id,
		Name:       s.name,
		Type:       s.stepType,
		Properties: s.properties(),
	}

	return json.Marshal(out)
}

//UnmarshalJSON implements JSON deserialization for StepCommandLine
func (s *StepCommandLine) UnmarshalJSON(data []byte) error {
	var aux stepJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Type != string(StepTypeCommandLine) {
		return fmt.Errorf("invalid type %s trying to deserialize into StepCommandLine entity", aux.Type)
	}
	s.name = aux.Name
	s.id = aux.ID
	s.stepType = StepTypeCommandLine

	props := aux.Properties
	if _, ok := props.GetOk("use.custom.script"); ok {
		s.isExecutable = false
		if v, ok := props.GetOk("script.content"); ok {
			s.CustomScript = v
		}
	}

	if v, ok := props.GetOk("command.executable"); ok {
		s.CommandExecutable = v
		if v, ok := props.GetOk("command.parameters"); ok {
			s.CommandParameters = v
		}
	}

	if v, ok := props.GetOk("teamcity.step.mode"); ok {
		s.ExecuteMode = StepExecuteMode(v)
	}
	return nil
}

package teamcity

import (
	"encoding/json"
	"errors"
	"fmt"
)

//StepPowershell represents a a build step of type "Powershell"
type StepPowershell struct {
	ID       string
	Name     string
	stepType string
	stepJSON *stepJSON
	isScript bool
	//ScriptFile holds the name of script to run for this step.
	ScriptFile string
	//Code is the inline powershell code to be ran for this step.
	Code string
	//ScriptArgs are the arguments that will be passed when using "ScriptFile"
	ScriptArgs string
	//ExecuteMode is the execute mode for the step. See StepExecuteMode for details.
	ExecuteMode StepExecuteMode
}

//NewStepPowershellScriptFile creates a powershell build step that runs a script file instead of inline code.
func NewStepPowershellScriptFile(name string, scriptFile string, scriptArgs string) (*StepPowershell, error) {
	if scriptFile == "" {
		return nil, errors.New("scriptFile is required")
	}

	return &StepPowershell{
		Name:        name,
		isScript:    true,
		stepType:    StepTypePowershell,
		ScriptFile:  scriptFile,
		ScriptArgs:  scriptArgs,
		ExecuteMode: StepExecuteModeDefault,
	}, nil
}

//NewStepPowershellCode creates a powershell build step that runs the inline code.
func NewStepPowershellCode(name string, code string) (*StepPowershell, error) {
	if code == "" {
		return nil, errors.New("code is required")
	}

	return &StepPowershell{
		Name:        name,
		stepType:    StepTypePowershell,
		Code:        code,
		ExecuteMode: StepExecuteModeDefault,
	}, nil
}

//GetID is a wrapper implementation for ID field, to comply with Step interface
func (s *StepPowershell) GetID() string {
	return s.ID
}

//GetName is a wrapper implementation for Name field, to comply with Step interface
func (s *StepPowershell) GetName() string {
	return s.Name
}

//Type returns the step type, in this case "StepTypePowershell".
func (s *StepPowershell) Type() BuildStepType {
	return StepTypePowershell
}

func (s *StepPowershell) properties() *Properties {
	props := NewPropertiesEmpty()
	props.AddOrReplaceValue("teamcity.step.mode", string(s.ExecuteMode))
	// Defaults
	props.AddOrReplaceValue("jetbrains_powershell_noprofile", "true")
	props.AddOrReplaceValue("jetbrains_powershell_execution", "PS1")

	if s.isScript {
		props.AddOrReplaceValue("jetbrains_powershell_script_mode", "FILE")
		props.AddOrReplaceValue("jetbrains_powershell_script_file", s.ScriptFile)

		if s.ScriptArgs != "" {
			props.AddOrReplaceValue("jetbrains_powershell_scriptArguments", s.ScriptArgs)
		}
	} else {
		props.AddOrReplaceValue("jetbrains_powershell_script_mode", "CODE")
		props.AddOrReplaceValue("jetbrains_powershell_script_code", s.Code)
	}

	return props
}

//MarshalJSON implements JSON serialization for StepPowershell
func (s *StepPowershell) MarshalJSON() ([]byte, error) {
	out := &stepJSON{
		ID:         s.ID,
		Name:       s.Name,
		Type:       s.stepType,
		Properties: s.properties(),
	}

	return json.Marshal(out)
}

//UnmarshalJSON implements JSON deserialization for StepPowershell
func (s *StepPowershell) UnmarshalJSON(data []byte) error {
	var aux stepJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Type != string(StepTypePowershell) {
		return fmt.Errorf("invalid type %s trying to deserialize into StepPowershell entity", aux.Type)
	}
	s.Name = aux.Name
	s.ID = aux.ID
	s.stepType = StepTypePowershell

	props := aux.Properties
	if v, ok := props.GetOk("jetbrains_powershell_script_file"); ok {
		s.ScriptFile = v
		if v, ok := props.GetOk("jetbrains_powershell_scriptArguments"); ok {
			s.ScriptArgs = v
		}
		s.isScript = true
	}

	if v, ok := props.GetOk("jetbrains_powershell_script_code"); ok {
		s.Code = v
		s.isScript = false
	}

	if v, ok := props.GetOk("teamcity.step.mode"); ok {
		s.ExecuteMode = StepExecuteMode(v)
	}
	return nil
}

package teamcity

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// StepOctopusPushPackage represents a a build step of type "octopus.push.package"
type StepOctopusPushPackage struct {
	ID       string
	Name     string
	stepType string
	stepJSON *stepJSON

	// Specify Octopus web portal URL.
	Host string

	// Specify Octopus API key.
	ApiKey string

	// Specify  Package path patterns.
	PackagePaths string

	// Force overwrite existing packages.
	ForcePush bool

	// Automatically publish any packages as TeamCity build artifacts.
	PublishArtifacts bool

	// Additional arguments to be passed to Octo.exe.
	AdditionalCommandLineArguments string
}

func NewStepOctopusPushPackage(name string) (*StepOctopusPushPackage, error) {
	return &StepOctopusPushPackage{
		Name:     name,
		stepType: StepTypeOctopusPushPackage,
	}, nil
}

func (s *StepOctopusPushPackage) GetID() string {
	return s.ID
}

func (s *StepOctopusPushPackage) GetName() string {
	return s.Name
}

func (s *StepOctopusPushPackage) Type() BuildStepType {
	return StepTypeOctopusPushPackage
}

func (s *StepOctopusPushPackage) properties() *Properties {
	props := NewPropertiesEmpty()
	props.AddOrReplaceValue("teamcity.step.mode", "default")
	props.AddOrReplaceValue("octopus_host", s.Host)
	props.AddOrReplaceValue("secure:octopus_apikey", s.ApiKey)
	props.AddOrReplaceValue("octopus_packagepaths", s.PackagePaths)
	props.AddOrReplaceValue("octopus_forcepush", strconv.FormatBool(s.ForcePush))
	props.AddOrReplaceValue("octopus_publishartifacts", strconv.FormatBool(s.PublishArtifacts))
	props.AddOrReplaceValue("octopus_additionalcommandlinearguments", s.AdditionalCommandLineArguments)

	return props
}

func (s *StepOctopusPushPackage) serializable() *stepJSON {
	return &stepJSON{
		ID:         s.ID,
		Name:       s.Name,
		Type:       s.stepType,
		Properties: s.properties(),
	}
}

//MarshalJSON implements JSON serialization for StepOctopusPushPackage
func (s *StepOctopusPushPackage) MarshalJSON() ([]byte, error) {
	out := s.serializable()
	return json.Marshal(out)
}

// UnmarshalJSON implements JSON deserialization for StepOctopusPushPackage
func (s *StepOctopusPushPackage) UnmarshalJSON(data []byte) error {
	var aux stepJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Type != string(StepTypeOctopusPushPackage) {
		return fmt.Errorf("invalid type %s trying to deserialize into StepOctopusPushPackage entity", aux.Type)
	}
	s.Name = aux.Name
	s.ID = aux.ID
	s.stepType = StepTypeOctopusPushPackage

	props := aux.Properties
	if v, ok := props.GetOk("octopus_host"); ok {
		s.Host = v
	}

	if v, ok := props.GetOk("secure:octopus_apikey"); ok {
		s.ApiKey = v
	}

	if v, ok := props.GetOk("octopus_packagepaths"); ok {
		s.PackagePaths = v
	}

	if v, ok := props.GetOk("octopus_forcepush"); ok {
		converted_value, _ := strconv.ParseBool(v)
		s.ForcePush = converted_value
	}

	if v, ok := props.GetOk("octopus_publishartifacts"); ok {
		converted_value, _ := strconv.ParseBool(v)
		s.PublishArtifacts = converted_value
	}

	if v, ok := props.GetOk("octopus_additionalcommandlinearguments"); ok {
		s.AdditionalCommandLineArguments = v
	}

	return nil

}

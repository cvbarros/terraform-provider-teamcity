package teamcity

import "github.com/lann/builder"

type stepPowershellBuilder builder.Builder

// ScriptFile sets properties required to run the powershell step as a script file
func (b stepPowershellBuilder) ScriptFile(scriptFile string) stepPowershellBuilder {
	out := addOrReplaceProperty(b, "jetbrains_powershell_script_mode", "FILE")
	return addOrReplaceProperty(out, "jetbrains_powershell_script_file", scriptFile).(stepPowershellBuilder)
}

func (b stepPowershellBuilder) Code(scriptCode string) stepPowershellBuilder {
	out := addOrReplaceProperty(b, "jetbrains_powershell_script_mode", "CODE")
	return addOrReplaceProperty(out, "jetbrains_powershell_script_code", scriptCode).(stepPowershellBuilder)
}

// Args sets properties required for script arguments
func (b stepPowershellBuilder) Args(args string) stepPowershellBuilder {
	return addOrReplaceProperty(b, "jetbrains_powershell_scriptArguments", args).(stepPowershellBuilder)
}

func (b stepPowershellBuilder) Build(name string) *Step {
	// Defaults
	b2 := addOrReplaceProperty(b, "teamcity.step.mode", "default")
	b2 = addOrReplaceProperty(b2, "jetbrains_powershell_noprofile", "true")
	b2 = addOrReplaceProperty(b2, "jetbrains_powershell_execution", "PS1")

	out := builder.GetStruct(b2).(Step)
	out.Type = StepTypes.Powershell
	out.Name = name
	return &out
}

func addOrReplaceProperty(b interface{}, name string, value string) interface{} {
	newProp := &Property{
		Name:  name,
		Value: value,
	}
	ret, exists := builder.Get(b, "Properties")
	if !exists {
		props := NewProperties(newProp)
		return builder.Set(b, "Properties", props)
	}

	props := ret.(*Properties)
	props.AddOrReplaceValue(newProp)
	return builder.Set(b, "Properties", props)
}

// StepPowershellBuilder is a convenience class for creating powershell build steps
var StepPowershellBuilder = builder.Register(stepPowershellBuilder{}, Step{}).(stepPowershellBuilder)

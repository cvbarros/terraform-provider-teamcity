package teamcity

type stepType = string

const (
	//Powershell step type
	Powershell stepType = "jetbrains_powershell"
	//DotnetCli step type
	DotnetCli stepType = "dotnet.cli"
	//CommandLine (shell/cmd) step type
	CommandLine stepType = "simpleRunner"
)

// StepTypes represents most common step types for build steps
var StepTypes = struct {
	Powershell  stepType
	DotnetCli   stepType
	CommandLine stepType
}{
	Powershell:  Powershell,
	DotnetCli:   DotnetCli,
	CommandLine: CommandLine,
}

package teamcity

type triggerType = string

const (
	//Vcs trigger type
	Vcs triggerType = "vcsTrigger"
	//Dependency is "Wait for build to finish trigger type"
	Dependency triggerType = "buildDependencyTrigger"
	//Schedule trigger type
	Schedule triggerType = "schedulingTrigger"
)

// TriggerTypes represents possible types for build triggers
var TriggerTypes = struct {
	Vcs        triggerType
	Dependency triggerType
	Schedule   triggerType
}{
	Vcs:        Vcs,
	Dependency: Dependency,
	Schedule:   Schedule,
}

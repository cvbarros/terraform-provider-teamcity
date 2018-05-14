package teamcity

type triggerType = string

const (
	//VCS Trigger type
	Vcs triggerType = "vcsTrigger"
	//Wait for build to finish trigger type
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

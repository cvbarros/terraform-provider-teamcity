package teamcity

//DefaultBuildNumberFormat is TC's default build number format setting for build configurations
const DefaultBuildNumberFormat = "%build.counter%"

//DefaultBuildConfigurationType is default build configuration type setting for build configurations.
//Other possible values for this setting would be "DEPLOYMENT" or "COMPOSITE"
const DefaultBuildConfigurationType = "REGULAR"

//BuildTypeOptions represents settings for a Build Configuration
type BuildTypeOptions struct {
	AllowPersonalBuildTriggering bool     `prop:"allowPersonalBuildTriggering" force:""`
	ArtifactRules                []string `prop:"artifactRules" separator:"\n"`
	EnableHangingBuildsDetection bool     `prop:"enableHangingBuildsDetection" force:""`
	EnableStatusWidget           bool     `prop:"allowExternalStatus"`
	BuildCounter                 int      `prop:"buildNumberCounter"`
	BuildNumberFormat            string   `prop:"buildNumberPattern"`
	BuildConfigurationType       string   `prop:"buildConfigurationType"`
	MaxSimultaneousBuilds        int      `prop:"maximumNumberOfBuilds"`
	Template                     bool
	BuildTypeID                  int
}

//NewBuildTypeOptionsWithDefaults returns a new instance of default settings, the same as presented in the TeamCity UI when a new build configuration is created.
func NewBuildTypeOptionsWithDefaults() *BuildTypeOptions {
	return &BuildTypeOptions{
		AllowPersonalBuildTriggering: true,
		ArtifactRules:                []string{},
		EnableHangingBuildsDetection: true,
		EnableStatusWidget:           false,
		MaxSimultaneousBuilds:        0,
		BuildConfigurationType:       DefaultBuildConfigurationType,
		BuildCounter:                 1,
		BuildNumberFormat:            DefaultBuildNumberFormat,
	}
}

//NewBuildTypeOptionsTemplate returns a new instance of settings for a BuildType Template.
func NewBuildTypeOptionsTemplate() *BuildTypeOptions {
	return &BuildTypeOptions{
		AllowPersonalBuildTriggering: true,
		ArtifactRules:                []string{},
		EnableHangingBuildsDetection: true,
		EnableStatusWidget:           false,
		MaxSimultaneousBuilds:        0,
		BuildConfigurationType:       DefaultBuildConfigurationType,
		BuildNumberFormat:            DefaultBuildNumberFormat,
		Template:                     true,
	}
}

func (o *BuildTypeOptions) properties() *Properties {
	props := serializeToProperties(o)

	//TeamCity API for build settings has a very weird behavior to omit some properties when they assume their "default" value.
	//In this case, in order to keep consistent behaviour between reads/writes, the property raw model is adjusted for this behaviour.

	//Omit allowPersonalBuildTriggering if equals to default 'true'
	if o.AllowPersonalBuildTriggering {
		props.Remove("allowPersonalBuildTriggering")
	}

	//Omit enableHangingBuildsDetection if equals to default 'true'
	if o.EnableHangingBuildsDetection {
		props.Remove("enableHangingBuildsDetection")
	}

	//Omit if buildConfigurationType == "REGULAR"
	if v, ok := props.GetOk("buildConfigurationType"); ok && v == DefaultBuildConfigurationType {
		props.Remove("buildConfigurationType")
	}

	//Omit if buildNumberPattern == "%build.counter%"
	if v, ok := props.GetOk("buildNumberPattern"); ok && v == DefaultBuildNumberFormat {
		props.Remove("buildNumberPattern")
	}

	if v, ok := props.GetOk("maximumNumberOfBuilds"); ok && v == "0" {
		props.Remove("maximumNumberOfBuilds")
	}

	//Build Type Templates do not have "buildNumberCounter" property, so remove that if these options are from a template
	if o.Template {
		props.Remove("buildNumberCounter")
	}
	return props
}

func (p *Properties) buildTypeOptions(template bool) (out *BuildTypeOptions) {
	if template {
		out = NewBuildTypeOptionsTemplate()
	} else {
		out = NewBuildTypeOptionsWithDefaults()
	}

	fillStructFromProperties(out, p)
	return out
}

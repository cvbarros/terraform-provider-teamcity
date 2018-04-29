package teamcity

import (
	"strconv"
	"strings"

	"github.com/lann/builder"
)

type buildTypeSettingsBuilder builder.Builder

// BuildTypeSettingsBuilder is a convenience class for creating build type settings
var BuildTypeSettingsBuilder = builder.Register(buildTypeSettingsBuilder{}, Properties{}).(buildTypeSettingsBuilder)

// ConfigurationType changes this builder to use the specified configuration type, which can be REGULAR, DEPLOYMENT or COMPOSITE
// Anything different than those strings (case insentive) will be ignored. No error returned due the fluent interface of builder
func (b buildTypeSettingsBuilder) ConfigurationType(t string) buildTypeSettingsBuilder {
	lcaseType := strings.ToLower(t)
	if lcaseType != "regular" && lcaseType != "deployment" && lcaseType != "composite" {
		return b
	}

	return builder.Set(b, "buildConfigurationType", strings.ToUpper(lcaseType)).(buildTypeSettingsBuilder)
}

func (b buildTypeSettingsBuilder) PersonalBuildTrigger(value bool) buildTypeSettingsBuilder {
	return builder.Set(b, "allowPersonalBuildTriggering", strconv.FormatBool(value)).(buildTypeSettingsBuilder)
}

func (b buildTypeSettingsBuilder) HangingBuildDetection(value bool) buildTypeSettingsBuilder {
	return builder.Set(b, "enableHangingBuildsDetection", strconv.FormatBool(value)).(buildTypeSettingsBuilder)
}

func (b buildTypeSettingsBuilder) ArtifactRules(value string) buildTypeSettingsBuilder {
	return builder.Set(b, "artifactRules", value).(buildTypeSettingsBuilder)
}

func (b buildTypeSettingsBuilder) MaxConcurrentBuilds(n int) buildTypeSettingsBuilder {
	return builder.Set(b, "maximumNumberOfBuilds", strconv.Itoa(n)).(buildTypeSettingsBuilder)
}

func (b buildTypeSettingsBuilder) BuildCounter(n int) buildTypeSettingsBuilder {
	return builder.Set(b, "buildNumberCounter", strconv.Itoa(n)).(buildTypeSettingsBuilder)
}

func (b buildTypeSettingsBuilder) BuildNumberPattern(pattern string) buildTypeSettingsBuilder {
	return builder.Set(b, "buildNumberPattern", pattern).(buildTypeSettingsBuilder)
}

func (b buildTypeSettingsBuilder) Build() *Properties {
	var props []*Property

	props = appendPropertyIfApplicable(b, props, "buildConfigurationType")
	props = appendPropertyIfApplicable(b, props, "allowPersonalBuildTriggering")
	props = appendPropertyIfApplicable(b, props, "enableHangingBuildsDetection")
	props = appendPropertyIfApplicable(b, props, "artifactRules")
	props = appendPropertyIfApplicable(b, props, "maximumNumberOfBuilds")
	props = appendPropertyIfApplicable(b, props, "buildCounter")
	props = appendPropertyIfApplicable(b, props, "buildNumberPattern")

	return NewProperties(props...)
}

func appendPropertyIfApplicable(b buildTypeSettingsBuilder, props []*Property, name string) []*Property {
	if value, ok := builder.Get(b, name); ok {
		prop := &Property{
			Name:  name,
			Value: value.(string),
		}
		props = append(props, prop)
	}
	return props
}

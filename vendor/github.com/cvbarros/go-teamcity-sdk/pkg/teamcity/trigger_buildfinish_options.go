package teamcity

import (
	"strconv"
	"strings"
)

// TriggerBuildFinishOptions represents optional settings for a 'Finish Build' Trigger type.
type TriggerBuildFinishOptions struct {
	AfterSuccessfulBuildOnly bool     `prop:"afterSucessfulBuildOnly"`
	BranchFilter             []string `prop:"branchFilter" separator:"\n"`
}

// NewTriggerBuildFinishOptions initialize a NewTriggerBuildFinishOptions
// branchFilter can be passed as "nil" to not filter on any specific branches
func NewTriggerBuildFinishOptions(afterSuccessfulBuildOnly bool, branchFilter []string) *TriggerBuildFinishOptions {
	return &TriggerBuildFinishOptions{
		AfterSuccessfulBuildOnly: afterSuccessfulBuildOnly,
		BranchFilter:             branchFilter,
	}
}

func (o *TriggerBuildFinishOptions) properties() *Properties {
	props := NewPropertiesEmpty()

	//Defaults to false, so ommit emitting the property if 'false'
	if o.AfterSuccessfulBuildOnly {
		props.AddOrReplaceValue("afterSuccessfulBuildOnly", strconv.FormatBool(o.AfterSuccessfulBuildOnly))
	}

	if o.BranchFilter != nil && len(o.BranchFilter) > 0 {
		props.AddOrReplaceValue("branchFilter", strings.Join(o.BranchFilter, "\n"))
	}

	return props
}

func (p *Properties) triggerBuildFinishOptions() *TriggerBuildFinishOptions {
	var out TriggerBuildFinishOptions

	fillStructFromProperties(&out, p)

	return &out
}

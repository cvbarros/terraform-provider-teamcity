package teamcity

import (
	"fmt"
	"strconv"
)

// VcsTriggerQuietPeriodMode specifies if the VCS Trigger will delay the start of a build after detecting a VCS change, used by VcsTriggerOptions.
type VcsTriggerQuietPeriodMode int

const (
	// QuietPeriodDoNotUse disables QuietPeriod on VCS Trigger
	QuietPeriodDoNotUse VcsTriggerQuietPeriodMode = 0
	// QuietPeriodUseDefault instructs the VCS Trigger to respect the server-wide quiet period
	QuietPeriodUseDefault VcsTriggerQuietPeriodMode = 1
	// QuietPeriodCustom allows specifying a period in seconds via VcsTriggerOptions.QuietPeriodInSeconds. When using Custom, VcsTriggerOptions.QuietPeriodInSeconds is mandatory.
	QuietPeriodCustom VcsTriggerQuietPeriodMode = 2
)

// VcsTriggerOptions represents optional settings for a VCS Trigger type.
type VcsTriggerOptions struct {
	perCheckinTriggering bool

	groupUserCheckins       bool
	enableQueueOptimization bool
	quietPeriodMode         VcsTriggerQuietPeriodMode
	quietPeriodInSeconds    int
}

// NewVcsTriggerOptions initialize a VcsTriggerOptions instance with same defaults as TeamCity UI
//
// Defaults:
//	//GroupCheckins = false
//	//EnableQueueOptimization = false
//	//QuietPeriodMode =  QuietPeriodDoNotUse (0)
func NewVcsTriggerOptions() *VcsTriggerOptions {
	return &VcsTriggerOptions{
		perCheckinTriggering:    false,
		groupUserCheckins:       false,
		enableQueueOptimization: true,
		quietPeriodMode:         QuietPeriodDoNotUse,
	}
}

//SetQueueOptimization toggles allowing the server to replace an already started build or a more recently queued one if new changes are detected. If set to true, PerCheckinTriggering will be disabled.
func (o *VcsTriggerOptions) SetQueueOptimization(enable bool) {
	o.enableQueueOptimization = enable
	if enable {
		o.perCheckinTriggering = false
	}
}

// SetPerCheckinTriggering specifies if VCS Trigger will fire a different build per checkin or commit for different committers. If set to true, EnableQueueOptimization will be disabled.
func (o *VcsTriggerOptions) SetPerCheckinTriggering(enable bool) {
	o.perCheckinTriggering = enable
	if enable {
		o.enableQueueOptimization = false
	}
}

// SetQuietPeriodMode controls quiet period mode for the VCSTrigger. See VcsTriggerQuietPeriodMode. If 'QuietPeriodCustom' (2), then QuietPeriodInSeconds is required to be greater than zero.
// seconds parameter is the period in seconds the build will be delayed after a change is detected. Only used if QuietPeriodMode = QuietPeriodCustom, otherwise this is always set to zero.
func (o *VcsTriggerOptions) SetQuietPeriodMode(mode VcsTriggerQuietPeriodMode, seconds int) error {
	quietPeriodInSeconds := 0
	if mode == QuietPeriodCustom {
		if seconds <= 0 {
			return fmt.Errorf("invalid valid %d for QuietPeriodInSeconds. Must be greater than zero when QuietPeriodModeCustom is used", o.quietPeriodInSeconds)
		}
		quietPeriodInSeconds = seconds
	}

	o.quietPeriodMode = mode
	o.quietPeriodInSeconds = quietPeriodInSeconds
	return nil
}

// SetGroupUserCheckins specifies if the server should trigger group user checkins before triggering a build.
func (o *VcsTriggerOptions) SetGroupUserCheckins(enable bool) {
	o.groupUserCheckins = enable
}

var quietPeriodModePropertyMap = map[VcsTriggerQuietPeriodMode]string{
	QuietPeriodDoNotUse:   "DO_NOT_USE",
	QuietPeriodUseDefault: "USE_DEFAULT",
	QuietPeriodCustom:     "USE_CUSTOM",
}

func (o *VcsTriggerOptions) vcsTriggerProperties() *Properties {
	var props []*Property

	p := NewProperty("quietPeriodMode", quietPeriodModePropertyMap[o.quietPeriodMode])
	props = append(props, p)

	if o.enableQueueOptimization {
		p := NewProperty("enableQueueOptimization", "true")
		props = append(props, p)
	}

	if o.perCheckinTriggering {
		p := NewProperty("perCheckinTriggering", "true")
		props = append(props, p)
	}

	if o.groupUserCheckins {
		p := NewProperty("groupCheckinsByCommitter", "true")
		props = append(props, p)
	}

	if o.quietPeriodInSeconds > 0 {
		p := NewProperty("quietPeriod", strconv.Itoa(o.quietPeriodInSeconds))
		props = append(props, p)
	}

	return NewProperties(props...)
}

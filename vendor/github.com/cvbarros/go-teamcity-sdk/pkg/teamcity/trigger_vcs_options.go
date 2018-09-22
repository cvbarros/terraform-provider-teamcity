package teamcity

import (
	"fmt"
	"strconv"
)

// VcsTriggerQuietPeriodMode specifies if the VCS Trigger will delay the start of a build after detecting a VCS change, used by TriggerVcsOptions.
type VcsTriggerQuietPeriodMode int

const (
	// QuietPeriodDoNotUse disables QuietPeriod on VCS Trigger
	QuietPeriodDoNotUse VcsTriggerQuietPeriodMode = 0
	// QuietPeriodUseDefault instructs the VCS Trigger to respect the server-wide quiet period
	QuietPeriodUseDefault VcsTriggerQuietPeriodMode = 1
	// QuietPeriodCustom allows specifying a period in seconds via TriggerVcsOptions.QuietPeriodInSeconds. When using Custom, TriggerVcsOptions.QuietPeriodInSeconds is mandatory.
	QuietPeriodCustom VcsTriggerQuietPeriodMode = 2
)

// TriggerVcsOptions represents optional settings for a VCS Trigger type.
type TriggerVcsOptions struct {
	enableQueueOptimization bool
	perCheckinTriggering    bool `prop:"branch"`

	GroupUserCheckins    bool
	QuietPeriodMode      VcsTriggerQuietPeriodMode
	QuietPeriodInSeconds int
}

// NewTriggerVcsOptions initialize a TriggerVcsOptions instance with same defaults as TeamCity UI
//
// Defaults:
//	- GroupCheckins = false
//	- EnableQueueOptimization = false
func NewTriggerVcsOptions(mode VcsTriggerQuietPeriodMode, seconds int) (*TriggerVcsOptions, error) {
	quietPeriodInSeconds := 0
	if mode == QuietPeriodCustom {
		if seconds <= 0 {
			return nil, fmt.Errorf("invalid valid %d for QuietPeriodInSeconds. Must be greater than zero when QuietPeriodModeCustom is used", quietPeriodInSeconds)
		}
		quietPeriodInSeconds = seconds
	}

	return &TriggerVcsOptions{
		perCheckinTriggering:    false,
		enableQueueOptimization: true,
		GroupUserCheckins:       false,
		QuietPeriodMode:         mode,
		QuietPeriodInSeconds:    quietPeriodInSeconds,
	}, nil
}

//QueueOptimization gets the value of enableQueueOptimization property
func (o *TriggerVcsOptions) QueueOptimization() bool {
	return o.enableQueueOptimization
}

//SetQueueOptimization toggles allowing the server to replace an already started build or a more recently queued one if new changes are detected. If set to true, PerCheckinTriggering will be disabled.
func (o *TriggerVcsOptions) SetQueueOptimization(enable bool) {
	o.enableQueueOptimization = enable
	if enable {
		o.perCheckinTriggering = false
	}
}

//PerCheckinTriggering gets the value of perCheckinTriggering property
func (o *TriggerVcsOptions) PerCheckinTriggering() bool {
	return o.perCheckinTriggering
}

// SetPerCheckinTriggering specifies if VCS Trigger will fire a different build per checkin or commit for different committers. If set to true, enableQueueOptimization will be disabled.
func (o *TriggerVcsOptions) SetPerCheckinTriggering(enable bool) {
	o.perCheckinTriggering = enable
	if enable {
		o.enableQueueOptimization = false
	}
}

var quietPeriodModePropertyMap = map[VcsTriggerQuietPeriodMode]string{
	QuietPeriodDoNotUse:   "DO_NOT_USE",
	QuietPeriodUseDefault: "USE_DEFAULT",
	QuietPeriodCustom:     "USE_CUSTOM",
}

var propertyToQuietPeriodModeMap = map[string]VcsTriggerQuietPeriodMode{
	"DO_NOT_USE":  QuietPeriodDoNotUse,
	"USE_DEFAULT": QuietPeriodUseDefault,
	"USE_CUSTOM":  QuietPeriodCustom,
}

func (o *TriggerVcsOptions) properties() *Properties {
	var props []*Property

	p := NewProperty("quietPeriodMode", quietPeriodModePropertyMap[o.QuietPeriodMode])
	props = append(props, p)

	if o.enableQueueOptimization {
		p := NewProperty("enableQueueOptimization", "true")
		props = append(props, p)
	}

	if o.perCheckinTriggering {
		p := NewProperty("perCheckinTriggering", "true")
		props = append(props, p)
	}

	if o.GroupUserCheckins {
		p := NewProperty("groupCheckinsByCommitter", "true")
		props = append(props, p)
	}

	if o.QuietPeriodInSeconds > 0 {
		p := NewProperty("quietPeriod", strconv.Itoa(o.QuietPeriodInSeconds))
		props = append(props, p)
	}

	return NewProperties(props...)
}

func (p *Properties) triggerVcsOptions() (*TriggerVcsOptions, error) {
	var out TriggerVcsOptions
	if v, ok := p.GetOk("quietPeriodMode"); ok {
		out.QuietPeriodMode = propertyToQuietPeriodModeMap[v]
	}
	if v, ok := p.GetOk("quietPeriod"); ok {
		v2, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		out.QuietPeriodInSeconds = v2
	}
	if v, ok := p.GetOk("groupCheckinsByCommitter"); ok {
		v2, err := strconv.ParseBool(v)
		if err != nil {
			return nil, err
		}
		out.GroupUserCheckins = v2
	}
	if v, ok := p.GetOk("perCheckinTriggering"); ok {
		v2, err := strconv.ParseBool(v)
		if err != nil {
			return nil, err
		}
		if v2 {
			out.SetPerCheckinTriggering(v2)
		}
	}
	if v, ok := p.GetOk("enableQueueOptimization"); ok {
		v2, err := strconv.ParseBool(v)
		if err != nil {
			return nil, err
		}
		if v2 {
			out.SetQueueOptimization(v2)
		}
	}

	return &out, nil
}

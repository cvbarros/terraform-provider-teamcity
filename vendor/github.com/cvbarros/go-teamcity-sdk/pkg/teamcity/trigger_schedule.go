package teamcity

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// TriggerSchedulingPolicy represents the shceduling policy for a scheduled trigger. Can be 'daily', 'weekly' or 'cron'
type TriggerSchedulingPolicy = string

const (
	//TriggerSchedulingDaily triggers every day
	TriggerSchedulingDaily TriggerSchedulingPolicy = "daily"
	//TriggerSchedulingWeekly triggers at specified day + time of the week, once per week
	TriggerSchedulingWeekly TriggerSchedulingPolicy = "weekly"
	//TriggerSchedulingCron triggers by matching a cron expression
	TriggerSchedulingCron TriggerSchedulingPolicy = "cron"
)

//TriggerSchedule represents a build trigger that fires on a time-bound schedule
type TriggerSchedule struct {
	triggerJSON *triggerJSON
	buildTypeID string

	SchedulingPolicy TriggerSchedulingPolicy `prop:"schedulingPolicy"`
	Rules            []string                `prop:"triggerRules" separator:"\n"`
	Timezone         string                  `prop:"timezone"`
	Hour             uint                    `prop:"hour"`
	Minute           uint                    `prop:"minute"`
	Weekday          time.Weekday
	Options          *TriggerScheduleOptions
}

//ID for this entity
func (t *TriggerSchedule) ID() string {
	return t.triggerJSON.ID
}

//Type returns TriggerTypes.Schedule ("schedulingTrigger")
func (t *TriggerSchedule) Type() string {
	return TriggerTypes.Schedule
}

//SetDisabled controls whether this trigger is disabled or not
func (t *TriggerSchedule) SetDisabled(disabled bool) {
	t.triggerJSON.Disabled = NewBool(disabled)
}

//Disabled gets the disabled status for this trigger
func (t *TriggerSchedule) Disabled() bool {
	return *t.triggerJSON.Disabled
}

//BuildTypeID gets the build type identifier
func (t *TriggerSchedule) BuildTypeID() string {
	return t.buildTypeID
}

//SetBuildTypeID sets the build type identifier
func (t *TriggerSchedule) SetBuildTypeID(id string) {
	t.buildTypeID = id
}

//NewTriggerScheduleDaily returns a TriggerSchedule that fires daily on the hour/minute specified
func NewTriggerScheduleDaily(sourceBuildID string, hour uint, minute uint, timezone string, rules []string) (*TriggerSchedule, error) {
	return NewTriggerSchedule(TriggerSchedulingDaily, sourceBuildID, time.Sunday, hour, minute, timezone, rules, NewTriggerScheduleOptions())
}

//NewTriggerScheduleWeekly returns a TriggerSchedule that fires weekly on the weekday and hour/minute specified
func NewTriggerScheduleWeekly(sourceBuildID string, weekday time.Weekday, hour uint, minute uint, timezone string, rules []string) (*TriggerSchedule, error) {
	return NewTriggerSchedule(TriggerSchedulingWeekly, sourceBuildID, weekday, hour, minute, timezone, rules, NewTriggerScheduleOptions())
}

//NewTriggerSchedule returns a TriggerSchedule with the scheduling policy and options specified
func NewTriggerSchedule(schedulingPolicy TriggerSchedulingPolicy, sourceBuildID string, weekday time.Weekday, hour uint, minute uint, timezone string, rules []string, opt *TriggerScheduleOptions) (*TriggerSchedule, error) {
	if hour > 23 {
		return nil, fmt.Errorf("invalid hour: %d, must be between 0-23", hour)
	}
	if minute > 59 {
		return nil, fmt.Errorf("invalid minute: %d, must be between 0-59", minute)
	}
	if weekday < time.Sunday || weekday > time.Saturday {
		return nil, fmt.Errorf("invalid weekday: %d, must be between time.Sunday and time.Saturday", weekday)
	}

	return &TriggerSchedule{
		SchedulingPolicy: schedulingPolicy,
		Timezone:         timezone,
		Rules:            rules,
		Weekday:          weekday,
		Hour:             hour,
		Minute:           minute,
		buildTypeID:      sourceBuildID,

		triggerJSON: &triggerJSON{
			Disabled: NewFalse(),
			Type:     TriggerTypes.Schedule,
		},

		Options: opt,
	}, nil
}

func (t *TriggerSchedule) read(dt *triggerJSON) error {
	if dt.Disabled != nil {
		t.SetDisabled(*dt.Disabled)
	}
	t.triggerJSON = dt

	if v, ok := dt.Properties.GetOk("schedulingPolicy"); ok {
		t.SchedulingPolicy = v
	} else {
		return fmt.Errorf("invalid 'schedulingPolicy' property")
	}

	if v, ok := dt.Properties.GetOk("triggerRules"); ok {
		t.Rules = strings.Split(v, "\n")
	}

	if v, ok := dt.Properties.GetOk("timezone"); ok {
		t.Timezone = v
	}

	t.Options = t.triggerJSON.Properties.triggerScheduleOptions()

	switch t.SchedulingPolicy {
	case TriggerSchedulingDaily, TriggerSchedulingWeekly:
		return t.readDailyOrWeekly(dt)
	}

	return nil
}

func (t *TriggerSchedule) readDailyOrWeekly(dt *triggerJSON) error {
	if v, ok := dt.Properties.GetOk("hour"); ok {
		p, _ := strconv.ParseUint(v, 10, 0)
		t.Hour = uint(p)
	} else {
		return fmt.Errorf("invalid 'hour' property")
	}
	if v, ok := dt.Properties.GetOk("minute"); ok {
		p, _ := strconv.ParseUint(v, 10, 0)
		t.Minute = uint(p)
	} else {
		return fmt.Errorf("invalid 'minute' property")
	}

	if t.SchedulingPolicy == TriggerSchedulingWeekly {
		if v, ok := dt.Properties.GetOk("dayOfWeek"); ok {
			w, err := parseWeekday(v)
			if err != nil {
				return fmt.Errorf("invalid 'dayOfWeek' property")
			}
			t.Weekday = w
		}
	}

	return nil
}

var daysOfWeek = [...]string{
	"Sunday",
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
}

func parseWeekday(v string) (time.Weekday, error) {
	for i := range daysOfWeek {
		if daysOfWeek[i] == v {
			return time.Weekday(i), nil
		}
	}

	return time.Sunday, fmt.Errorf("invalid weekday '%s'", v)
}

func (t *TriggerSchedule) properties() *Properties {
	props := serializeToProperties(t)
	optProps := t.Options.properties()
	if t.SchedulingPolicy == TriggerSchedulingWeekly {
		props.AddOrReplaceValue("dayOfWeek", t.Weekday.String())
	}
	props = props.Concat(optProps)
	return props
}

//MarshalJSON implements JSON serialization for TriggerSchedule
func (t *TriggerSchedule) MarshalJSON() ([]byte, error) {
	props := t.properties()
	optProps := t.Options.properties()
	for _, p := range optProps.Items {
		props.AddOrReplaceProperty(p)
	}

	out := &triggerJSON{
		ID:         t.ID(),
		Type:       t.Type(),
		Disabled:   NewBool(t.Disabled()),
		Properties: props,
	}

	return json.Marshal(out)
}

//UnmarshalJSON implements JSON deserialization for TriggerSchedule
func (t *TriggerSchedule) UnmarshalJSON(data []byte) error {
	var aux triggerJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Type != TriggerTypes.Schedule {
		return fmt.Errorf("invalid type %s trying to deserialize into TriggerSchedule entity", aux.Type)
	}

	if err := t.read(&aux); err != nil {
		return err
	}

	return nil
}

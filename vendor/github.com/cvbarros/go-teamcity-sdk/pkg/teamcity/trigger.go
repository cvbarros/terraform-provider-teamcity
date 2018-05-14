package teamcity

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
)

type Triggers struct {
	// count
	Count int32 `json:"count,omitempty" xml:"count"`

	// property
	Items []*Trigger `json:"trigger"`
}

// Trigger represents a build trigger to be associated with a build configuration. Use the constructor methods to create new instances.
type Trigger struct {

	// disabled
	Disabled *bool `json:"disabled,omitempty" xml:"disabled"`

	// href
	Href string `json:"href,omitempty" xml:"href"`

	// id
	ID string `json:"id,omitempty" xml:"id"`

	// inherited
	Inherited *bool `json:"inherited,omitempty" xml:"inherited"`

	// properties
	Properties *Properties `json:"properties,omitempty"`

	// type
	Type string `json:"type,omitempty" xml:"type"`
}

// NewVcsTrigger returns a VCS trigger type with the triggerRules specified. triggerRules is required, but branchFilter can be optional if the VCS root uses multiple branches.
func NewVcsTrigger(triggerRules string, branchFilter string) *Trigger {
	props := NewProperties(
		&Property{
			Name:  "triggerRules",
			Value: triggerRules,
		},
		&Property{
			Name:  "enableQueueOptimization",
			Value: "true",
		},
		&Property{
			Name:  "quietPeriodMode",
			Value: "DO_NOT_USE",
		},
	)

	if branchFilter != "" {
		props.Add(&Property{
			Name:  "branchFilter",
			Value: branchFilter,
		})
	}

	return &Trigger{
		Disabled:   NewFalse(),
		Type:       TriggerTypes.Vcs,
		Properties: props,
	}
}

// TriggerService provides operations for managing build triggers for a buildType
type TriggerService struct {
	BuildTypeID string
	httpClient  *http.Client
	base        *sling.Sling
}

func newTriggerService(buildTypeId string, c *http.Client, base *sling.Sling) *TriggerService {
	return &TriggerService{
		BuildTypeID: buildTypeId,
		httpClient:  c,
		base:        base.Path(fmt.Sprintf("buildTypes/%s/triggers/", Locator(buildTypeId).String())),
	}
}

//AddTrigger adds a new build trigger to a build type
func (s *TriggerService) AddTrigger(t *Trigger) error {
	var out Trigger
	if t == nil {
		return errors.New("t can't be nil")
	}

	_, err := s.base.New().Post("").BodyJSON(t).ReceiveSuccess(&out)
	if err != nil {
		return err
	}

	return nil
}

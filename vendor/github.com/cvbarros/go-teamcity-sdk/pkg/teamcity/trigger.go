package teamcity

import (
	"errors"
	"fmt"
	"io/ioutil"
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
	BuildTypeID string `json:"-"`

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

//Rules is a getter for triggerRules read-only property. No check performed since it's a required property.
func (t *Trigger) Rules() string {
	v, _ := t.Properties.GetOk("triggerRules")
	return v
}

//BranchFilterOk is a getter for branchFilter property. Returns false if propery doesnt exist
func (t *Trigger) BranchFilterOk() (string, bool) {
	if t.Properties == nil {
		return "", false
	}
	return t.Properties.GetOk("branchFilter")
}

//SetBranchFilter is s setter for branchFilter property
func (t *Trigger) SetBranchFilter(value string) {
	if t.Properties == nil {
		t.Properties = &Properties{Count: 0, Items: make([]*Property, 0)}
	}

	t.Properties.AddOrReplaceValue("branchFilter", value)
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
func (s *TriggerService) AddTrigger(t *Trigger) (*Trigger, error) {
	var out Trigger
	if t == nil {
		return nil, errors.New("t can't be nil")
	}

	resp, err := s.base.New().Post("").BodyJSON(t).ReceiveSuccess(&out)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error creating trigger for build type (id: %s)", s.BuildTypeID)
	}
	out.BuildTypeID = s.BuildTypeID
	return &out, nil
}

//GetById returns a dependency by its id
func (s *TriggerService) GetById(id string) (*Trigger, error) {
	var out Trigger
	resp, err := s.base.New().Get(id).ReceiveSuccess(&out)

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("404 Not Found - Trigger (id: %s) for buildTypeId (id: %s) was not found", id, s.BuildTypeID)
	}

	if err != nil {
		return nil, err
	}
	out.BuildTypeID = s.BuildTypeID
	return &out, nil
}

//Delete removes a snapshot dependency from the build configuration by its id
func (s *TriggerService) Delete(id string) error {
	request, _ := s.base.New().Delete(id).Request()
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	if response.StatusCode == 204 {
		return nil
	}

	if response.StatusCode != 200 && response.StatusCode != 204 {
		respData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Error '%d' when deleting trigger: %s", response.StatusCode, string(respData))
	}

	return nil
}

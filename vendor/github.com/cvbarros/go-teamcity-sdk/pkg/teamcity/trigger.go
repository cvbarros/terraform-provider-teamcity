package teamcity

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dghubble/sling"
)

//Triggers represents a typed, serializable collection of Trigger
type Triggers struct {
	// count
	Count int32 `json:"count,omitempty" xml:"count"`

	// property
	Items []*Trigger `json:"trigger"`
}

// Trigger represents a build trigger to be associated with a build configuration. Use the constructor methods to create new instances.
type Trigger struct {
	// build type id
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
	opt := NewVcsTriggerOptions()
	i, _ := NewVcsTriggerWithOptions(triggerRules, branchFilter, opt)
	return i
}

// NewVcsTriggerWithOptions returns a VCS trigger type with VcsTriggerOptions. See also NewVcsTrigger for other parameters.
func NewVcsTriggerWithOptions(triggerRules string, branchFilter string, opt *VcsTriggerOptions) (*Trigger, error) {
	if opt == nil {
		return nil, fmt.Errorf("opt parameter must be valid VcsTriggerOptions, not nil")
	}

	props := NewProperties(
		&Property{
			Name:  "triggerRules",
			Value: triggerRules,
		},
	)

	optProps := opt.vcsTriggerProperties()
	for _, p := range optProps.Items {
		props.AddOrReplaceProperty(p)
	}

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
	}, nil
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

func newTriggerService(buildTypeID string, c *http.Client, base *sling.Sling) *TriggerService {
	return &TriggerService{
		BuildTypeID: buildTypeID,
		httpClient:  c,
		base:        base.Path(fmt.Sprintf("buildTypes/%s/triggers/", Locator(buildTypeID).String())),
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

//GetByID returns a dependency by its id
func (s *TriggerService) GetByID(id string) (*Trigger, error) {
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
	response, err := s.httpClient.Do(request)
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

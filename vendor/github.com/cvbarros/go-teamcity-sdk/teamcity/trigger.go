package teamcity

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dghubble/sling"
)

//Triggers represents a typed, serializable collection of Trigger
type Triggers struct {
	// count
	Count int32 `json:"count,omitempty" xml:"count"`

	// property
	Items []*Trigger `json:"trigger"`
}

// TriggerService provides operations for managing build triggers for a buildType
type TriggerService struct {
	BuildTypeID string
	httpClient  *http.Client
	base        *sling.Sling
	restHelper  *restHelper
}

func newTriggerService(buildTypeID string, c *http.Client, base *sling.Sling) *TriggerService {
	sling := base.Path(fmt.Sprintf("buildTypes/%s/triggers/", Locator(buildTypeID).String()))
	return &TriggerService{
		BuildTypeID: buildTypeID,
		httpClient:  c,
		base:        sling,
		restHelper:  newRestHelperWithSling(c, sling),
	}
}

//AddTrigger adds a new build trigger to a build type
func (s *TriggerService) AddTrigger(t Trigger) (Trigger, error) {
	var created Trigger
	err := s.restHelper.postCustom("", t, &created, "build trigger", triggerReadingFunc)
	if err != nil {
		//Duplicate vcsTrigger for the buildConfiguration - Can't add more than one vcsTrigger
		if strings.Contains(err.Error(), "Trigger with id 'vcsTrigger'already exists") {
			return nil, fmt.Errorf("unable to add two VCS triggers to build configuration")
		}
		return nil, err
	}

	created.SetBuildTypeID(s.BuildTypeID)
	return created, nil
}

//GetByID returns a build trigger by its id
func (s *TriggerService) GetByID(id string) (Trigger, error) {
	var out Trigger
	err := s.restHelper.getCustom(id, &out, "build trigger", triggerReadingFunc)

	if err != nil {
		return nil, err
	}
	out.SetBuildTypeID(s.BuildTypeID)
	return out, nil
}

//Delete removes a build trigger from the build configuration by its id
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

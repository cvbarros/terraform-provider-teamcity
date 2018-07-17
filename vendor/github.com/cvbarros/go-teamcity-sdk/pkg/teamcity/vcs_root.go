package teamcity

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dghubble/sling"
)

// VcsRoot represents a detailed VCS Root entity
type VcsRoot struct {

	// href
	Href string `json:"href,omitempty" xml:"href"`

	// id
	ID string `json:"id,omitempty" xml:"id"`

	// internal Id
	InternalID string `json:"internalId,omitempty" xml:"internalId"`

	// locator
	Locator string `json:"locator,omitempty" xml:"locator"`

	// modification check interval
	ModificationCheckInterval int32 `json:"modificationCheckInterval,omitempty" xml:"modificationCheckInterval"`

	// name
	Name string `json:"name,omitempty" xml:"name"`

	// project
	Project *ProjectReference `json:"project,omitempty"`

	// project locator
	ProjectLocator string `json:"projectLocator,omitempty" xml:"projectLocator"`

	// properties
	Properties *Properties `json:"properties,omitempty"`

	// uuid
	UUID string `json:"uuid,omitempty" xml:"uuid"`

	// vcs name
	VcsName string `json:"vcsName,omitempty" xml:"vcsName"`
}

// VcsRootReference represents a subset detail of a VCS Root
type VcsRootReference struct {

	// href
	Href string `json:"href,omitempty" xml:"href"`

	// id
	ID string `json:"id,omitempty" xml:"id"`

	// name
	Name string `json:"name,omitempty" xml:"name"`

	// project Id
	Project *ProjectReference `json:"project,omitempty" xml:"project"`
}

// VcsRootService has operations for handling vcs roots
type VcsRootService struct {
	sling      *sling.Sling
	httpClient *http.Client
}

func newVcsRootService(base *sling.Sling, httpClient *http.Client) *VcsRootService {
	return &VcsRootService{
		sling:      base.Path("vcs-roots/"),
		httpClient: httpClient,
	}
}

// Create creates a new vcs root
func (s *VcsRootService) Create(projectID string, vcsRoot *VcsRoot) (*VcsRootReference, error) {
	var created VcsRootReference

	success, err := s.Validate(projectID, vcsRoot)
	if success == false {
		return nil, err
	}

	_, err = s.sling.New().Post("").BodyJSON(vcsRoot).ReceiveSuccess(&created)

	if err != nil {
		return nil, err
	}

	return &created, nil
}

// GetByID Retrieves a vcs root by id using the id: locator
func (s *VcsRootService) GetByID(id string) (*VcsRoot, error) {
	var out VcsRoot

	resp, err := s.sling.New().Get(id).ReceiveSuccess(&out)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error when retrieving VcsRoot id = '%s', status: %d", id, resp.StatusCode)
	}

	return &out, err
}

//Delete a VCS Root resource using id: locator
func (s *VcsRootService) Delete(id string) error {
	request, _ := s.sling.New().Delete(id).Request()

	//TODO: Expose the same httpClient used by sling
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
		return fmt.Errorf("Error '%d' when deleting vcsRoot: %s", response.StatusCode, string(respData))
	}

	return nil
}

// Validate verifies if a VcsRoot model is valid for updating/creation before sending to upstream API.
func (s *VcsRootService) Validate(projectID string, vcsRoot *VcsRoot) (bool, error) {
	if vcsRoot == nil {
		return false, errors.New("vcsRoot must not be nil")
	}
	if vcsRoot.Project == nil {
		return false, errors.New("vcsRoot.Project must not be nil")
	}
	if vcsRoot.VcsName == "" {
		return false, errors.New("vcsRoot.VcsName must be defined")
	}

	props := vcsRoot.Properties.Map()
	if _, ok := props["url"]; !ok {
		return false, errors.New("'url' property must be defined in VcsRoot.Properties")
	}

	if _, ok := props["branch"]; !ok {
		return false, errors.New("'branch' property must be defined in VcsRoot.Properties")
	}
	return true, nil
}

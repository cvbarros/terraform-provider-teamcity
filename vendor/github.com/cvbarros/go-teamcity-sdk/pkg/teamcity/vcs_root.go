package teamcity

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dghubble/sling"
)

//VcsRoot interface represents a base type of VCSRoot
type VcsRoot interface {
	//GetID returns the ID of this VCS Root
	GetID() string

	//VcsName returns the type of VCS Root. See VcsNames for possible values returned.
	//In addition, this can be used to type assert to the appropriate concrete VCS Root type.
	VcsName() string

	//Properties returns the Properties collection for this VCS Root. This should be used for querying only.
	Properties() *Properties
}

type vcsRootJSON struct {
	// id
	ID string `json:"id,omitempty" xml:"id"`

	// internal Id
	InternalID string `json:"internalId,omitempty" xml:"internalId"`

	// ModificationCheckInterval value in seconds to override the global server setting.
	ModificationCheckInterval int32 `json:"modificationCheckInterval,omitempty" xml:"modificationCheckInterval"`

	// name
	Name string `json:"name,omitempty" xml:"name"`

	// project
	Project *ProjectReference `json:"project,omitempty"`

	// Properties for the VCS Root. Do not set directly, instead use NewVcsRoot... constructors.
	Properties *Properties `json:"properties,omitempty"`

	// VcsName is the VCS Type used for this VCS Root. See VcsNames for allowed values.
	// Use NewVcsRoot... constructors to avoid setting this directly.
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
	restHelper *restHelper
}

func newVcsRootService(base *sling.Sling, httpClient *http.Client) *VcsRootService {
	return &VcsRootService{
		sling:      base.Path("vcs-roots/"),
		httpClient: httpClient,
		restHelper: newRestHelper(httpClient),
	}
}

// Create creates a new vcs root
func (s *VcsRootService) Create(projectID string, vcsRoot VcsRoot) (*VcsRootReference, error) {
	var created VcsRootReference

	err := s.restHelper.postJSONWithSling("", s.sling, vcsRoot, &created, "VcsRoot")

	if err != nil {
		return nil, err
	}

	return &created, nil
}

// GetByID Retrieves a vcs root by id using the id: locator
func (s *VcsRootService) GetByID(id string) (VcsRoot, error) {
	req, err := s.sling.New().Get(id).Request()

	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error when retrieving VcsRoot id = '%s', status: %d", id, resp.StatusCode)
	}

	return s.readVcsRootResponse(resp)
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

func (s *VcsRootService) readVcsRootResponse(resp *http.Response) (VcsRoot, error) {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var payload vcsRootJSON
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		return nil, err
	}

	var out VcsRoot
	switch payload.VcsName {
	case VcsNames.Git:
		var git GitVcsRoot
		if err := git.UnmarshalJSON(bodyBytes); err != nil {
			return nil, err
		}
		out = &git
	default:
		return nil, fmt.Errorf("Unsupported VCS Root type: '%s' (id:'%s') for projectID: %s", payload.VcsName, payload.ID, payload.Project.ID)
	}

	return out, nil
}

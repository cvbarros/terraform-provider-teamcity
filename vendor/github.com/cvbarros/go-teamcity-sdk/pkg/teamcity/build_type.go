package teamcity

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dghubble/sling"
)

// BuildType represents a build configuration or a build configuration template
type BuildType struct {

	// agent requirements
	AgentRequirements *AgentRequirements `json:"agent-requirements,omitempty"`

	// // artifact dependencies
	// ArtifactDependencies *ArtifactDependencies `json:"artifact-dependencies,omitempty"`

	// description
	Description string `json:"description,omitempty" xml:"description"`

	// href
	Href string `json:"href,omitempty" xml:"href"`

	// id
	ID string `json:"id,omitempty" xml:"id"`

	// inherited
	// Inherited *bool `json:"inherited,omitempty" xml:"inherited"`

	// internal Id
	InternalID string `json:"internalId,omitempty" xml:"internalId"`

	// locator
	Locator string `json:"locator,omitempty" xml:"locator"`

	// name
	Name string `json:"name,omitempty" xml:"name"`

	// Parameters for the build configuration. Read-only, only useful when retrieving project details
	Parameters *Properties `json:"parameters,omitempty"`

	// paused
	Paused *bool `json:"paused,omitempty" xml:"paused"`

	// project
	Project *Project `json:"project,omitempty"`

	// project Id
	ProjectID string `json:"projectId,omitempty" xml:"projectId"`

	// project internal Id
	ProjectInternalID string `json:"projectInternalId,omitempty" xml:"projectInternalId"`

	// project name
	ProjectName string `json:"projectName,omitempty" xml:"projectName"`

	// settings
	Settings *Properties `json:"settings,omitempty"`

	// snapshot dependencies
	SnapshotDependencies *SnapshotDependencies `json:"snapshot-dependencies,omitempty"`

	// steps
	Steps *Steps `json:"steps,omitempty"`

	// template flag
	TemplateFlag *bool `json:"templateFlag,omitempty" xml:"templateFlag"`

	// triggers
	Triggers *Triggers `json:"triggers,omitempty"`

	// type
	Type string `json:"type,omitempty" xml:"type"`

	// uuid
	UUID string `json:"uuid,omitempty" xml:"uuid"`

	// vcs root entries
	VcsRootEntries *VcsRootEntries `json:"vcs-root-entries,omitempty"`

	// web Url
	WebURL string `json:"webUrl,omitempty" xml:"webUrl"`
}

// BuildTypeReference represents a subset detail of a Build Type
type BuildTypeReference struct {

	// href
	Href string `json:"href,omitempty" xml:"href"`

	// id
	ID string `json:"id,omitempty" xml:"id"`

	// name
	Name string `json:"name,omitempty" xml:"name"`

	// project Id
	ProjectID string `json:"projectId,omitempty" xml:"projectId"`

	// project name
	ProjectName string `json:"projectName,omitempty" xml:"projectName"`
}

// Reference converts a BuildType entity to a BuildType reference
func (b *BuildType) Reference() *BuildTypeReference {
	return &BuildTypeReference{
		ID:          b.ID,
		Name:        b.Name,
		ProjectID:   b.ProjectID,
		Href:        b.Href,
		ProjectName: b.ProjectName,
	}
}

// BuildTypeService has operations for handling build configurations and templates
type BuildTypeService struct {
	sling *sling.Sling
}

func newBuildTypeService(base *sling.Sling) *BuildTypeService {
	return &BuildTypeService{
		sling: base.Path("buildTypes/"),
	}
}

// Create Creates a new build type under a project
func (s *BuildTypeService) Create(projectId string, buildType *BuildType) (*BuildTypeReference, error) {
	var created BuildTypeReference

	_, err := s.sling.New().Post("").BodyJSON(buildType).ReceiveSuccess(&created)

	if err != nil {
		return nil, err
	}

	return &created, nil
}

// GetById Retrieves a build type resource by ID
func (s *BuildTypeService) GetById(id string) (*BuildType, error) {
	var out BuildType

	resp, err := s.sling.New().Get(id).ReceiveSuccess(&out)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error when retrieving BuildType id = '%s', status: %d", id, resp.StatusCode)
	}

	return &out, err
}

//Delete a build type resource
func (s *BuildTypeService) Delete(id string) error {
	request, _ := s.sling.New().Delete(id).Request()
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
		return fmt.Errorf("Error '%d' when deleting build type: %s", response.StatusCode, string(respData))
	}

	return nil
}

// AttachVcsRoot adds the VcsRoot reference to this build type
func (s *BuildTypeService) AttachVcsRoot(id string, vcsRoot *VcsRootReference) error {
	var vcsEntry = NewVcsRootEntry(vcsRoot)
	return s.AttachVcsRootEntry(id, vcsEntry)
}

// AttachVcsRootEntry adds the VcsRootEntry to this build type
func (s *BuildTypeService) AttachVcsRootEntry(id string, entry *VcsRootEntry) error {
	var created VcsRootEntry
	_, err := s.sling.New().Post(fmt.Sprintf("%s/vcs-root-entries/", LocatorId(id))).BodyJSON(entry).ReceiveSuccess(&created)

	if err != nil {
		return err
	}

	return nil
}

// AddStep creates a new build steo for this build type
func (s *BuildTypeService) AddStep(id string, step *Step) error {
	var created Step
	_, err := s.sling.New().Post(fmt.Sprintf("%s/steps/", LocatorId(id))).BodyJSON(step).ReceiveSuccess(&created)

	if err != nil {
		return err
	}

	return nil
}

// UpdateSettings will do a remote call for each setting being updated. Operation is not atomic, and the list of settings is processed in the order sent.
// Will return the error of the first failure and not process the rest
func (s *BuildTypeService) UpdateSettings(id string, settings *Properties) error {
	for _, item := range settings.Items {
		bodyProvider := textPlainBodyProvider{payload: item.Value}
		req, err := s.sling.New().Put(fmt.Sprintf("%s/settings/%s", LocatorId(id), item.Name)).BodyProvider(bodyProvider).Add("Accept", "text/plain").Request()
		response, err := http.DefaultClient.Do(req)
		response.Body.Close()
		if err != nil {
			return fmt.Errorf("error updating buildType id: '%s' setting '%s': %s", id, item.Name, err)
		}
	}

	return nil
}

//DeleteStep removes a build step from this build type by its id
func (s *BuildTypeService) DeleteStep(id string, stepId string) error {
	_, err := s.sling.New().Delete(fmt.Sprintf("%s/steps/%s", LocatorId(id), stepId)).ReceiveSuccess(nil)

	if err != nil {
		return err
	}

	return nil
}

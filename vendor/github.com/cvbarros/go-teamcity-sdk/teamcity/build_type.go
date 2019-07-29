package teamcity

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dghubble/sling"
)

type buildTypeJSON struct {
	Description          string                `json:"description,omitempty" xml:"description"`
	Href                 string                `json:"href,omitempty" xml:"href"`
	ID                   string                `json:"id,omitempty" xml:"id"`
	InternalID           string                `json:"internalId,omitempty" xml:"internalId"`
	Locator              string                `json:"locator,omitempty" xml:"locator"`
	Name                 string                `json:"name,omitempty" xml:"name"`
	Parameters           *Parameters           `json:"parameters,omitempty"`
	Paused               *bool                 `json:"paused,omitempty" xml:"paused"`
	Project              *Project              `json:"project,omitempty"`
	ProjectID            string                `json:"projectId,omitempty" xml:"projectId"`
	ProjectInternalID    string                `json:"projectInternalId,omitempty" xml:"projectInternalId"`
	ProjectName          string                `json:"projectName,omitempty" xml:"projectName"`
	SnapshotDependencies *SnapshotDependencies `json:"snapshot-dependencies,omitempty"`
	TemplateFlag         *bool                 `json:"templateFlag,omitempty" xml:"templateFlag"`
	Type                 string                `json:"type,omitempty" xml:"type"`
	UUID                 string                `json:"uuid,omitempty" xml:"uuid"`
	Settings             *Properties           `json:"settings,omitempty"`
	Templates            *Templates            `json:"templates,omitempty"`
	Steps                *stepsJSON            `json:"steps,omitempty"`
	VcsRootEntries       *VcsRootEntries       `json:"vcs-root-entries,omitempty"`
	WebURL               string                `json:"webUrl,omitempty" xml:"webUrl"`

	// inherited
	// Inherited *bool `json:"inherited,omitempty" xml:"inherited"`
}

// Templates represents a collection of BuildTypeReference that are templates attached to a build configuration.
type Templates struct {

	// count
	Count int32 `json:"count,omitempty" xml:"count"`

	// buildType
	Items []*BuildTypeReference `json:"buildType"`
}

// BuildType represents a build configuration or a build configuration template
type BuildType struct {
	ProjectID   string
	ID          string
	Name        string
	Description string
	Options     *BuildTypeOptions
	Disabled    bool
	IsTemplate  bool
	Steps       []Step
	Templates   *Templates

	VcsRootEntries []*VcsRootEntry
	Parameters     *Parameters
	buildTypeJSON  *buildTypeJSON
}

//NewBuildType returns a build configuration with default options
func NewBuildType(projectID string, name string) (*BuildType, error) {
	if projectID == "" || name == "" {
		return nil, fmt.Errorf("projectID and name are required")
	}

	opt := NewBuildTypeOptionsWithDefaults()
	return &BuildType{
		ProjectID:  projectID,
		Name:       name,
		Options:    opt,
		Parameters: NewParametersEmpty(),
		IsTemplate: false,
		buildTypeJSON: &buildTypeJSON{
			ProjectID: projectID,
			Settings:  opt.properties(),
		},
		Steps: []Step{},
	}, nil
}

//NewBuildTypeTemplate returns a build configuration template with default options
func NewBuildTypeTemplate(projectID string, name string) (*BuildType, error) {
	if projectID == "" || name == "" {
		return nil, fmt.Errorf("projectID and name are required")
	}

	opt := NewBuildTypeOptionsTemplate()
	return &BuildType{
		ProjectID:  projectID,
		Name:       name,
		Options:    opt,
		IsTemplate: true,
		Parameters: NewParametersEmpty(),
		buildTypeJSON: &buildTypeJSON{
			ProjectID: projectID,
			Settings:  opt.properties(),
		},
	}, nil
}

//MarshalJSON implements JSON serialization for BuildType
func (b *BuildType) MarshalJSON() ([]byte, error) {
	optProps := b.Options.properties()

	out := &buildTypeJSON{
		ID:           b.ID,
		ProjectID:    b.ProjectID,
		Name:         b.Name,
		Settings:     optProps,
		Parameters:   b.Parameters,
		TemplateFlag: NewBool(b.IsTemplate),
		Templates:    b.Templates,
	}

	//TODO: TeamCity API doesn't support "description" property if creating a template. Need to manually update it after, like projects.
	if !b.IsTemplate {
		out.Description = b.Description
	}
	if len(b.Steps) > 0 {
		out.Steps = b.serializeSteps()
	}
	return json.Marshal(out)
}

//UnmarshalJSON implements JSON deserialization for TriggerSchedule
func (b *BuildType) UnmarshalJSON(data []byte) error {
	var aux buildTypeJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if err := b.read(&aux); err != nil {
		return err
	}

	return nil
}

func (b *BuildType) read(dt *buildTypeJSON) error {
	var isTemplate bool
	if dt.TemplateFlag != nil {
		isTemplate = *dt.TemplateFlag
		b.IsTemplate = isTemplate
	}
	opts := dt.Settings.buildTypeOptions(isTemplate)

	b.ID = dt.ID
	b.Name = dt.Name
	b.Description = dt.Description
	b.Options = opts
	b.ProjectID = dt.ProjectID
	b.VcsRootEntries = dt.VcsRootEntries.Items
	b.Parameters = dt.Parameters
	b.Templates = dt.Templates

	steps := make([]Step, dt.Steps.Count)
	for i := range dt.Steps.Items {
		dt, err := json.Marshal(dt.Steps.Items[i])
		if err != nil {
			return err
		}
		stepReadingFunc(dt, &steps[i])
	}
	b.Steps = steps

	return nil
}

func (b *BuildType) serializeSteps() *stepsJSON {
	out := &stepsJSON{Count: int32(len(b.Steps)), Items: make([]*stepJSON, len(b.Steps))}
	for i := 0; i < len(b.Steps); i++ {
		out.Items[i] = b.Steps[i].serializable()
	}
	return out
}

// BuildTypeReference represents a subset detail of a Build Type
type BuildTypeReference struct {

	// id
	ID string `json:"id,omitempty" xml:"id"`

	// name
	Name string `json:"name,omitempty" xml:"name"`

	// project Id
	ProjectID string `json:"projectId,omitempty" xml:"projectId"`
}

// BuildTypeReferences represents a collection of *BuildTypeReference
type BuildTypeReferences struct {
	// count
	Count int32 `json:"count,omitempty" xml:"count"`

	// buildType
	Items []*BuildTypeReference `json:"buildType"`
}

// Reference converts a BuildType entity to a BuildType reference
func (b *BuildType) Reference() *BuildTypeReference {
	return &BuildTypeReference{
		ID:        b.ID,
		Name:      b.Name,
		ProjectID: b.ProjectID,
	}
}

// BuildTypeService has operations for handling build configurations and templates
type BuildTypeService struct {
	sling      *sling.Sling
	httpClient *http.Client
	restHelper *restHelper
}

func newBuildTypeService(base *sling.Sling, httpClient *http.Client) *BuildTypeService {
	sling := base.Path("buildTypes/")
	return &BuildTypeService{
		httpClient: httpClient,
		sling:      sling,
		restHelper: newRestHelperWithSling(httpClient, sling),
	}
}

// Create Creates a new build type under a project
func (s *BuildTypeService) Create(projectID string, buildType *BuildType) (*BuildTypeReference, error) {
	var created BuildTypeReference

	err := s.restHelper.post("", buildType, &created, "Build Type")

	if err != nil {
		return nil, err
	}

	return &created, nil
}

// GetByID Retrieves a build type resource by ID
func (s *BuildTypeService) GetByID(id string) (*BuildType, error) {
	var out BuildType

	resp, err := s.sling.New().Get(id).ReceiveSuccess(&out)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error when retrieving BuildType id = '%s', status: %d", id, resp.StatusCode)
	}

	//For now, filter all inherited parameters, until figuring out a proper way of exposing filtering options to the caller
	out.Parameters = out.Parameters.NonInherited()

	return &out, err
}

//Update changes the resource in-place for this build configuration.
//TeamCity API does not support "PUT" on the whole Build Configuration resource, so the only updateable field is "Description". Other field updates will be ignored.
//This method also updates Settings and Parameters, but this is not an atomic operation. If an error occurs, it will be returned to caller what was updated or not.
func (s *BuildTypeService) Update(buildType *BuildType) (*BuildType, error) {
	_, err := s.restHelper.putTextPlain(buildType.ID+"/description", buildType.Description, "build type description")

	if err != nil {
		return nil, err
	}

	//Update settings
	var settings BuildTypeOptions
	err = s.restHelper.put(buildType.ID+"/settings", buildType.Options.properties(), &settings, "build type settings")
	if err != nil {
		return nil, err
	}

	//Update Parameters
	var parameters *Properties
	err = s.restHelper.put(buildType.ID+"/parameters", buildType.Parameters, &parameters, "build type parameters")
	if err != nil {
		return nil, err
	}

	//Update Steps
	if buildType.Steps != nil && len(buildType.Steps) > 0 {
		var steps []Step
		err = s.restHelper.putCustom(buildType.ID+"/steps", buildType.serializeSteps(), &steps, "build type steps", stepsReadingFunc)
		if err != nil {
			return nil, err
		}
	}

	out, err := s.GetByID(buildType.ID) //Refresh after update
	if err != nil {
		return nil, err
	}

	return out, nil
}

//Delete a build type resource
func (s *BuildTypeService) Delete(id string) error {
	request, _ := s.sling.New().Delete(id).Request()
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
	_, err := s.sling.New().Post(fmt.Sprintf("%s/vcs-root-entries/", LocatorID(id))).BodyJSON(entry).ReceiveSuccess(&created)

	if err != nil {
		return err
	}

	return nil
}

// AddStep creates a new build step for the build configuration with given id.
func (s *BuildTypeService) AddStep(id string, step Step) (Step, error) {
	var created Step
	path := fmt.Sprintf("%s/steps/", LocatorID(id))

	err := s.restHelper.postCustom(path, step, &created, "build step", stepReadingFunc)
	if err != nil {
		return nil, err
	}

	return created, nil
}

//GetSteps return the list of steps for a Build configuration with given id.
func (s *BuildTypeService) GetSteps(id string) ([]Step, error) {
	var aux stepsJSON
	path := fmt.Sprintf("%s/steps/", LocatorID(id))
	err := s.restHelper.get(path, &aux, "build steps")
	if err != nil {
		return nil, err
	}
	steps := make([]Step, aux.Count)
	for i := range aux.Items {
		dt, err := json.Marshal(aux.Items[i])
		if err != nil {
			return nil, err
		}
		stepReadingFunc(dt, &steps[i])
	}

	return steps, nil
}

// UpdateSettings will do a remote call for each setting being updated. Operation is not atomic, and the list of settings is processed in the order sent.
// Will return the error of the first failure and not process the rest
func (s *BuildTypeService) UpdateSettings(id string, settings *Properties) error {
	for _, item := range settings.Items {
		bodyProvider := textPlainBodyProvider{payload: item.Value}
		req, err := s.sling.New().Put(fmt.Sprintf("%s/settings/%s", LocatorID(id), item.Name)).BodyProvider(bodyProvider).Add("Accept", "text/plain").Request()
		response, err := s.httpClient.Do(req)
		response.Body.Close()
		if err != nil {
			return fmt.Errorf("error updating buildType id: '%s' setting '%s': %s", id, item.Name, err)
		}
	}

	return nil
}

//DeleteStep removes a build step from this build type by its id
func (s *BuildTypeService) DeleteStep(id string, stepID string) error {
	_, err := s.sling.New().Delete(fmt.Sprintf("%s/steps/%s", LocatorID(id), stepID)).ReceiveSuccess(nil)

	if err != nil {
		return err
	}

	return nil
}

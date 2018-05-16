package teamcity

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dghubble/sling"
)

//DependencyService provides operations for managing dependencies for a buildType
type DependencyService struct {
	BuildTypeID   string
	httpClient    *http.Client
	artifactSling *sling.Sling
	snapshotSling *sling.Sling
}

//NewDependencyService constructs and instance of DependencyService scoped to a given buildTypeId
func NewDependencyService(buildTypeId string, c *http.Client, base *sling.Sling) *DependencyService {
	return &DependencyService{
		BuildTypeID:   buildTypeId,
		httpClient:    c,
		artifactSling: base.New().Path(fmt.Sprintf("buildTypes/%s/artifact-dependencies/", buildTypeId)),
		snapshotSling: base.New().Path(fmt.Sprintf("buildTypes/%s/snapshot-dependencies/", buildTypeId)),
	}
}

//AddSnapshotDependency adds a new snapshot dependency to build type
func (s *DependencyService) AddSnapshotDependency(dep *SnapshotDependency) (*SnapshotDependency, error) {
	var out SnapshotDependency
	if dep == nil {
		return nil, errors.New("dep can't be nil")
	}

	resp, err := s.snapshotSling.New().Post("").BodyJSON(dep).ReceiveSuccess(&out)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Unknown error when adding snapshot dependency, statusCode: %d", resp.StatusCode)
	}
	out.BuildTypeID = s.BuildTypeID
	return &out, nil
}

//GetById returns a dependency by its id
func (s *DependencyService) GetById(depId string) (*SnapshotDependency, error) {
	var out SnapshotDependency
	resp, err := s.snapshotSling.New().Get(depId).ReceiveSuccess(&out)

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("404 Not Found - Snapshot dependency (id: %s) for buildTypeId (id: %s) was not found", depId, s.BuildTypeID)
	}

	if err != nil {
		return nil, err
	}
	out.BuildTypeID = s.BuildTypeID
	return &out, nil
}

//Delete removes a snapshot dependency from the build configuration by its id
func (s *DependencyService) Delete(depId string) error {
	request, _ := s.snapshotSling.New().Delete(depId).Request()
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
		return fmt.Errorf("Error '%d' when deleting snapshot dependency: %s", response.StatusCode, string(respData))
	}

	return nil
}

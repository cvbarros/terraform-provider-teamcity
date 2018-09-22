package teamcity

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
)

//DependencyService provides operations for managing dependencies for a buildType
type DependencyService struct {
	BuildTypeID   string
	httpClient    *http.Client
	artifactSling *sling.Sling
	snapshotSling *sling.Sling

	artifactHelper *restHelper
	snapshotHelper *restHelper
}

//NewDependencyService constructs and instance of DependencyService scoped to a given buildTypeId
func NewDependencyService(buildTypeID string, c *http.Client, base *sling.Sling) *DependencyService {
	artifactSling := base.New().Path(fmt.Sprintf("buildTypes/%s/artifact-dependencies/", buildTypeID))
	snapshotSling := base.New().Path(fmt.Sprintf("buildTypes/%s/snapshot-dependencies/", buildTypeID))
	return &DependencyService{
		BuildTypeID:    buildTypeID,
		httpClient:     c,
		artifactSling:  artifactSling,
		snapshotSling:  snapshotSling,
		artifactHelper: newRestHelperWithSling(c, artifactSling),
		snapshotHelper: newRestHelperWithSling(c, snapshotSling),
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

//AddArtifactDependency adds a new artifact dependency to build type
func (s *DependencyService) AddArtifactDependency(dep *ArtifactDependency) (*ArtifactDependency, error) {
	var out ArtifactDependency
	if dep == nil {
		return nil, errors.New("dep can't be nil")
	}

	resp, err := s.artifactSling.New().Post("").BodyJSON(dep).ReceiveSuccess(&out)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Unknown error when adding artifact dependency, statusCode: %d", resp.StatusCode)
	}

	out.SetBuildTypeID(s.BuildTypeID)
	return &out, nil
}

//GetSnapshotByID returns a snapshot dependency by its id
func (s *DependencyService) GetSnapshotByID(depID string) (*SnapshotDependency, error) {
	var out SnapshotDependency
	resp, err := s.snapshotSling.New().Get(depID).ReceiveSuccess(&out)

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("404 Not Found - Snapshot dependency (id: %s) for buildTypeId (id: %s) was not found", depID, s.BuildTypeID)
	}

	if err != nil {
		return nil, err
	}
	out.BuildTypeID = s.BuildTypeID
	return &out, nil
}

//GetArtifactByID returns an artifact dependency by its id
func (s *DependencyService) GetArtifactByID(depID string) (*ArtifactDependency, error) {
	var out ArtifactDependency
	err := s.artifactHelper.get(depID, &out, "artifact dependency")

	if err != nil {
		return nil, err
	}
	out.SetBuildTypeID(s.BuildTypeID)
	return &out, nil
}

//DeleteSnapshot removes a snapshot dependency from the build configuration by its id
func (s *DependencyService) DeleteSnapshot(depID string) error {
	return s.snapshotHelper.deleteByIDWithSling(s.snapshotSling, depID, "snapshot dependency")
}

//DeleteArtifact removes an artifact dependency from the build configuration by its id
func (s *DependencyService) DeleteArtifact(depID string) error {
	return s.artifactHelper.deleteByIDWithSling(s.artifactSling, depID, "artifact dependency")
}

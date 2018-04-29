package teamcity

import (
	"errors"
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

//AddSnapshotDependency adds a new snapshot dependency to build type
func (s *DependencyService) AddSnapshotDependency(dep *SnapshotDependency) error {
	var out SnapshotDependency
	if dep == nil {
		return errors.New("dep can't be nil")
	}

	_, err := s.snapshotSling.New().Post("").BodyJSON(dep).ReceiveSuccess(&out)
	if err != nil {
		return err
	}

	return nil
}

package teamcity

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
)

//BuildTemplateService provides operations for managing attaching/detaching build configuration templates to/from build configurations
type BuildTemplateService struct {
	BuildTypeID string
	httpClient  *http.Client

	restHelper *restHelper
}

//NewBuildTemplateService constructs an instance of NewBuildTemplateService scoped to a given buildTypeId
func NewBuildTemplateService(buildTypeID string, c *http.Client, base *sling.Sling) *BuildTemplateService {
	sling := base.New().Path(fmt.Sprintf("buildTypes/%s/templates/", buildTypeID))
	return &BuildTemplateService{
		BuildTypeID: buildTypeID,
		httpClient:  c,
		restHelper:  newRestHelperWithSling(c, sling),
	}
}

//Attach is an idempotent operation that attaches the build template with given ID to the build configuration fo this service.
func (s *BuildTemplateService) Attach(buildTemplateID string) (*BuildTypeReference, error) {
	var out BuildTypeReference
	dt := &BuildTypeReference{
		ID: buildTemplateID,
	}
	err := s.restHelper.post("", dt, &out, "attach build template")

	if err != nil {
		return nil, err
	}

	return &out, nil
}

//Detach disassociates the build template with given ID from the build configuration fo this service.
func (s *BuildTemplateService) Detach(buildTemplateID string) error {
	return s.restHelper.delete(buildTemplateID, "detach build template")
}

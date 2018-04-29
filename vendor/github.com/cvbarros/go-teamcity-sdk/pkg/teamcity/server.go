package teamcity

import "github.com/dghubble/sling"

// Server holds information about the TeamCity server
type Server struct {

	// agent pools
	AgentPools string `json:"agentPools,omitempty"`

	// agents
	Agents string `json:"agents,omitempty"`

	// build date
	BuildDate string `json:"buildDate,omitempty" xml:"buildDate"`

	// build number
	BuildNumber string `json:"buildNumber,omitempty" xml:"buildNumber"`

	// build queue
	BuildQueue string `json:"buildQueue,omitempty"`

	// builds
	Builds string `json:"builds,omitempty"`

	// current time
	CurrentTime string `json:"currentTime,omitempty" xml:"currentTime"`

	// internal Id
	InternalID string `json:"internalId,omitempty" xml:"internalId"`

	// investigations
	Investigations string `json:"investigations,omitempty"`

	// mutes
	Mutes string `json:"mutes,omitempty"`

	// projects
	Projects string `json:"projects,omitempty"`

	// role
	Role string `json:"role,omitempty" xml:"role"`

	// start time
	StartTime string `json:"startTime,omitempty" xml:"startTime"`

	// user groups
	UserGroups string `json:"userGroups,omitempty"`

	// users
	Users string `json:"users,omitempty"`

	// vcs roots
	VcsRoots string `json:"vcsRoots,omitempty"`

	// version
	Version string `json:"version,omitempty" xml:"version"`

	// version major
	VersionMajor int32 `json:"versionMajor,omitempty" xml:"versionMajor"`

	// version minor
	VersionMinor int32 `json:"versionMinor,omitempty" xml:"versionMinor"`

	// web Url
	WebURL string `json:"webUrl,omitempty" xml:"webUrl"`
}

// ServerService allows retrieving information about the server
type ServerService struct {
	sling *sling.Sling
}

func newServerService(base *sling.Sling) *ServerService {
	return &ServerService{
		sling: base.Get("server/"),
	}
}

// Get returns a struct with server information
func (s *ServerService) Get() (*Server, error) {

	var out Server

	_, err := s.sling.ReceiveSuccess(&out)

	if err != nil {
		return nil, err
	}

	return &out, nil
}

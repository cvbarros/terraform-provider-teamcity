package teamcity

import "fmt"

// StatusPublisherGithubOptions represents parameters used to create Github Commit Status Publisher Feature
type StatusPublisherGithubOptions struct {
	//AuthenticationType can be 'password' or 'token'
	AuthenticationType string
	//Host is the Github URL, for instance "https://api.github.com" or "https://hostname/api/v3/" for Github Enterprise
	Host string
	//Username is required if AuthenticationType is 'password'
	Username string
	//Password is required if AuthenticationType is 'password'
	Password string
	//AccessToken is required if AuthenticationType is 'token'
	AccessToken string
}

//NewCommitStatusPublisherGithubOptionsPassword returns options created for AuthenticationType = 'password'. No validation is performed, parameters indicate mandatory fields.
func NewCommitStatusPublisherGithubOptionsPassword(host string, username string, password string) StatusPublisherGithubOptions {
	return StatusPublisherGithubOptions{
		Host:               host,
		AuthenticationType: "password",
		Username:           username,
		Password:           password,
	}
}

//NewCommitStatusPublisherGithubOptionsToken returns options created for AuthenticationType = 'token'. No validation is performed, parameters indicate mandatory fields.
func NewCommitStatusPublisherGithubOptionsToken(host string, accessToken string) StatusPublisherGithubOptions {
	return StatusPublisherGithubOptions{
		Host:               host,
		AuthenticationType: "token",
		AccessToken:        accessToken,
	}
}

//NewFeatureCommitStatusPublisherGithub creates a Build Feature Commit status Publisher to Github with the given options and validates the required properties
func NewFeatureCommitStatusPublisherGithub(opt StatusPublisherGithubOptions) (*FeatureCommitStatusPublisher, error) {
	if opt.AuthenticationType == "" {
		return nil, fmt.Errorf("AuthenticationType is required")
	}

	if opt.AuthenticationType != "password" && opt.AuthenticationType != "token" {
		return nil, fmt.Errorf("invalid AuthenticationType, must be 'password' or 'token'")
	}

	if opt.Host == "" {
		return nil, fmt.Errorf("Host is required")
	}

	if opt.AuthenticationType == "password" {
		if opt.Username == "" || opt.Password == "" {
			return nil, fmt.Errorf("username/password required for auth type 'password'")
		}
	}

	if opt.AuthenticationType == "token" {
		if opt.AccessToken == "" {
			return nil, fmt.Errorf("accesstoken required for auth type 'token'")
		}
	}

	out := &FeatureCommitStatusPublisher{
		Options:    opt,
		properties: opt.Properties(),
	}

	return out, nil
}

//Properties returns a *Properties collection with properties filled related to this commit publisher parameters to be used in build features
func (s StatusPublisherGithubOptions) Properties() *Properties {
	props := NewPropertiesEmpty()

	props.AddOrReplaceValue("publisherId", "githubStatusPublisher")
	props.AddOrReplaceValue("github_authentication_type", s.AuthenticationType)
	props.AddOrReplaceValue("github_host", s.Host)

	if s.AuthenticationType == "password" {
		props.AddOrReplaceValue("github_username", s.Username)
		props.AddOrReplaceValue("secure:github_password", s.Password)
	}

	if s.AuthenticationType == "token" {
		props.AddOrReplaceValue("secure:github_access_token", s.AccessToken)
	}

	return props
}

//CommitStatusPublisherGithubOptionsFromProperties grabs a Properties collection and transforms back to a StatusPublisherGithubOptions
func CommitStatusPublisherGithubOptionsFromProperties(p *Properties) (*StatusPublisherGithubOptions, error) {
	var out StatusPublisherGithubOptions
	if host, ok := p.GetOk("github_host"); ok {
		out.Host = host
	} else {
		return nil, fmt.Errorf("Properties do not have 'github_host' key")
	}

	if authType, ok := p.GetOk("github_authentication_type"); ok {
		out.AuthenticationType = authType
		switch authType {
		case "password":
			u, _ := p.GetOk("github_username")
			out.Username = u

			//Password or AccessToken is never returned from properties, because it is secure. Once set, we cannot read it back
		}
	} else {
		return nil, fmt.Errorf("Properties do not have 'github_authentication_type' key")
	}

	return &out, nil
}

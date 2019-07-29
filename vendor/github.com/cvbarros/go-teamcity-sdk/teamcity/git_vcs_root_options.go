package teamcity

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

//GitAuthMethod enum is to specify the authentication method when connecting to Git VCS.
type GitAuthMethod string

const (
	//GitAuthMethodAnonymous is used for anonymously connecting to Git VCS.
	GitAuthMethodAnonymous GitAuthMethod = "ANONYMOUS"

	//GitAuthMethodPassword is used for connecting using username/password to Git VCS.
	GitAuthMethodPassword GitAuthMethod = "PASSWORD"

	//GitAuthSSHUploadedKey is used for connecting using SSH with a Private Key uploaded to TeamCity
	GitAuthSSHUploadedKey GitAuthMethod = "TEAMCITY_SSH_KEY"

	//GitAuthSSHDefaultKey is used for connecting using SSH and uses mapping specified in the file /root/.ssh/config if that file exists.
	GitAuthSSHDefaultKey GitAuthMethod = "PRIVATE_KEY_DEFAULT"

	//GitAuthSSHCustomKey is used for connecting using SSH with a custom private key in the given path.
	GitAuthSSHCustomKey GitAuthMethod = "PRIVATE_KEY_FILE"
)

//GitAgentCleanPolicy enum specifies when the "git clean" should be run on the agent.
type GitAgentCleanPolicy string

const (
	//CleanPolicyBranchChange run clean whenever a branch change is detected
	CleanPolicyBranchChange GitAgentCleanPolicy = "ON_BRANCH_CHANGE"

	//CleanPolicyAlways always run 'git clean'
	CleanPolicyAlways GitAgentCleanPolicy = "ALWAYS"

	//CleanPolicyNever never run 'git clean'
	CleanPolicyNever GitAgentCleanPolicy = "NEVER"
)

//GitAgentCleanFilesPolicy enum specifies which files will be removed when "git clean" command is run on agent.
type GitAgentCleanFilesPolicy string

const (

	//CleanFilesPolicyAllUntracked will clean all untracked files
	CleanFilesPolicyAllUntracked GitAgentCleanFilesPolicy = "ALL_UNTRACKED"

	//CleanFilesPolicyIgnoredOnly will clean all ignored files
	CleanFilesPolicyIgnoredOnly GitAgentCleanFilesPolicy = "IGNORED_ONLY"

	//CleanFilesPolicyIgnoredUntracked will clean all non-ignored untracked files
	CleanFilesPolicyIgnoredUntracked GitAgentCleanFilesPolicy = "NON_IGNORED_ONLY"
)

// GitAgentSettings are agent-specific settings that are used in case of agent checkout
type GitAgentSettings struct {

	//GitPath is the path to a git executable on the agent. If blank, the location set up in TEAMCITY_GIT_PATH environment variable is used.
	GitPath string `prop:"agentGitPath"`

	//CleanPolicy specifies when the "git clean" command is run on the agent. Defaults to 'CleanPolicyBranchChange'
	CleanPolicy GitAgentCleanPolicy `prop:"agentCleanPolicy"`

	//CleanFilesPolicy specifies which files will be removed when "git clean" command is run on agent.
	CleanFilesPolicy GitAgentCleanFilesPolicy `prop:"agentCleanFilesPolicy"`

	//UseMirrors when enabled, TeamCity creates a separate clone of the repository on each agent and uses it in the checkout directory via git alternates.
	UseMirrors bool `prop:"useAlternates"`
}

// GitVcsUsernameStyle defines a way TeamCity binds VCS changes to the user.
// With selected style and the following content of ~/.gitconfig:
// [user]
// name = Joe Coder
// email = joe.coder@acme.com
type GitVcsUsernameStyle string

const (
	//GitVcsUsernameStyleUserID should configure 'joe.coder' in the VCS profile for the User.
	GitVcsUsernameStyleUserID GitVcsUsernameStyle = "USERID"
	//GitVcsUsernameStyleAuthorName should configure 'Joe Coder' in the VCS profile for the User.
	GitVcsUsernameStyleAuthorName GitVcsUsernameStyle = "NAME"
	//GitVcsUsernameStyleAuthorEmail should configure 'joe.coder@acme.com' in the VCS profile for the User.
	GitVcsUsernameStyleAuthorEmail GitVcsUsernameStyle = "EMAIL"
	//GitVcsUsernameStyleAuthorNameAndEmail should configure 'Joe Coder <joe.coder@acme.com>' in the VCS profile for the User.
	GitVcsUsernameStyleAuthorNameAndEmail GitVcsUsernameStyle = "FULL"
)

// GitVcsRootOptions represents parameters used when manipulating VCS Roots of type "Git"
type GitVcsRootOptions struct {
	//DefaultBranch indicates which main branch or tag to be monitored by the VCS Root. Requied.
	DefaultBranch string `prop:"branch"`

	//BrancSpec are monitor besides the default one as a newline-delimited set of rules in the form of +|-:branch name (with the optional * placeholder)
	//Set separately, outside constructor.
	BranchSpec []string `prop:"teamcity:branchSpec" separator:"\\n"`

	//FetchURL is used for fetching data from the repository. Required.
	FetchURL string `prop:"url"`

	//PushURL is used for pushing tags to the remote repository. If blank, the fetch url is used.
	PushURL string `prop:"push_url"`

	//EnableTagsInBranchSpec enable/disable use tags in branch specification. Defaults to false. Set separately, outside constructor.
	EnableTagsInBranchSpec bool `prop:"reportTagRevisions"`

	//AuthMethod controls how the TeamCity server will authenticate against the VCS Git provider. Required.
	AuthMethod GitAuthMethod `prop:"authMethod"`

	//Username is used for methods other than "GitAuthMethodAnonymous". For SSH, it overrides the username used in the Fetch/Push URLs.
	Username string `prop:"username"`

	//Password represents the user password when "GitAuthMethodPassword" auth method is used.
	//Password is the key passphrase when "GitAuthSSHUploadedKey" or "GitAuthSSHCustomKey" with auth methods are used.
	Password string `prop:"secure:password"`

	//PrivateKeySource is used for "GitAuthSSHCustomKey" as the key path on disk.
	//PrivateKeySource is used for "GitAuthSSHUploadedKey" as the name of the SSH Key uploaded for the project in "SSH Keys"
	PrivateKeySource string `prop:"secure:passphrase"`

	//AgentSettings control agent-specific settings that are used in case of agent checkout
	AgentSettings *GitAgentSettings

	//SubModuleCheckout specifies whether checkout Git submodules or ignore.
	//Possible values are 'CHECKOUT' or 'IGNORE'. Defaults to 'CHECKOUT'. Set separately, outside constructor.
	SubModuleCheckout string `prop:"submoduleCheckout"`

	// UsernameStyle defines a way TeamCity binds VCS changes to the user.
	// Defaults to 'GitVcsUsernameStyleUserID'. Set separately, outside constructor.
	UsernameStyle GitVcsUsernameStyle `prop:"usernameStyle"`
}

//NewGitVcsRootOptions returns a new instance of GitVcsRootOptions with default GitAgentSettings
func NewGitVcsRootOptions(defaultBranch string, fetchURL string, pushURL string, auth GitAuthMethod, username string, password string) (*GitVcsRootOptions, error) {
	return NewGitVcsRootOptionsWithAgentSettings(defaultBranch, fetchURL, pushURL, auth, username, password, nil)
}

//NewGitVcsRootOptionsDefaults returns a new instance of GitVcsRootOptions with default values
//Anonymous auth method
//Default AgentSettings
func NewGitVcsRootOptionsDefaults(defaultBranch string, fetchURL string) (*GitVcsRootOptions, error) {
	return NewGitVcsRootOptions(defaultBranch, fetchURL, "", GitAuthMethodAnonymous, "", "")
}

//NewGitVcsRootOptionsWithAgentSettings returns a new instance of GitVcsRootOptions with specified GitAgentSettings
func NewGitVcsRootOptionsWithAgentSettings(defaultBranch string, fetchURL string, pushURL string, auth GitAuthMethod, username string, password string, agentSettings *GitAgentSettings) (*GitVcsRootOptions, error) {
	if auth == "" {
		return nil, errors.New("auth is required")
	}
	if defaultBranch == "" {
		return nil, errors.New("defaultBranch is required")
	}
	if fetchURL == "" {
		return nil, errors.New("fetchURL is required")
	}

	if auth == GitAuthMethodPassword {
		if username == "" {
			return nil, fmt.Errorf("username is required if using auth method '%s'", auth)
		}
	}
	if agentSettings == nil {
		agentSettings = &GitAgentSettings{
			UseMirrors:       true,
			CleanPolicy:      CleanPolicyBranchChange,
			CleanFilesPolicy: CleanFilesPolicyAllUntracked,
		}
	}

	opt := &GitVcsRootOptions{
		DefaultBranch:          defaultBranch,
		FetchURL:               fetchURL,
		PushURL:                pushURL,
		AuthMethod:             auth,
		Username:               username,
		Password:               password,
		EnableTagsInBranchSpec: false,
		UsernameStyle:          GitVcsUsernameStyleUserID,
		SubModuleCheckout:      "CHECKOUT",
		AgentSettings:          agentSettings,
	}

	if opt.PushURL == "" {
		opt.PushURL = opt.FetchURL
	}

	return opt, nil
}
func (o *GitVcsRootOptions) properties() *Properties {
	p := NewPropertiesEmpty()

	p.AddOrReplaceValue("branch", o.DefaultBranch)
	p.AddOrReplaceValue("authMethod", string(o.AuthMethod))
	p.AddOrReplaceValue("url", o.FetchURL)
	p.AddOrReplaceValue("push_url", o.PushURL)
	p.AddOrReplaceValue("usernameStyle", string(o.UsernameStyle))
	p.AddOrReplaceValue("submoduleCheckout", o.SubModuleCheckout)

	p.AddOrReplaceValue("ignoreKnownHosts", "true") // This is always true. Couldn't find an option on the UI to change this setting

	switch o.AuthMethod {
	case GitAuthMethodPassword:
		p.AddOrReplaceValue("username", o.Username)
		p.AddOrReplaceValue("secure:password", o.Password)
	case GitAuthSSHUploadedKey:
		p.AddOrReplaceValue("username", o.Username)
		p.AddOrReplaceValue("secure:passphrase", o.Password)
		p.AddOrReplaceValue("teamcitySshKey", o.PrivateKeySource)
	case GitAuthSSHCustomKey:
		p.AddOrReplaceValue("username", o.Username)
		p.AddOrReplaceValue("secure:passphrase", o.Password)
		p.AddOrReplaceValue("privateKeyPath", o.PrivateKeySource)
	case GitAuthSSHDefaultKey:
		p.AddOrReplaceValue("username", o.Username)
	}

	if len(o.BranchSpec) > 0 {
		// Some properties use \\r\\n to split. But this one only uses \\n, conversely
		p.AddOrReplaceValue("teamcity:branchSpec", strings.Join(o.BranchSpec, "\\n"))
	}

	if o.EnableTagsInBranchSpec {
		p.AddOrReplaceValue("reportTagRevisions", "true")
	}

	agentP := o.AgentSettings.properties()
	for _, ap := range agentP.Items {
		p.AddOrReplaceProperty(ap)
	}

	return p
}

func (s *GitAgentSettings) properties() *Properties {
	p := NewPropertiesEmpty()

	if s.GitPath != "" {
		p.AddOrReplaceValue("agentGitPath", s.GitPath)
	}

	p.AddOrReplaceValue("agentCleanPolicy", string(s.CleanPolicy))
	p.AddOrReplaceValue("agentCleanFilesPolicy", string(s.CleanFilesPolicy))
	p.AddOrReplaceValue("useAlternates", strconv.FormatBool(s.UseMirrors))

	return p
}

func (p *Properties) gitVcsOptions() *GitVcsRootOptions {
	var out GitVcsRootOptions
	var agt GitAgentSettings

	fillStructFromProperties(&out, p)
	fillStructFromProperties(&agt, p)

	out.AgentSettings = &agt

	return &out
}

func (p *Properties) gitAgentSettings() *GitAgentSettings {
	var out GitAgentSettings
	fillStructFromProperties(&out, p)
	return &out
}

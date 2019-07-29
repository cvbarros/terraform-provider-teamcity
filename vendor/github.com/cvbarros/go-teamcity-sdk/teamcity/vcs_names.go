package teamcity

type vcsName = string

const (
	//Git vcs type
	Git vcsName = "jetbrains.git"
)

// VcsNames represents possible vcsNames for VCS Roots
var VcsNames = struct {
	Git vcsName
}{
	Git: Git,
}

package teamcity

import (
	"errors"
	"strconv"
	"strings"
)

//ArtifactDependencyRevision has some allowed values listed per constants.
type ArtifactDependencyRevision string

const (
	// LatestSuccessfulBuild grabs the artifacts produced by the last sucessful build for the source build configuration.
	LatestSuccessfulBuild ArtifactDependencyRevision = "lastSuccessful"

	// LatestPinnedBuild grabs the artifacts produced by the last pinned build for the source build configuration.
	LatestPinnedBuild ArtifactDependencyRevision = "lastPinned"

	// LatestFinishedBuild grabs the artifacts produced by the last finished build, sucessful or not, for the source build configuration.
	LatestFinishedBuild ArtifactDependencyRevision = "lastFinished"

	// BuildFromSameChain grabs the artifacts produced by the source build triggered within the same build chain.
	BuildFromSameChain ArtifactDependencyRevision = "sameChainOrLastFinished"

	// BuildWithSpecifiedNumber grabs the artifacts produced by the source build that has a specific build number.
	BuildWithSpecifiedNumber ArtifactDependencyRevision = "buildNumber"

	// LastBuildFinishedWithTag grabs the artifacts produced by the source build that have a specific VCS tag
	LastBuildFinishedWithTag ArtifactDependencyRevision = "buildTag"
)

//ArtifactDependencyOptions represents options when creating an artifact dependency for a build configuration.
//For more information see: https://confluence.jetbrains.com/display/TCD10/Artifact+Dependencies
type ArtifactDependencyOptions struct {
	//ArtifactRevisionType maps to the TeamCity UI's "Get artifacts from", indicating which build should be used as artifact source.
	ArtifactRevisionType ArtifactDependencyRevision `prop:"revisionName"`

	//PathRules is a list of rules to match files that will have to be dowloaded from the source build that output artifacts.
	PathRules []string `prop:"pathRules"`

	//CleanDestination
	CleanDestination bool `prop:"cleanDestinationDirectory"`

	//RevisionNumber is used in conjunction with `BuildWithSpecifiedNumber` (as the build number value) or `LastBuildFinishedWithTag` (as the tag value to look for)
	RevisionNumber string
}

//NewArtifactDependencyOptions creates an instance of ArtifactDependencyOptions with default values.
//
//(required) pathRules - list of rules to match files that will have to be dowloaded from the source build that output artifacts.
//They can be specified in the format [+:|-:]SourcePath[!ArchivePath][=>DestinationPath]
//
//(required) revisionType - Which kind of revision type the artifact dependency will be based upon. See ArtifactDependencyRevision enum for options.
//
//(optional) revisionValue - Required if using `BuildWithSpecifiedNumber` or `LastBuildFinishedWithTag`
func NewArtifactDependencyOptions(pathRules []string, revisionType ArtifactDependencyRevision, cleanDestination bool, revisionValue string) (*ArtifactDependencyOptions, error) {
	if len(pathRules) == 0 {
		return nil, errors.New("pathRules is required")
	}

	if revisionType == "" {
		return nil, errors.New("revisionType is required")
	}

	if revisionType == BuildWithSpecifiedNumber && revisionValue == "" {
		return nil, errors.New("revisionValue is required is using 'BuildWithSpecifiedNumber'")
	}

	if revisionType == LastBuildFinishedWithTag && revisionValue == "" {
		return nil, errors.New("revisionValue is required is using 'LastBuildFinishedWithTag'")
	}

	return &ArtifactDependencyOptions{
		PathRules:            pathRules,
		CleanDestination:     cleanDestination,
		ArtifactRevisionType: revisionType,
		RevisionNumber:       revisionValue,
	}, nil
}

func (o *ArtifactDependencyOptions) properties() *Properties {
	p := NewPropertiesEmpty()

	p.AddOrReplaceValue("pathRules", strings.Join(o.PathRules, "\r\n"))
	p.AddOrReplaceValue("cleanDestinationDirectory", strconv.FormatBool(o.CleanDestination))

	switch o.ArtifactRevisionType {
	case BuildWithSpecifiedNumber:
		p.AddOrReplaceValue("revisionValue", o.RevisionNumber)
		break
	case LastBuildFinishedWithTag:
		p.AddOrReplaceValue("revisionValue", o.RevisionNumber+".tcbuildtag")
	default:
		p.AddOrReplaceValue("revisionValue", "latest."+string(o.ArtifactRevisionType))
	}

	p.AddOrReplaceValue("revisionName", string(o.ArtifactRevisionType))
	return p
}

func (p *Properties) artifactDependencyOptions() *ArtifactDependencyOptions {
	var out ArtifactDependencyOptions

	fillStructFromProperties(&out, p)

	if v, ok := p.GetOk("revisionValue"); ok {
		s := strings.TrimSuffix(v, ".tcbuildtag")
		out.RevisionNumber = s
	}

	return &out
}

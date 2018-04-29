package teamcity

// Step is a build configuration/template build step. Use constructor functions NewStep* to create those
type Step struct {

	// disabled
	Disabled *bool `json:"disabled,omitempty" xml:"disabled"`

	// href
	Href string `json:"href,omitempty" xml:"href"`

	// id
	ID string `json:"id,omitempty" xml:"id"`

	// inherited
	Inherited *bool `json:"inherited,omitempty" xml:"inherited"`

	// name
	Name string `json:"name,omitempty" xml:"name"`

	// properties
	Properties *Properties `json:"properties,omitempty"`

	// type
	Type string `json:"type,omitempty" xml:"type"`
}

// Steps represents a collection of Steps
type Steps struct {

	// count
	Count int32 `json:"count,omitempty" xml:"count"`

	// step
	Items []*Step `json:"step"`
}

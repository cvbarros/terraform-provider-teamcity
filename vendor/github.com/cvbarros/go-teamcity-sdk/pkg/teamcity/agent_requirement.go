package teamcity

import "errors"

// AgentRequirement is a condition evaluated per agent to see if a build type is compatible or not
type AgentRequirement struct {

	// id
	ID string `json:"id,omitempty" xml:"id"`

	// inherited
	Inherited *bool `json:"inherited,omitempty" xml:"inherited"`

	// inherited
	Disabled *bool `json:"disabled,omitempty" xml:"disabled"`

	// type
	Condition string `json:"type,omitempty"`

	// Do not use this directly, build this struct via NewAgentRequirement
	Properties *Properties `json:"properties,omitempty"`
}

// AgentRequirements is a collection of AgentRequirement
type AgentRequirements struct {

	// count
	Count int32 `json:"count,omitempty" xml:"count"`

	// href
	Href string `json:"href,omitempty" xml:"href"`

	// property
	Items []*AgentRequirement `json:"agent-requirement"`
}

//Conditions - Possible conditions for requirements. Do not change the values.
var Conditions = struct {
	Exists            string
	Equals            string
	DoesNotEqual      string
	MoreThan          string
	NoMoreThan        string
	LessThan          string
	NoLessThan        string
	StartsWith        string
	Contains          string
	DoesNotContain    string
	EndsWith          string
	Matches           string
	DoesNotMatch      string
	VersionMoreThan   string
	VersionNoMoreThan string
	VersionLessThan   string
	VersionNoLessThan string
}{
	Exists:            "exists",
	Equals:            "equals",
	DoesNotEqual:      "does-not-equal",
	MoreThan:          "more-than",
	NoMoreThan:        "no-more-than",
	LessThan:          "less-than",
	NoLessThan:        "no-less-than",
	StartsWith:        "starts-with",
	Contains:          "contains",
	DoesNotContain:    "does-not-contain",
	EndsWith:          "ends-with",
	Matches:           "matches",
	DoesNotMatch:      "does-not-match",
	VersionMoreThan:   "ver-more-than",
	VersionNoMoreThan: "ver-no-more-than",
	VersionLessThan:   "ver-less-than",
	VersionNoLessThan: "ver-no-less-than",
}

// NewAgentRequirement creates AgentRequirement structure with correct representation. Use this instead of creating the struct manually.
func NewAgentRequirement(condition string, paramName string, paramValue string) (*AgentRequirement, error) {

	// Sample structure for a requirement
	// The "property-name" and "property-value" properties nested are used as operands for the condition
	// {
	// 	"id": "RQ_17",
	// 	"type": "ver-no-more-than",
	// 	"properties": {
	// 		"count": 2,
	// 		"property": [
	// 			{
	// 				"name": "property-name",
	// 				"value": "r"
	// 			},
	// 			{
	// 				"name": "property-value",
	// 				"value": "a"
	// 			}
	// 		]
	// 	}
	// },

	if condition != Conditions.Exists && paramValue == "" {
		return nil, errors.New("paramValue is required except for 'exists' condition")
	}

	propertyNameProp := &Property{Name: "property-name", Value: paramName}
	props := NewProperties(propertyNameProp)

	// 'exists' uses only "property-name" operand
	if condition != Conditions.Exists {
		propertyValueProp := &Property{Name: "property-value", Value: paramValue}
		props.Add(propertyValueProp)
	}

	return &AgentRequirement{
		Condition:  condition,
		Properties: props,
	}, nil
}

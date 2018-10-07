package teamcity

import (
	"encoding/json"
	"fmt"
	"strings"
)

type paramType = string

const (
	configParamType = "configuration"
	systemParamType = "system"
	envVarParamType = "env"
)

// ParameterTypes represent the possible parameter types
var ParameterTypes = struct {
	Configuration       paramType
	System              paramType
	EnvironmentVariable paramType
}{
	Configuration:       configParamType,
	System:              systemParamType,
	EnvironmentVariable: envVarParamType,
}

//Parameters is a strongly-typed collection of "Parameter" suitable for serialization
type Parameters struct {
	Count int32        `json:"count,omitempty" xml:"count"`
	Href  string       `json:"href,omitempty" xml:"href"`
	Items []*Parameter `json:"property"`
}

//Parameter represents a project or build configuration parameter that may be defined as "configuration", "system" or "environment variable"
type Parameter struct {
	Inherited *bool `json:"inherited,omitempty" xml:"inherited"`

	Name string `json:"name,omitempty" xml:"name"`

	Value string `json:"value" xml:"value"`

	Type string `json:"-"`
}

//NewParametersEmpty returns an empty collection of Parameters
func NewParametersEmpty() *Parameters {
	return &Parameters{
		Count: 0,
		Items: make([]*Parameter, 0),
	}
}

// NewParameters returns an instance of Parameters collection with the given parameters slice
func NewParameters(items ...*Parameter) *Parameters {
	count := len(items)
	return &Parameters{
		Count: int32(count),
		Items: items,
	}
}

//NewParameter creates a new instance of a parameter with the given type
func NewParameter(t string, name string, value string) (*Parameter, error) {
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if t != ParameterTypes.Configuration && t != ParameterTypes.EnvironmentVariable && t != ParameterTypes.System {
		return nil, fmt.Errorf("invalid parameter type, use one of the values defined in ParameterTypes")
	}

	return &Parameter{
		Type:  string(t),
		Name:  name,
		Value: value,
	}, nil
}

//MarshalJSON implements JSON serialization for Parameter
func (p *Parameter) MarshalJSON() ([]byte, error) {
	out := p.Property()

	return json.Marshal(out)
}

//UnmarshalJSON implements JSON deserialization for Parameter
func (p *Parameter) UnmarshalJSON(data []byte) error {
	var aux Property
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	var name, paramType string
	if strings.HasPrefix(aux.Name, "system.") {
		name = strings.TrimPrefix(aux.Name, "system.")
		paramType = ParameterTypes.System
	} else if strings.HasPrefix(aux.Name, "env.") {
		name = strings.TrimPrefix(aux.Name, "env.")
		paramType = ParameterTypes.EnvironmentVariable
	} else {
		name = aux.Name
		paramType = ParameterTypes.Configuration
	}
	p.Name = name
	p.Inherited = aux.Inherited
	p.Value = aux.Value
	p.Type = paramType
	return nil
}

//Properties convert a Parameters collection to a Properties collection
func (p *Parameters) Properties() *Properties {
	out := NewPropertiesEmpty()
	for _, i := range p.Items {
		out.AddOrReplaceProperty(i.Property())
	}
	return out
}

//Property converts a Parameter instance to a Property
func (p *Parameter) Property() *Property {
	return &Property{
		Name:      fmt.Sprintf("%s%s", paramPrefixByType[p.Type], p.Name),
		Value:     p.Value,
		Inherited: p.Inherited,
	}
}

// AddOrReplaceValue will update a parameter value if it exists, or add if it doesnt
func (p *Parameters) AddOrReplaceValue(t string, n string, v string) {
	for _, elem := range p.Items {
		if elem == nil {
			continue
		}

		if elem.Name == n {
			elem.Value = v
			return
		}
	}
	param, _ := NewParameter(t, n, v)
	p.Add(param)
}

// AddOrReplaceParameter will update a parameter value if another parameter with the same name exists. It won't replace the Parameter struct within the Parameters collection.
func (p *Parameters) AddOrReplaceParameter(param *Parameter) {
	p.AddOrReplaceValue(param.Type, param.Name, param.Value)
}

// Add a new parameter to this collection
func (p *Parameters) Add(param *Parameter) {
	p.Count++
	p.Items = append(p.Items, param)
}

// Concat appends the source Parameters collection to this collection and returns the appended collection
func (p *Parameters) Concat(source *Parameters) *Parameters {
	for _, item := range source.Items {
		p.AddOrReplaceParameter(item)
	}
	return p
}

//Remove a parameter if it exists in the collection
func (p *Parameters) Remove(t string, n string) {
	removed := -1
	for i := range p.Items {
		if p.Items[i].Name == n && p.Items[i].Type == t {
			removed = i
			break
		}
	}
	if removed >= 0 {
		p.Count--
		p.Items = append(p.Items[:removed], p.Items[removed+1:]...)
	}
}

var paramPrefixByType = map[string]string{
	string(ParameterTypes.Configuration):       "",
	string(ParameterTypes.System):              "system.",
	string(ParameterTypes.EnvironmentVariable): "env.",
}

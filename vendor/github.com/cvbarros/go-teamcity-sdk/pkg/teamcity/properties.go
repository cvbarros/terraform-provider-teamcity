package teamcity

// Property represents a key/value/type structure used by several resources to extend their representation
type Property struct {

	// inherited
	Inherited *bool `json:"inherited,omitempty" xml:"inherited"`

	// name
	Name string `json:"name,omitempty" xml:"name"`

	// type
	Type *Type `json:"type,omitempty"`

	// value
	Value string `json:"value,omitempty" xml:"value"`
}

// Type represents a parameter type . The rawValue is the parameter specification as defined in the UI.
type Type struct {
	// raw value
	RawValue string `json:"rawValue,omitempty" xml:"rawValue"`
}

// Properties represents a collection of key/value properties for a resource
type Properties struct {

	// count
	Count int32 `json:"count,omitempty" xml:"count"`

	// href
	Href string `json:"href,omitempty" xml:"href"`

	// property
	Items []*Property `json:"property"`
}

// NewProperties returns an instance of Properties collection
func NewProperties(items ...*Property) *Properties {
	count := len(items)
	return &Properties{
		Count: int32(count),
		Items: items,
	}
}

// NewProperty returns an instance of Property
func NewProperty(name string, value string) *Property {
	return &Property{
		Name:  name,
		Value: value,
	}
}

// Add a new property to this collection
func (p *Properties) Add(prop *Property) {
	p.Count++
	p.Items = append(p.Items, prop)
}

// AddOrReplace will update a property value if it exists, or add if it doesnt
func (p *Properties) AddOrReplaceValue(prop *Property) {
	for _, elem := range p.Items {
		if elem == nil {
			continue
		}

		if elem.Name == prop.Name {
			elem.Value = prop.Value
			return
		}
	}

	p.Add(prop)
}

// Map converts Properties to a key/value dictionary as map[string]string
func (p *Properties) Map() map[string]string {
	out := make(map[string]string)
	for _, item := range p.Items {
		out[item.Name] = item.Value
	}

	return out
}

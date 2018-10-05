package teamcity

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Property represents a key/value/type structure used by several resources to extend their representation
type Property struct {

	// inherited
	Inherited *bool `json:"inherited,omitempty" xml:"inherited"`

	// name
	Name string `json:"name,omitempty" xml:"name"`

	// type
	Type *Type `json:"type,omitempty"`

	// value
	Value string `json:"value" xml:"value"`
}

func (p *Property) String() string {
	return fmt.Sprintf("Name: '%s', Value: '%s'", p.Name, p.Value)
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

// NewPropertiesEmpty returns an instance of Properties collection with no properties
func NewPropertiesEmpty() *Properties {
	return &Properties{
		Count: 0,
		Items: make([]*Property, 0),
	}
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

//Remove a property if it exists in the collection
func (p *Properties) Remove(n string) {
	removed := -1
	for i := range p.Items {
		if p.Items[i].Name == n {
			removed = i
			break
		}
	}
	if removed >= 0 {
		p.Count--
		p.Items = append(p.Items[:removed], p.Items[removed+1:]...)
	}
}

// AddOrReplaceValue will update a property value if it exists, or add if it doesnt
func (p *Properties) AddOrReplaceValue(n string, v string) {
	for _, elem := range p.Items {
		if elem == nil {
			continue
		}

		if elem.Name == n {
			elem.Value = v
			return
		}
	}

	p.Add(&Property{Name: n, Value: v})
}

// AddOrReplaceProperty will update a property value if another property with the same name exists. It won't replace the Property struct within the Properties collection.
func (p *Properties) AddOrReplaceProperty(prop *Property) {
	p.AddOrReplaceValue(prop.Name, prop.Value)
}

// Concat appends the source Properties collection to this collection and returns the appended collection
func (p *Properties) Concat(source *Properties) *Properties {
	for _, item := range source.Items {
		p.AddOrReplaceProperty(item)
	}
	return p
}

// GetOk returns the value of the propery and true if found, otherwise ""/false
func (p *Properties) GetOk(key string) (string, bool) {
	if len(p.Items) == 0 {
		return "", false
	}

	for _, v := range p.Items {
		if v.Name == key {
			return v.Value, true
		}
	}

	return "", false
}

// Map converts Properties to a key/value dictionary as map[string]string
func (p *Properties) Map() map[string]string {
	out := make(map[string]string)
	for _, item := range p.Items {
		out[item.Name] = item.Value
	}

	return out
}

func fillStructFromProperties(data interface{}, p *Properties) {
	t := reflect.TypeOf(data).Elem()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if v, ok := f.Tag.Lookup("prop"); ok {
			sf := reflect.ValueOf(data).Elem().Field(i)
			if pv, pok := p.GetOk(v); pok {
				switch sf.Kind() {
				case reflect.Uint:
					bv, _ := strconv.ParseUint(pv, 10, 0)
					sf.SetUint(bv)
				case reflect.Int:
					bv, _ := strconv.ParseInt(pv, 10, 0)
					sf.SetInt(bv)
				case reflect.Bool:
					bv, _ := strconv.ParseBool(pv)
					sf.SetBool(bv)
				case reflect.String:
					sf.SetString(pv)
				case reflect.Slice:
					var sep string
					sep, ok = f.Tag.Lookup("separator")
					if !ok {
						sep = "\\r\\n" // Use default
					}
					sVal := reflect.ValueOf(strings.Split(pv, sep))
					sf.Set(sVal)
				default:
					//TODO: Panic if cannot set value
					continue
				}
			}
		}
	}
}

func serializeToProperties(data interface{}) *Properties {
	props := NewPropertiesEmpty()
	t := reflect.TypeOf(data).Elem()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if v, ok := f.Tag.Lookup("prop"); ok {
			pv := reflect.ValueOf(data).Elem().Field(i)
			switch pv.Kind() {
			case reflect.Slice:
				sep := "\\r\\n" // Use default
				sep, _ = f.Tag.Lookup("separator")
				pVal := strings.Join(pv.Interface().([]string), sep)
				props.AddOrReplaceValue(v, pVal)
			case reflect.Bool:
				pVal := pv.Bool()
				_, force := f.Tag.Lookup("force")
				if pVal || force {
					props.AddOrReplaceValue(v, strconv.FormatBool(pVal))
				}
			case reflect.Int:
				props.AddOrReplaceValue(v, fmt.Sprint(pv.Int()))
			case reflect.Uint:
				props.AddOrReplaceValue(v, fmt.Sprint(pv.Uint()))
			default:
				pVal := pv.String()
				if pVal != "" {
					props.AddOrReplaceValue(v, pVal)
				}
			}
		}
	}
	return props
}

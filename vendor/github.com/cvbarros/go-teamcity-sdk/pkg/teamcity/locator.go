package teamcity

import "net/url"

//Locator represents a arbitraty locator to be used when querying resources, such as id: or type:
//These are used in GET requests within the URL so must be properly escaped
type Locator string

//LocatorID creates a locator for a Project/BuildType by Id
func LocatorID(id string) Locator {
	return Locator(url.QueryEscape("id:") + id)
}

//LocatorName creates a locator for Project/BuildType by Name
func LocatorName(name string) Locator {
	return Locator(url.QueryEscape("name:") + url.PathEscape(name))
}

func (l Locator) String() string {
	return string(l)
}

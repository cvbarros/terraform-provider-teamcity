package teamcity

import "net/url"

//Locator represents a arbitraty locator to be used when querying resources, such as id: or type:
//These are used in GET requests within the URL so must be properly escaped
type Locator string

//LocatorId creates a locator for a Project by Id
func LocatorId(id string) Locator {
	return Locator("id:" + id)
}

func (l Locator) String() string {
	return url.QueryEscape(string(l))
}

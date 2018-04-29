package teamcity

// VcsRootEntries represents a collection of VCS Roots attached to a resource
type VcsRootEntries struct {
	// count
	Count int32 `json:"count,omitempty" xml:"count"`

	// href
	Href string `json:"href,omitempty" xml:"href"`

	// property
	Items []*VcsRootEntry `json:"vcs-root-entry"`
}

// VcsRootEntry represents a single VCS Root attached to a resource
type VcsRootEntry struct {
	// id
	Id string `json:"id,omitempty" xmld:"id"`

	// inherited
	Inherited *bool `json:"inherited,omitempty" xml:"inherited"`

	// checkout rules
	CheckoutRules string `json:"checkout-rules,omitempty"`

	// vcs root
	VcsRoot *VcsRootReference `json:"vcs-root,omitempty"`
}

// NewVcsRootEntries returns an instance of VcsRootEntries collection
func NewVcsRootEntries(items ...*VcsRootReference) *VcsRootEntries {
	count := len(items)
	entries := make([]*VcsRootEntry, count)
	for i, item := range items {
		entries[i] = NewVcsRootEntry(item)
	}

	return &VcsRootEntries{
		Count: int32(count),
		Items: entries,
	}
}

// NewVcsRootEntry is a convenience function for attaching a VcsRootReference to a Build configuration, represented as a VcsRootEntry in Teamcity API
func NewVcsRootEntry(vcsRef *VcsRootReference) *VcsRootEntry {
	return &VcsRootEntry{
		VcsRoot: vcsRef,
	}
}

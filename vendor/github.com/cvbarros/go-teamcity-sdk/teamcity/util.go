package teamcity

// NewTrue is a helper function to return a *bool to true
func NewTrue() *bool {
	return NewBool(true)
}

// NewFalse is a helper function to return a *bool to true
func NewFalse() *bool {
	return NewBool(false)
}

// NewBool is a helper function to return a *bool to the specified value
func NewBool(b bool) *bool {
	out := b
	return &out
}

//NewInt32 is a helper function to return a *int32 to the specified value
func NewInt32(i int32) *int32 {
	return &i
}

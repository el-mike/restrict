package restrict

// Context - alias type for map of any values.
type Context map[string]interface{}

// AccessRequest - describes a Subject's intention to perform some Actions against
// given Resource.
type AccessRequest struct {
	Subject  Subject
	Resource Resource
	Actions  []string
	Context  Context
}

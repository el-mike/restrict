package restrict

// Context - alias type for a map of any values.
type Context map[string]interface{}

// AccessRequest - describes a Subject's intention to perform some Actions against
// given Resource.
type AccessRequest struct {
	// Subject - subject (typically a user) that wants to perform given Actions.
	// Needs to implement Subject interface.
	Subject Subject
	// Resource - resource that given Subject wants to interact with.
	// Needs to implement Resource interface.
	Resource Resource
	// Actions - list of operations Subject wants to perform on given Resource.
	Actions []string
	// Context - map of any additional values needed while checking the access.
	Context Context
	// SkipConditions - allows to skip Conditions while checking the access.
	SkipConditions bool
	// CompleteValidation - when true, validation will not return early, and all possible errors
	// will be returned, including all Conditions checks.
	CompleteValidation bool
}

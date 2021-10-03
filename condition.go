package restrict

// Condition - additional requirement that needs to be satisfied
// to grant given permission.
type Condition interface {
	// Name - returns Condition's name, which is it's unique
	// identifier.
	Name() string

	// Check - returns true if Condition is satisfied by
	// given request, false otherwise.
	Check(interface{}, *AccessRequest) bool
}

// Conditions - alias type for Conditions map.
type Conditions map[string]Condition

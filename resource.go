package restrict

// Resource - interface that needs to be implemented by any entity
// which acts as a resource in the system.
type Resource interface {
	// GetResourceName - returns a Resource's name. Should be the same as the one
	// used in PolicyDefinition.
	GetResourceName() string
}

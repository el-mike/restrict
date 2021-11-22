package restrict

// Resource - interface that needs to be implemented by any entity
// which acts as a resource in the system.
type Resource interface {
	// GetResourceName - returns a Resource's name. Should be the same as the one
	// used in PolicyDefinition.
	GetResourceName() string
}

// OwnableResource - interface that can be implemented by any resource
// when its ownership needs to tested in some Condition.
type OwnableResource interface {
	// GetOwner - returns a Resource's owner id.
	GetOwner() SubjectId
}

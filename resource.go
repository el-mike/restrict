package restrict

// Resource - interface that needs to be implemented by any entity
// which acts as a resource in the system.
type Resource interface {
	// GetResourceName - returns a Resource's name. Should be the same as the one
	// used in PolicyDefinition.
	GetResourceName() string
}

// baseResource - Resource implementation, to be used when proper Resource
// is impossible or not feasible to obtain.
type baseResource struct {
	name string
}

// GetResourceName - Resource interface implementation.
func (br *baseResource) GetResourceName() string {
	return br.name
}

// UseResource - returns baseResource instance.
func UseResource(name string) *baseResource {
	return &baseResource{
		name: name,
	}
}

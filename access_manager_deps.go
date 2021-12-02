package restrict

// PolicyProvider - interface for an entity that will provide Role configuration
// for AccessProvider.
type PolicyProvider interface {
	GetRole(roleID string) (*Role, error)
}

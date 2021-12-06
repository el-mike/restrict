package restrict

// StorageAdapter - interface for an entity that will provide persistance
// logic for PolicyDefinition.
type StorageAdapter interface {
	// LoadPolicy - loads and returns PolicyDefinition from underlying
	// storage provider.
	LoadPolicy() (*PolicyDefinition, error)

	// SavePolicy - saves PolicyDefinition in underlying storage provider.
	SavePolicy(policy *PolicyDefinition) error
}

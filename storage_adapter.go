package restrict

// StorageAdapter - allows to implement logic for persisting and loading
// given PolicyDefinition.
type StorageAdapter interface {
	// LoadPolicy - loads and returns PolicyDefinition from underlying
	// storage provider.
	LoadPolicy() (*PolicyDefinition, error)

	// SavePolicy - saves PolicyDefinition in underlying storage provider.
	SavePolicy(policy *PolicyDefinition) error
}

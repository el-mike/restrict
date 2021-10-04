package restrict

// PolicyManager - an entity responsible for managing policy. It uses passed StorageAdapter
// to save any changes made to policy.
type PolicyManager struct {
	// StorageAdapter used to load and save policy.
	adapter StorageAdapter

	// If set to true, PolicyManager will use it's StorageAdapter to save
	// the policy every time any change is being made.
	autoUpdate bool

	// PolicyDefinition currently loaded into memory and managed by
	// PolicyManager. Comes from StorageAdapter passed while creating PolicyManager.
	policy *PolicyDefinition
}

// NewPolicyManager - returns new PolicyManager instance and loads PolicyDefinition
// using passed StorageAdapter.
func NewPolicyManager(adapter StorageAdapter, autoUpdate bool) (*PolicyManager, error) {
	policy, err := adapter.LoadPolicy()
	if err != nil {
		return nil, err
	}

	manager := &PolicyManager{
		adapter:    adapter,
		autoUpdate: autoUpdate,
		policy:     policy,
	}

	return manager, nil
}

// GetPolicy - returns currently loaded PolicyDefinition.
func (pm *PolicyManager) GetPolicy() *PolicyDefinition {
	return pm.policy
}

// DisableAutoUpdate - disables automatic update.
func (pm *PolicyManager) DisableAutoUpdate() {
	pm.autoUpdate = false
}

// EnableAutoUpdate - enabled automatic update.
func (pm *PolicyManager) EnableAutoUpdate() {
	pm.autoUpdate = true
}

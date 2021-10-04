package adapters

import "github.com/el-Mike/restrict"

// InMemoryAdapter - simple in-memory adapter, for handling hard-coded policy.
type InMemoryAdapter struct {
	policy *restrict.PolicyDefinition
}

// NewInMemoryAdapter - returns new InMemoryAdapter instance.
func NewInMemoryAdapter(policy *restrict.PolicyDefinition) *InMemoryAdapter {
	return &InMemoryAdapter{
		policy: policy,
	}
}

// LoadPolicy - returns policy from memory.
func (ia *InMemoryAdapter) LoadPolicy() (*restrict.PolicyDefinition, error) {
	return ia.policy, nil
}

// SavePolicy - saves policy to memory.
func (ia *InMemoryAdapter) SavePolicy(policy *restrict.PolicyDefinition) error {
	ia.policy = policy

	return nil
}

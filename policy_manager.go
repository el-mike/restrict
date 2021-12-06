package restrict

import "sync"

// PolicyManager - an entity responsible for managing PolicyDefinition. It uses passed StorageAdapter
// for policy persistance.
type PolicyManager struct {
	// StorageAdapter used to load and save policy.
	adapter StorageAdapter

	// If set to true, PolicyManager will use it's StorageAdapter to save
	// the policy every time any change is made.
	autoUpdate bool

	// PolicyDefinition currently loaded into memory.
	policy *PolicyDefinition

	// PolicyManager should thread-safe for writing operations, therefore it uses RWMutex.
	sync.RWMutex
}

// NewPolicyManager - returns new PolicyManager instance and loads PolicyDefinition
// using passed StorageAdapter.
func NewPolicyManager(adapter StorageAdapter, autoUpdate bool) (*PolicyManager, error) {
	manager := &PolicyManager{
		adapter:    adapter,
		autoUpdate: autoUpdate,
	}

	// Load and initialize the policy.
	if err := manager.LoadPolicy(); err != nil {
		return nil, err
	}

	return manager, nil
}

// LoadPolicy - proxy method for loading the policy via StorageAdapter set
// when creating PolicyManager instance.
// Calling this method will override currently loaded policy.
func (pm *PolicyManager) LoadPolicy() error {
	pm.Lock()
	defer pm.Unlock()

	policy, err := pm.adapter.LoadPolicy()
	if err != nil {
		return err
	}

	pm.policy = policy

	if err := pm.applyPresets(); err != nil {
		return err
	}

	return nil
}

// SavePolicy - proxy method for saving the policy via StorageAdapter set
// when creating PolicyManager instance.
func (pm *PolicyManager) SavePolicy() error {
	return pm.adapter.SavePolicy(pm.policy)
}

// GetPolicy - returns currently loaded PolicyDefinition.
func (pm *PolicyManager) GetPolicy() *PolicyDefinition {
	pm.RLock()
	defer pm.RUnlock()

	return pm.policy
}

// applyPresets - applies defined presets to Permissions that are not yet merged.
func (pm *PolicyManager) applyPresets() error {
	// For every Role, iterate over all Permissions for given Resource and
	// merge Permission with it's preset if defined.
	for _, role := range pm.policy.Roles {
		for _, grants := range role.Grants {
			for _, permission := range grants {
				if permission.Preset != "" {
					if err := pm.applyPreset(permission); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// applyPreset - applies defined preset to Permission.
func (pm *PolicyManager) applyPreset(permission *Permission) error {
	permissionPreset := pm.policy.PermissionPresets[permission.Preset]

	// If given preset does not exist, return an error.
	if permissionPreset == nil {
		return newPermissionPresetNotFoundError(permission.Preset)
	}

	// Otherwise, merge found preset into Permission.
	permission.mergePreset(permissionPreset)

	return nil
}

// GetRole - returns a Role with given ID from currently loaded PolicyDefiniton.
func (pm *PolicyManager) GetRole(roleID string) (*Role, error) {
	pm.RLock()
	defer pm.RUnlock()

	role := pm.getRole(roleID)
	// If given Role does not exists, return an error.
	if role == nil {
		return nil, newRoleNotFoundError(roleID)
	}

	return role, nil
}

// AddRole - adds a new role to the policy.
// Saves with StorageAdapter if autoUpdate is set to true.
func (pm *PolicyManager) AddRole(role *Role) error {
	pm.Lock()
	defer pm.Unlock()

	// Check if role already exists - if yes, return an error.
	if r := pm.getRole(role.ID); r != nil {
		return newRoleAlreadyExistsError(role.ID)
	}

	pm.policy.Roles[role.ID] = role

	// Since new Permissions with presets could be added, run ApplyPresets.
	if err := pm.applyPresets(); err != nil {
		return err
	}

	if pm.autoUpdate {
		return pm.adapter.SavePolicy(pm.policy)
	}

	return nil
}

// UpdateRole - updates existing Role in currently loaded policy.
// Saves with StorageAdapter if autoUpdate is set to true.
func (pm *PolicyManager) UpdateRole(role *Role) error {
	pm.Lock()
	defer pm.Unlock()

	// If given Role does not exists, return an error.
	if r := pm.getRole(role.ID); r == nil {
		return newRoleNotFoundError(role.ID)
	}

	pm.policy.Roles[role.ID] = role

	// Since new Permissions with presets could be added, run ApplyPresets.
	if err := pm.applyPresets(); err != nil {
		return err
	}

	if pm.autoUpdate {
		return pm.adapter.SavePolicy(pm.policy)
	}

	return nil
}

// UpsertRole - updates a Role if exists, adds new Role otherwise.
// Saves with StorageAdapter if autoUpdate is set to true.
func (pm *PolicyManager) UpsertRole(role *Role) error {
	if err := pm.UpdateRole(role); err != nil {
		if _, ok := err.(*RoleNotFoundError); ok {
			return pm.AddRole(role)
		}

		return err
	}

	return nil
}

// DeleteRole - removes a Role with given ID.
// Saves with StorageAdapter if autoUpdate is set to true.
func (pm *PolicyManager) DeleteRole(roleID string) error {
	pm.Lock()
	defer pm.Unlock()

	if pm.policy.Roles == nil {
		pm.policy.Roles = Roles{}
	}

	// If Role with given ID does not exist, return an error.
	if r := pm.getRole(roleID); r == nil {
		return newRoleNotFoundError(roleID)
	}

	delete(pm.policy.Roles, roleID)

	if pm.autoUpdate {
		return pm.adapter.SavePolicy(pm.policy)
	}

	return nil
}

// AddPermission - adds a new Permission for the Role and Resource with passed IDs.
// Saves with StorageAdapter if autoUpdate is set to true.
func (pm *PolicyManager) AddPermission(roleID, resourceID string, permission *Permission) error {
	pm.Lock()
	defer pm.Unlock()

	role := pm.getRole(roleID)
	// If role does not exist, return an error.
	if role == nil {
		return newRoleNotFoundError(roleID)
	}

	pm.ensurePermissionsArray(role, resourceID)

	role.Grants[resourceID] = append(role.Grants[resourceID], permission)

	// If added Permission has preset defined, apply it immediately.
	if permission.Preset != "" {
		if err := pm.applyPreset(permission); err != nil {
			return err
		}
	}

	if pm.autoUpdate {
		return pm.adapter.SavePolicy(pm.policy)
	}

	return nil
}

// DeletePermission - removes a Permission with given name for Role and Resource with
// passed IDs.
// Saves with StorageAdapter if autoUpdate is set to true.
func (pm *PolicyManager) DeletePermission(roleID, resourceID, action string) error {
	pm.Lock()
	defer pm.Unlock()

	role := pm.getRole(roleID)

	// If role does not exist, return an error.
	if role == nil {
		return newRoleNotFoundError(roleID)
	}

	pm.ensurePermissionsArray(role, resourceID)

	for i, permission := range role.Grants[resourceID] {
		if permission.Action == action {
			role.Grants[resourceID] = pm.deletePermissionFromSlice(role.Grants[resourceID], i)
		}
	}

	if pm.autoUpdate {
		return pm.adapter.SavePolicy(pm.policy)
	}

	return nil
}

// deletePermissionFromSlice - helper function for removing Permission under given index
// from Permissions slice.
func (pm *PolicyManager) deletePermissionFromSlice(grants []*Permission, index int) []*Permission {
	if index >= 0 {
		newGrants := make([]*Permission, 0)
		newGrants = append(newGrants, grants[:index]...)
		newGrants = append(newGrants, grants[index+1:]...)

		return newGrants
	}

	return grants
}

// AddPermissionPreset - adds new PermissionPreset to PolicyDefinition.
// Saves with StorageAdapter if autoUpdate is set to true.
func (pm *PolicyManager) AddPermissionPreset(preset *PermissionPreset) error {
	pm.Lock()
	defer pm.Unlock()

	// If there is already a preset with given name, return an error.
	if p := pm.getPermissionPreset(preset.Name); p != nil {
		return newPermissionPresetAlreadyExistsError(preset.Name)
	}

	if pm.policy.PermissionPresets == nil {
		pm.policy.PermissionPresets = PermissionPresets{}
	}

	pm.policy.PermissionPresets[preset.Name] = preset

	if pm.autoUpdate {
		return pm.adapter.SavePolicy(pm.policy)
	}

	return nil
}

// UpdatePermissionPreset - updates a PermissionPreset in PolicyDefinition.
// Saves with StorageAdapter if autoUpdate is set to true.
func (pm *PolicyManager) UpdatePermissionPreset(preset *PermissionPreset) error {
	pm.Lock()
	defer pm.Unlock()

	// If there is no preset with given name, return an error.
	if p := pm.getPermissionPreset(preset.Name); p == nil {
		return newPermissionPresetNotFoundError(preset.Name)
	}

	pm.policy.PermissionPresets[preset.Name] = preset

	if pm.autoUpdate {
		return pm.adapter.SavePolicy(pm.policy)
	}

	return nil
}

// UpsertPermissionPreset - updates Permission preset if exists, adds a new otherwise.
// Saves with StorageAdapter if autoUpdate is set to true.
func (pm *PolicyManager) UpsertPermissionPreset(preset *PermissionPreset) error {
	if err := pm.UpdatePermissionPreset(preset); err != nil {
		if _, ok := err.(*PermissionPresetNotFoundError); ok {
			return pm.AddPermissionPreset(preset)
		}

		return err
	}

	return nil
}

// DeletePermissionPreset - removes PermissionPreset with given name.
// Saves with StorageAdapter if autoUpdate is set to true.
func (pm *PolicyManager) DeletePermissionPreset(name string) error {
	pm.Lock()
	defer pm.Unlock()

	// If there is no preset with given name, return an error.
	if p := pm.getPermissionPreset(name); p == nil {
		return newPermissionPresetNotFoundError(name)
	}

	delete(pm.policy.PermissionPresets, name)

	if pm.autoUpdate {
		return pm.adapter.SavePolicy(pm.policy)
	}

	return nil
}

// DisableAutoUpdate - disables automatic update.
func (pm *PolicyManager) DisableAutoUpdate() {
	pm.autoUpdate = false
}

// EnableAutoUpdate - enables automatic update.
func (pm *PolicyManager) EnableAutoUpdate() {
	pm.autoUpdate = true
}

// ensurePermissionsArray - helper function for setting GrantsMap and Permissions array
// for given Role if they don't exist (i.e. are equal to nil).
func (pm *PolicyManager) ensurePermissionsArray(role *Role, resourceID string) {
	if role.Grants == nil {
		role.Grants = GrantsMap{}
	}

	if role.Grants[resourceID] == nil {
		role.Grants[resourceID] = []*Permission{}
	}
}

// getRole - helper function for getting a Role with given ID.
func (pm *PolicyManager) getRole(roleID string) *Role {
	role, ok := pm.policy.Roles[roleID]

	if !ok {
		return nil
	}

	return role
}

// getPermissionPreset - helper function for getting PermissionPreset from PolicyDefinition.
func (pm *PolicyManager) getPermissionPreset(name string) *PermissionPreset {
	preset, ok := pm.policy.PermissionPresets[name]

	if !ok {
		return nil
	}

	return preset
}

package restrict

import "sync"

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

	if err := pm.ApplyPresets(); err != nil {
		return err
	}

	return nil
}

// ApplyPresets - applies defined presets to Permissions that are not yet merged.
func (pm *PolicyManager) ApplyPresets() error {
	// For every Role, iterate over all Permissions for given Resource and
	// merge Permission with it's preset if defined.
	for _, role := range pm.policy.Roles {
		for _, grants := range role.Grants {
			for _, permission := range grants {
				if permission.Preset != "" {
					return pm.ApplyPreset(permission)
				}
			}
		}
	}

	return nil
}

func (pm *PolicyManager) ApplyPreset(permission *Permission) error {
	permissionPreset := pm.policy.PermissionPresets[permission.Preset]

	// If given preset does not exist, return an error.
	if permissionPreset == nil {
		return NewPermissionPresetNotFoundError(permission.Preset)
	}

	// Otherwise, merge found preset into Permission.
	permission.MergePreset(permissionPreset)

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

// GetPermissionPresets - returns a map of Permission presets defined in PolicyDefinition.
func (pm *PolicyManager) GetPermissionPresets() *PermissionPresets {
	pm.RLock()
	defer pm.RUnlock()

	return &pm.policy.PermissionPresets
}

// GetRole - returns a Role with given ID from currently loaded PolicyDefiniton.
func (pm *PolicyManager) GetRole(roleID string) (*Role, error) {
	pm.RLock()
	defer pm.RUnlock()

	role := pm.getRole(roleID)
	// If given Role does not exists, return an error.
	if role == nil {
		return nil, NewRoleNotFoundError(roleID)
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
		return NewRoleAlreadyExistsError(role.ID)
	}

	pm.policy.Roles[role.ID] = role

	// Since new Permissions with presets could be added, run ApplyPresets.
	if err := pm.ApplyPresets(); err != nil {
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
		return NewRoleNotFoundError(role.ID)
	}

	pm.policy.Roles[role.ID] = role

	// Since new Permissions with presets could be added, run ApplyPresets.
	if err := pm.ApplyPresets(); err != nil {
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

// DeleteRole - removes a role with given ID.
// Saves with StorageAdapter if autoUpdate is set to true.
func (pm *PolicyManager) DeleteRole(roleID string) error {
	pm.Lock()
	defer pm.Unlock()

	if pm.policy.Roles == nil {
		pm.policy.Roles = Roles{}
	}

	// If Role with given ID does not exist, return an error.
	if r := pm.getRole(roleID); r == nil {
		return NewRoleNotFoundError(roleID)
	}

	delete(pm.policy.Roles, roleID)

	if pm.autoUpdate {
		return pm.adapter.SavePolicy(pm.policy)
	}

	return nil
}

// AddPermission - adds a new Permission for the Role and resource with passed IDs.
// Saves with StorageAdapter if autoUpdate is set to true.
func (pm *PolicyManager) AddPermission(roleID, resourceID string, permission *Permission) error {
	pm.Lock()
	defer pm.Unlock()

	// If name for new Permission could not be resolved, return an error.
	if err := permission.ResolveName(); err != nil {
		return err
	}

	role := pm.getRole(roleID)
	// If role does not exist, return an error.
	if role == nil {
		return NewRoleNotFoundError(role.ID)
	}

	pm.ensurePermissionsArray(role, resourceID)

	// If there is already a permission with given name for given resource,
	// return an error.
	if p := pm.getPermission(role.ID, resourceID, permission.Name); p != nil {
		return NewPermissionAlreadyExistsError(resourceID, permission.Name)
	}

	role.Grants[resourceID] = append(role.Grants[resourceID], permission)

	// If added Permission has preset defined, apply it immediately.
	if permission.Preset != "" {
		if err := pm.ApplyPreset(permission); err != nil {
			return err
		}
	}

	if pm.autoUpdate {
		return pm.adapter.SavePolicy(pm.policy)
	}

	return nil
}

// UpdatePermission - updates existing Permission in currently loaded policy.
// Saves with StorageAdapter if autoUpdate is set to true.
func (pm *PolicyManager) UpdatePermission(roleID, resourceID string, permission *Permission) error {
	pm.Lock()
	defer pm.Unlock()

	// If name for new Permission could not be resolved, return an error.
	if err := permission.ResolveName(); err != nil {
		return err
	}

	role := pm.getRole(roleID)
	// If role does not exist, return an error.
	if role == nil {
		return NewRoleNotFoundError(role.ID)
	}

	pm.ensurePermissionsArray(role, resourceID)

	if p := pm.getPermission(role.ID, resourceID, permission.Name); p == nil {
		return NewPermissionNotFoundError(resourceID, permission.Name)
	}

	index := -1

	for i, perm := range role.Grants[resourceID] {
		if perm.Name == permission.Name {
			index = i
			break
		}
	}

	if index >= 0 {
		role.Grants[resourceID][index] = permission
	}

	// If updated Permission has preset defined, apply it immediately.
	if err := pm.ApplyPreset(permission); err != nil {
		return err
	}

	if pm.autoUpdate {
		return pm.adapter.SavePolicy(pm.policy)
	}

	return nil
}

// UpsertRole - updates a Permission if exists for given resource, adds new Permission otherwise.
// Saves with StorageAdapter if autoUpdate is set to true.
func (pm *PolicyManager) UpsertPermission(roleID, resourceID string, permission *Permission) error {
	if err := pm.UpdatePermission(roleID, resourceID, permission); err != nil {
		if _, ok := err.(*PermissionNotFoundError); ok {
			return pm.AddPermission(roleID, resourceID, permission)
		}

		return err
	}

	return nil
}

// DeletePermission - removes a Permission with given name for Role and resource with
// passed IDs.
// Saves with StorageAdapter if autoUpdate is set to true.
func (pm *PolicyManager) DeletePermission(roleID, resourceID, name string) error {
	pm.Lock()
	defer pm.Unlock()

	role := pm.getRole(roleID)
	// If role does not exist, return an error.
	if role == nil {
		return NewRoleNotFoundError(role.ID)
	}

	pm.ensurePermissionsArray(role, resourceID)

	index := -1

	for i, permission := range role.Grants[resourceID] {
		if permission.Name == name {
			index = i
			break
		}
	}

	if index >= 0 {
		grants := role.Grants[resourceID]

		newGrants := make([]*Permission, 0)
		newGrants = append(newGrants, grants[:index]...)
		newGrants = append(newGrants, grants[index+1:]...)

		role.Grants[resourceID] = newGrants
	}

	if pm.autoUpdate {
		return pm.adapter.SavePolicy(pm.policy)
	}

	return nil
}

// AddPermissionPreset - adds new Permission preset to PolicyDefinition.
// Saves with StorageAdapter if autoUpdate is set to true.
func (pm *PolicyManager) AddPermissionPreset(permission *Permission) error {
	pm.Lock()
	defer pm.Unlock()

	// If name for new Permission preset could not be resolved, return an error.
	if err := permission.ResolveName(); err != nil {
		return err
	}

	// If there is already a preset with given name, return an error.
	if p := pm.getPermissionPreset(permission.Name); p != nil {
		return NewPermissionPresetAlreadyExistsError(permission.Name)
	}

	pm.policy.PermissionPresets[permission.Name] = permission

	if pm.autoUpdate {
		return pm.adapter.SavePolicy(pm.policy)
	}

	return nil
}

// UpdatePermissionPreset - updates a Permission preeset in PolicyDefinition.
// Saves with StorageAdapter if autoUpdate is set to true.
func (pm *PolicyManager) UpdatePermissionPreset(permission *Permission) error {
	pm.Lock()
	defer pm.Unlock()

	// If name for new Permission preset could not be resolved, return an error.
	if err := permission.ResolveName(); err != nil {
		return err
	}

	// If there is no preset with given name, return an error.
	if p := pm.getPermissionPreset(permission.Name); p == nil {
		return NewPermissionPresetNotFoundError(permission.Name)
	}

	pm.policy.PermissionPresets[permission.Name] = permission

	if pm.autoUpdate {
		return pm.adapter.SavePolicy(pm.policy)
	}

	return nil
}

// UpsertPermissionPreset - updates Permission preset if exists, add a new otherwise.
// Saves with StorageAdapter if autoUpdate is set to true.
func (pm *PolicyManager) UpsertPermissionPreset(permission *Permission) error {
	if err := pm.UpdatePermissionPreset(permission); err != nil {
		if _, ok := err.(*PermissionPresetNotFoundError); ok {
			return pm.AddPermissionPreset(permission)
		}

		return err
	}

	return nil
}

// DeletePermissionPreset - removes Permission preset with given name.
// Saves with StorageAdapter if autoUpdate is set to true.
func (pm *PolicyManager) DeletePermissionPreset(name string) error {
	pm.Lock()
	defer pm.Unlock()

	// If there is no preset with given name, return an error.
	if p := pm.getPermissionPreset(name); p == nil {
		return NewPermissionPresetNotFoundError(name)
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

// EnableAutoUpdate - enabled automatic update.
func (pm *PolicyManager) EnableAutoUpdate() {
	pm.autoUpdate = true
}

// ensurePermissionsArray - helper function for setting GrantsMap and Permissions array
// for given role if they don't exist (i.e. are equal to nil).
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

// getPermission - helper function for getting a Permission with given name
// under resource and Role with passed IDs.
func (pm *PolicyManager) getPermission(roleID string, resourceID string, name string) *Permission {
	for _, permission := range pm.policy.Roles[roleID].Grants[resourceID] {
		if permission.Name == name {
			return permission
		}
	}

	return nil
}

// getPermissionPreset - helper function for getting Permission preset from PolicyDefinition.
func (pm *PolicyManager) getPermissionPreset(name string) *Permission {
	preset, ok := pm.policy.PermissionPresets[name]

	if !ok {
		return nil
	}

	return preset
}

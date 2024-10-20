package restrict

// Subject - interface that has to be implemented by any entity
// which authorization needs to be checked.
type Subject interface {
	// GetRoles - returns a Subject's role.
	GetRoles() []string
}

// baseSubject - Subject implementation, to be used when proper entity
// is impossible or not feasible to obtain.
type baseSubject struct {
	roles []string
}

// GetRoles - Subject interface implementation.
func (bs *baseSubject) GetRoles() []string {
	return bs.roles
}

// UseSubject - returns baseSubject instance.
func UseSubject(roles []string) *baseSubject {
	return &baseSubject{
		roles,
	}
}

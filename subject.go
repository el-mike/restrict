package restrict

// Subject - interface that has to be implemented by any entity
// which authorization needs to be checked.
type Subject interface {
	// GetRole - returns a Subject's role.
	GetRole() string
}

// baseSubject - Subject implementation, to be used when proper entity
// is impossible or not feasible to obtain.
type baseSubject struct {
	role string
}

// GetRole - Subject interface implementation.
func (bs *baseSubject) GetRole() string {
	return bs.role
}

// UseSubject - returns baseSubject instance.
func UseSubject(role string) *baseSubject {
	return &baseSubject{
		role: role,
	}
}

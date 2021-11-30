package restrict

// SubjectId - alias type representing any type that could use as Subject's id.
type SubjectId interface{}

// Subject - interface that needs to be implemented by any entity
// which can perform Actions against Resources.
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

package restrict

// SubjectId - alias type representing any type that could use as Subject's id.
type SubjectId interface{}

// Subject - interface that needs to be implemented by any entity
// which can perform Actions against Resources.
type Subject interface {
	// GetRole - returns a Subject's role.
	GetRole() string
}

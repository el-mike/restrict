package examples

type User struct {
	ID string
}

// Subject interface implementation.
func (u *User) GetRoles() []string {
	return []string{"User"}
}

type Conversation struct {
	ID            string
	CreatedBy     string
	Participants  []string
	MessagesCount int
	Active        bool
}

// Resource interface implementation.
func (c *Conversation) GetResourceName() string {
	return "Conversation"
}

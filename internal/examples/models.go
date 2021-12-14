package examples

type User struct {
	ID string
}

// Subject interface implementation.
func (u *User) GetRole() string {
	return "User"
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

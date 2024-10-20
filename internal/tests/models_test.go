package tests

const (
	BasicUserRole = "BasicUser"
	UserRole      = "User"
	AdminRole     = "Admin"

	UserResource         = "User"
	ConversationResource = "Conversation"
	MessageResource      = "Message"
)

type User struct {
	ID string

	Roles []string
}

// Subject interface implementation.
func (u *User) GetRoles() []string {
	return u.Roles
}

// Resource interface implementation. User can be both Subject and Resource.
func (u *User) GetResourceName() string {
	return UserResource
}

type Conversation struct {
	ID string

	CreatedBy     string
	Participants  []string
	MessagesCount int
	Active        bool
}

// Resource interface implementation.
func (c *Conversation) GetResourceName() string {
	return ConversationResource
}

type Message struct {
	ID string

	CreatedBy     string
	CoversationId string
}

func (c *Message) GetResourceName() string {
	return MessageResource
}

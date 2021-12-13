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

	Role string
}

// Subject interface implementation.
func (u *User) GetRole() string {
	return u.Role
}

// Resource interface implementation. User can be both Subject and Resource.
func (u *User) GetResourceName() string {
	return UserResource
}

type Conversation struct {
	ID string

	CreatedBy    string
	Participants []string
	Active       bool
	MessageIds   []string
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

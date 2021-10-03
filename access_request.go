package restrict

type Context map[string]interface{}

type AccessRequest struct {
	Role     string
	Resource string
	Actions  []string
	Context  Context
}

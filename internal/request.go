package internal

// TODO add replies from server

type RequestType int

const (
	Identity RequestType = iota
	ListClients
	Message
)

type Request struct {
	Type    RequestType
	Message string // optional
}

func NewRequest() Request {
	return Request{}
}

func BuildIdentityRequest()    {}
func BuildListClientsRequest() {}
func BuildMessageRequest()     {}

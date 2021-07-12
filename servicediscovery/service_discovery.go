package servicediscovery

type Server struct {
	ID             string
	Type           string
	Addr           string
	ClientHandlers []string
	ServerHandlers []string
}

type ServerListener interface {
	OnAddServer(*Server)
	OnRemoveServer(*Server)
}

type ServiceDiscovery interface {
	GetRandomServer(serverType string) (*Server, error)
	GetServer(serverID string) (*Server, error)
	PullServers(first bool) error
	AddServerListener(sl ServerListener)
}

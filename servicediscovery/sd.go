package servicediscovery

type Server struct {
	ID             string
	Type           string
	Addr           string
	ClientHandlers []string
	ServerHandlers []string
}

type ServiceDiscovery interface {
	GetRandomServer(serverType string) (*Server, error)
	GetServer(serverID string) (*Server, error)
}

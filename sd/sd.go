package sd

type Server struct {
	ID             string
	Type           string
	Addr           string
	ClientHandlers []string
	ServerHandlers []string
}

var SD = newEtcdSD()

type ServiceDiscovery interface {
	GetServersByType(serverType string) (map[string]*Server, error)
	GetServer(serverID string) (*Server, error)
	GetServers() []*Server
}

func GetRandomServer(map[string]*Server) *Server {
	return nil
}

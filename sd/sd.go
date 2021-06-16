package sd

type Server struct {
	ID             string
	ServerType     string
	Addr           string
	ClientHandlers []string
	ServerHandlers []string
}

var SD ServiceDiscovery

// todo watcher and local cache
type ServiceDiscovery interface {
	GetServersByType(serverType string) (map[string]*Server, error)
	GetServer(id string) (*Server, error)
	GetServers() []*Server
	SyncServers(firstSync bool) error
}

func GetRandomServer(map[string]*Server) *Server {
	return nil
}

func GetMinLoadServer(map[string]*Server) *Server {
	return nil
}

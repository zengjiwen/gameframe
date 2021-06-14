package servicediscovery

type Server struct {
	ID             string
	ServerType     string
	Addr           string
	ClientHandlers []string
	ServerHandlers []string
}

// todo watcher and local cache
type ServiceDiscovery interface {
	GetServersByType(serverType string) (map[string]*Server, error)
	GetServer(id string) (*Server, error)
	GetServers() []*Server
	SyncServers(firstSync bool) error
}

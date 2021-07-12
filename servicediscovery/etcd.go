package servicediscovery

type etcd struct {
}

func NewEtcd() ServiceDiscovery {
	return &etcd{}
}

func (e *etcd) GetRandomServer(serverType string) (*Server, error) {
	panic("implement me")
}

func (e *etcd) GetServer(id string) (*Server, error) {
	panic("implement me")
}

func (e *etcd) PullServers(first bool) error {
	panic("implement me")
}

func (e *etcd) AddServerListener(sl ServerListener) {
	panic("implement me")
}

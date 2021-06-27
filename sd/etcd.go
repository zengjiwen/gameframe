package sd

type etcdSD struct {
}

func newEtcdSD() *etcdSD {
	return &etcdSD{}
}

func (e *etcdSD) GetServersByType(serverType string) (map[string]*Server, error) {
	panic("implement me")
}

func (e *etcdSD) GetServer(id string) (*Server, error) {
	panic("implement me")
}

func (e *etcdSD) GetServers() []*Server {
	panic("implement me")
}

func (e *etcdSD) SyncServers(firstSync bool) error {
	panic("implement me")
}

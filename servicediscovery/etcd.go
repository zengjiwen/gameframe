package servicediscovery

import (
	"context"
	"github.com/zengjiwen/gameframe/env"
	"go.etcd.io/etcd/clientv3"
	"time"
)

type etcd struct {
	client  *clientv3.Client
	leaseID clientv3.LeaseID
	sl      ServerListener
	dieChan chan struct{}
}

func NewEtcd(addr string) ServiceDiscovery {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{addr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	sd := &etcd{
		client:  client,
		leaseID: clientv3.NoLease,
		dieChan: make(chan struct{}),
	}
	return sd
}

func (e *etcd) Init() error {
	resp, err := e.client.Create(context.Background(), int64(5*time.Second))
	if err != nil {
		return err
	}

	e.leaseID = clientv3.LeaseID(resp.ID)
	c, err := e.client.KeepAlive(context.Background(), e.leaseID)
	if err != nil {
		return err
	}

	<-c
	go e.keepAlive(c)
	return nil
}

func (e *etcd) keepAlive(c <-chan *clientv3.LeaseKeepAliveResponse) {
	for {
		select {
		case _, ok := <-c:
			if !ok {
				env.DieChan <- struct{}{}
				return
			}
		case <-e.dieChan:
			return
		}
	}
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
	e.sl = sl
}

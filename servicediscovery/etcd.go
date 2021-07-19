package servicediscovery

import (
	"context"
	"github.com/coreos/etcd/storage/storagepb"
	"github.com/zengjiwen/gameframe/env"
	"go.etcd.io/etcd/clientv3"
	"math/rand"
	"sync"
	"time"
)

type etcd struct {
	client  *clientv3.Client
	leaseID clientv3.LeaseID
	sl      []ServerListener
	dieChan chan struct{}
	err     error

	mu                sync.RWMutex
	serverInfos       map[string]*ServerInfo
	serverInfosByType map[string][]*ServerInfo
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
		client:            client,
		leaseID:           clientv3.NoLease,
		dieChan:           make(chan struct{}),
		serverInfos:       make(map[string]*ServerInfo),
		serverInfosByType: make(map[string][]*ServerInfo),
	}
	return sd
}

func (e *etcd) Init() error {
	_serverInfo = newServerInfo(env.ServerID, env.ServerType, env.ServiceAddr)

	e.putMyself()
	e.lease()
	e.getServers()
	go e.watch()
	return e.err
}

func (e *etcd) lease() {
	if e.err != nil {
		return
	}

	resp, err := e.client.Create(context.Background(), int64(5*time.Second))
	if err != nil {
		e.err = err
		return
	}

	e.leaseID = clientv3.LeaseID(resp.ID)
	c, err := e.client.KeepAlive(context.Background(), e.leaseID)
	if err != nil {
		e.err = err
		return
	}

	<-c
	go e.keepAlive(c)
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

func (e *etcd) watch() {
	if e.err != nil {
		return
	}

	c := e.client.Watch(context.Background(), _sdPrefix, clientv3.WithPrefix())
	for {
		select {
		case resp, ok := <-c:
			if !ok || resp.Err() != nil {
				env.DieChan <- struct{}{}
				return
			}

			e.mu.Lock()
			for _, event := range resp.Events {
				serverInfo, err := parseSDValue(event.Kv.Value)
				if err != nil {
					continue
				}

				switch event.Type {
				case storagepb.PUT:
					e.addServer(serverInfo)
				case storagepb.DELETE:
					e.removeServer(serverInfo)
				}
			}
			e.mu.Unlock()
		case <-e.dieChan:
			return
		}
	}
}

func (e *etcd) putMyself() {
	if e.err != nil {
		return
	}

	value, err := genSDValue(_serverInfo)
	if err != nil {
		e.err = err
		return
	}

	_, err = e.client.Put(context.Background(),
		genSDKey(env.ServerID, env.ServerType),
		value, clientv3.WithLease(e.leaseID))
}

func (e *etcd) getServers() {
	if e.err != nil {
		return
	}

	resp, err := e.client.Get(context.Background(), _sdPrefix, clientv3.WithPrefix())
	if err != nil {
		e.err = err
		return
	}

	for _, kv := range resp.Kvs {
		serverInfo, err := parseSDValue(kv.Value)
		if err != nil {
			continue
		}
		if serverInfo.ID == _serverInfo.ID {
			continue
		}

		e.doAddServer(serverInfo)
	}
}

func (e *etcd) addServer(serverInfo *ServerInfo) {
	e.mu.Lock()
	e.doOnlyAddServer(serverInfo)
	e.mu.Unlock()

	for _, listener := range e.sl {
		listener.OnAddServer(serverInfo)
	}
}

func (e *etcd) doOnlyAddServer(serverInfo *ServerInfo) {
	e.serverInfos[serverInfo.ID] = serverInfo
	typeServerInfos, ok := e.serverInfosByType[serverInfo.Type]
	if !ok {
		typeServerInfos = make([]*ServerInfo, 0)
		e.serverInfosByType[serverInfo.Type] = typeServerInfos
	}
	typeServerInfos = append(typeServerInfos, serverInfo)
}

func (e *etcd) doAddServer(serverInfo *ServerInfo) {
	e.doOnlyAddServer(serverInfo)
	for _, listener := range e.sl {
		listener.OnAddServer(serverInfo)
	}
}

func (e *etcd) removeServer(serverInfo *ServerInfo) {
	e.mu.Lock()
	e.doOnlyRemoveServer(serverInfo)
	e.mu.Unlock()

	for _, listener := range e.sl {
		listener.OnRemoveServer(serverInfo)
	}
}

func (e *etcd) doOnlyRemoveServer(serverInfo *ServerInfo) {
	delete(e.serverInfos, serverInfo.ID)
	if typeServerInfos, ok := e.serverInfosByType[serverInfo.Type]; ok {
		for i, info := range typeServerInfos {
			if info.ID == serverInfo.ID {
				typeServerInfos = append(typeServerInfos[:i], typeServerInfos[i+1:]...)
				break
			}
		}
	}
}

func (e *etcd) doRemoveServer(serverInfo *ServerInfo) {
	e.doOnlyRemoveServer(serverInfo)
	for _, listener := range e.sl {
		listener.OnRemoveServer(serverInfo)
	}
}

func (e *etcd) GetRandomServer(serverType string) (*ServerInfo, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	typeServerInfos, ok := e.serverInfosByType[serverType]
	if !ok {
		return nil, false
	}

	return typeServerInfos[rand.Intn(len(typeServerInfos))], true
}

func (e *etcd) GetServer(serverID string) (*ServerInfo, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	serverInfo, ok := e.serverInfos[serverID]
	return serverInfo, ok
}

func (e *etcd) AddServerListener(sl ServerListener) {
	e.sl = append(e.sl, sl)
}

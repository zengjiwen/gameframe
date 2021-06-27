package rpc

import (
	"github.com/zengjiwen/gameframe/sd"
	"google.golang.org/grpc"
	"sync"
)

var (
	_mu      sync.RWMutex
	_clients map[string]*grpc.ClientConn
)

func ClientByServerID(serverID string) (*grpc.ClientConn, bool) {
	_mu.RLock()
	defer _mu.RUnlock()

	client, ok := _clients[serverID]
	return client, ok
}

func OnServerAdded(server *sd.Server) {
	_mu.RLock()
	_, ok := _clients[server.ID]
	if ok {
		_mu.RUnlock()
		return
	}
	_mu.RUnlock()

	conn, err := grpc.Dial(server.Addr, grpc.WithInsecure())
	if err != nil {
		return
	}

	_mu.Lock()
	_clients[server.ID] = conn
	_mu.Unlock()
}

func OnServerRemoved(server *sd.Server) {
	_mu.Lock()
	delete(_clients, server.ID)
	_mu.Unlock()
}

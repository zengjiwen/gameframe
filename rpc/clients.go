package rpc

import (
	"github.com/zengjiwen/gameframe/rpc/protos"
	"github.com/zengjiwen/gameframe/servicediscovery"
	"google.golang.org/grpc"
	"sync"
)

var (
	_mu      sync.RWMutex
	_clients map[string]protos.RPCClient
)

func ClientByServerID(serverID string) (protos.RPCClient, bool) {
	_mu.RLock()
	defer _mu.RUnlock()

	client, ok := _clients[serverID]
	return client, ok
}

func OnAddServer(server *servicediscovery.Server) {
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
	_clients[server.ID] = protos.NewRPCClient(conn)
	_mu.Unlock()
}

func OnRemoveServer(server *servicediscovery.Server) {
	_mu.Lock()
	delete(_clients, server.ID)
	_mu.Unlock()
}

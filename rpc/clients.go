package rpc

import (
	"google.golang.org/grpc"
	"sync"
)

var _mu sync.RWMutex
var _clients map[string]*grpc.ClientConn

func ClientByServerID(serverID string) (*grpc.ClientConn, bool) {
	_mu.RLock()
	defer _mu.RUnlock()

	client, ok := _clients[serverID]
	return client, ok
}

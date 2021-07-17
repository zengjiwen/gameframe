package rpc

import (
	"github.com/zengjiwen/gameframe/rpc/protos"
	"github.com/zengjiwen/gameframe/servicediscovery"
	"google.golang.org/grpc"
	"sync"
)

var Clients = &clients{
	conns: make(map[string]*grpc.ClientConn),
}

type clients struct {
	mu    sync.RWMutex
	conns map[string]*grpc.ClientConn
}

func (c *clients) ClientByServerID(serverID string) (protos.RPCClient, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	conn, ok := c.conns[serverID]
	return protos.NewRPCClient(conn), ok
}

func (c *clients) OnAddServer(server *servicediscovery.Server) {
	c.mu.RLock()
	_, ok := c.conns[server.ID]
	if ok {
		c.mu.RUnlock()
		return
	}
	c.mu.RUnlock()

	conn, err := grpc.Dial(server.Addr, grpc.WithInsecure())
	if err != nil {
		return
	}

	c.mu.Lock()
	c.conns[server.ID] = conn
	c.mu.Unlock()
}

func (c *clients) OnRemoveServer(server *servicediscovery.Server) {
	c.mu.Lock()
	conn, ok := c.conns[server.ID]
	delete(c.conns, server.ID)
	c.mu.Unlock()

	if ok {
		conn.Close()
	}
}

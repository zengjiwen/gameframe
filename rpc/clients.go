package rpc

import (
	"context"
	"github.com/zengjiwen/gameframe/rpc/protos"
	"github.com/zengjiwen/gameframe/servicediscovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
	"time"
)

var _clients *clients

type clients struct {
	mu    sync.RWMutex
	conns map[string]*grpc.ClientConn
}

func (c *clients) OnAddServer(server *servicediscovery.ServerInfo) {
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

func (c *clients) OnRemoveServer(server *servicediscovery.ServerInfo) {
	c.mu.Lock()
	conn, ok := c.conns[server.ID]
	delete(c.conns, server.ID)
	c.mu.Unlock()

	if ok {
		conn.Close()
	}
}

func InitClients() {
	_clients = &clients{
		conns: make(map[string]*grpc.ClientConn),
	}
}

func CloseClients() error {
	var err error
	_clients.mu.Lock()
	for _, conn := range _clients.conns {
		if err == nil {
			err = conn.Close()
		} else {
			conn.Close()
		}
	}
	_clients.mu.Unlock()
	_clients.conns = make(map[string]*grpc.ClientConn)
	return err
}

func GetConn(serverID string) (*grpc.ClientConn, bool) {
	_clients.mu.RLock()
	defer _clients.mu.RUnlock()

	conn, ok := _clients.conns[serverID]
	return conn, ok
}

func WatchServer() {
	servicediscovery.Get().AddServerWatcher(_clients)
}

func RemoveConn(serverID string) {
	_clients.mu.Lock()
	conn, ok := _clients.conns[serverID]
	delete(_clients.conns, serverID)
	_clients.mu.Unlock()

	if ok {
		conn.Close()
	}
}

func TryBestCall(serverID string, client protos.RPCClient, request *protos.CallRequest) (*protos.CallRespond, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var tmpDelay time.Duration
	var retryCount int
	for {
		resp, err := client.Call(ctx, request)
		if err != nil {
			if s, ok := status.FromError(err); ok && s.Code() == codes.DeadlineExceeded {
				retryCount++
				if retryCount > 10 {
					RemoveConn(serverID)
					return nil, err
				}

				if tmpDelay == 0 {
					tmpDelay = 5 * time.Millisecond
				} else {
					tmpDelay *= 2
				}
				if max := 1 * time.Second; tmpDelay > max {
					tmpDelay = max
				}
				time.Sleep(tmpDelay)
				continue
			}
			return nil, err
		}
		return resp, nil
	}
}

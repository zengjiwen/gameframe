package proxy

import (
	"context"
	"github.com/zengjiwen/gameframe/env"
	"github.com/zengjiwen/gameframe/rpc"
	"github.com/zengjiwen/gameframe/rpc/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type backend struct {
	serverID  string
	rpcClient protos.RPCClient
}

func NewBackend(serverID string, conn *grpc.ClientConn) Proxy {
	return &backend{
		serverID:  serverID,
		rpcClient: protos.NewRPCClient(conn),
	}
}

func (b *backend) Send(route string, payload []byte) error {
	request := &protos.SendRequest{
		Route:   route,
		Payload: payload,
	}

	ctx, cancel := context.WithTimeout(context.Background(), env.RPCTimeout)
	defer cancel()

	var tmpDelay time.Duration
	var retryCount int
	for {
		_, err := b.rpcClient.Send(ctx, request)
		if err != nil {
			if s, ok := status.FromError(err); ok && s.Code() == codes.DeadlineExceeded {
				retryCount++
				if retryCount > 10 {
					rpc.RemoveConn(b.serverID)
					return err
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
			return err
		}
		return nil
	}
}

package proxy

import (
	"context"
	"errors"
	"github.com/zengjiwen/gameframe/rpc"
	"github.com/zengjiwen/gameframe/rpc/protos"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type backend struct {
	frontendServerID string
}

func NewBackend(frontendServerID string) Proxy {
	return &backend{
		frontendServerID: frontendServerID,
	}
}

func (b *backend) Send(route string, payload []byte) error {
	rpcConn, ok := rpc.GetConn(b.frontendServerID)
	if !ok {
		return errors.New("rpc conn not exist!")
	}
	rpcClient := protos.NewRPCClient(rpcConn)

	request := &protos.SendRequest{
		Route:   route,
		Payload: payload,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var tmpDelay time.Duration
	var retryCount int
	for {
		_, err := rpcClient.Send(ctx, request)
		if err != nil {
			if s, ok := status.FromError(err); ok && s.Code() == codes.DeadlineExceeded {
				retryCount++
				if retryCount > 10 {
					rpc.RemoveConn(b.frontendServerID)
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

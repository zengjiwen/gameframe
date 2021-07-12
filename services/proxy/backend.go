package proxy

import (
	"context"
	"github.com/zengjiwen/gameframe/rpc/protos"
)

type backend struct {
	rpcClient protos.RPCClient
}

func NewBackend(rpcClient protos.RPCClient) Proxy {
	return &backend{
		rpcClient: rpcClient,
	}
}

func (b *backend) Send(route string, payload []byte) error {
	request := &protos.SendRequest{
		Route:   route,
		Payload: payload,
	}

	// todo timeout
	_, err := b.rpcClient.Send(context.Background(), request)
	return err
}

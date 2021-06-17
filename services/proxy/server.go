package proxy

import (
	"context"
	"github.com/zengjiwen/gameframe/rpc/protos"
	"google.golang.org/grpc"
)

type server struct {
	conn      *grpc.ClientConn
	rpcClient protos.RPCClient
}

func NewServer(conn *grpc.ClientConn) Proxy {
	return &server{
		conn:      conn,
		rpcClient: protos.NewRPCClient(conn),
	}
}

func (s *server) Send(route string, payload []byte) error {
	request := &protos.SendRequest{
		Route:   route,
		Payload: payload,
	}

	_, err := s.rpcClient.Send(context.Background(), request)
	return err
}

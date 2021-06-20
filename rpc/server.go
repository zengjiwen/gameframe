package rpc

import (
	"github.com/zengjiwen/gameframe/env"
	"github.com/zengjiwen/gameframe/rpc/protos"
	"google.golang.org/grpc"
	"net"
)

var _server *grpc.Server

func StartServer(service protos.RPCServer) error {
	listener, err := net.Listen("tcp", env.ServiceAddr)
	if err != nil {
		return err
	}

	_server = grpc.NewServer()
	protos.RegisterRPCServer(_server, service)
	// todo add wait group
	go _server.Serve(listener)
	return nil
}

func StopServer() {
	_server.GracefulStop()
}

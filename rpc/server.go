package rpc

import (
	"github.com/zengjiwen/gameframe/rpc/protos"
	"google.golang.org/grpc"
	"net"
)

var _server = &server{}

type server struct {
	grpcServer *grpc.Server
}

func StartServer(addr string, service protos.RPCServer) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	_server.grpcServer = grpc.NewServer()
	protos.RegisterRPCServer(_server.grpcServer, service)
	// todo add wait group
	go _server.grpcServer.Serve(listener)
	return nil
}

func StopServer() {
	_server.grpcServer.GracefulStop()
}

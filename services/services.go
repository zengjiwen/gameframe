package services

import "google.golang.org/grpc"

type service struct {
	server *grpc.Server
}

func NewService() *service {
	return &service{}
}

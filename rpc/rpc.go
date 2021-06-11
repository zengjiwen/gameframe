package rpc

import "context"

type RPCClient interface {
	Send(route string, data []byte) error
	SendPush(userID string, frontendSv *Server, push *protos.Push) error
	SendKick(userID string, serverType string, kick *protos.KickMsg) error
	BroadcastSessionBind(uid string) error
	Call(ctx context.Context, rpcType protos.RPCType, route *route.Route, session *session.Session, msg *message.Message, server *Server) (*protos.Response, error)
}

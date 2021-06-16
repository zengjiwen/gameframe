package peers

import (
	"context"
	"github.com/zengjiwen/gameframe/codecs"
	"github.com/zengjiwen/gameframe/rpc/proto"
	"google.golang.org/grpc"
	"sync"
)

var RemotesMu sync.RWMutex
var Remotes = make(map[string]*remote)

type remote struct {
	conn *grpc.ClientConn
	rpc  proto.RPCClient
}

func (r *remote) Send(route string, arg interface{}) {

}

func (r *remote) Call(m *codecs.Message) (*proto.Respond, error) {
	request := &proto.Request{
		Route:   m.Route,
		Payload: m.Payload,
	}

	respond, err := r.rpc.Call(context.Background(), request)
	return respond, err
}

package env

import (
	"github.com/zengjiwen/gameframe/codec"
	"github.com/zengjiwen/gameframe/marshaler"
	"github.com/zengjiwen/gameframe/servicediscovery"
	"time"
)

var (
	ServerID    string
	ServerType  string
	ServiceAddr string
	DieChan     = make(chan struct{})
	Codec       = codec.NewPlain()
	Marshaler   = marshaler.NewProtobuf()
	SD          servicediscovery.ServiceDiscovery
	RPCTimeout  = 5 * time.Second
)

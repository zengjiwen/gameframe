package env

import (
	"github.com/zengjiwen/gameframe/codec"
	"github.com/zengjiwen/gameframe/marshaler"
	"github.com/zengjiwen/gameframe/servicediscovery"
)

var (
	ServerType  string
	ServiceAddr string
	Codec       = codec.NewPlain()
	Marshaler   = marshaler.NewProtobuf()
	SD          = servicediscovery.NewEtcd()
)

package env

import (
	"github.com/zengjiwen/gameframe/codecs"
	"github.com/zengjiwen/gameframe/marshalers"
	"github.com/zengjiwen/gameframe/servicediscovery"
)

var (
	ServerType  string
	ServiceAddr string
	Codec       = codecs.NewPlain()
	Marshaler   = marshalers.NewProtobuf()
	SD          = servicediscovery.NewEtcd()
)

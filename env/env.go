package env

import (
	"github.com/zengjiwen/gameframe/codecs"
	"github.com/zengjiwen/gameframe/marshalers"
)

var (
	ServerType     string
	Addr           string
	IsFrontend     bool
	ConcurrentMode                      = "actor"
	Codec          codecs.Codec         = codecs.NewPlain()
	Marshaler      marshalers.Marshaler = marshalers.NewProtobuf()
)

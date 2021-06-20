package env

import (
	"github.com/zengjiwen/gameframe/codecs"
	"github.com/zengjiwen/gameframe/marshalers"
)

var (
	ServerType  string
	ServiceAddr string
	Codec       codecs.Codec
	Marshaler   marshalers.Marshaler
)

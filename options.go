package gameframe

import (
	"github.com/zengjiwen/gameframe/codecs"
	"github.com/zengjiwen/gameframe/marshalers"
	"github.com/zengjiwen/gameframe/servicediscovery"
)

var _opts = options{
	concurrentMode: "actor",
}

type options struct {
	concurrentMode string
	codec          codecs.Codec
	marshaler      marshalers.Marshaler
	clientAddr     string
	sd             servicediscovery.ServiceDiscovery
}

func WithConcurrentMode(concurrentMode string) func(*options) {
	return func(opts *options) {
		opts.concurrentMode = concurrentMode
	}
}

func WithCodec(codec codecs.Codec) func(*options) {
	return func(opts *options) {
		opts.codec = codec
	}
}

func WithMarshaler(marshaler marshalers.Marshaler) func(*options) {
	return func(opts *options) {
		opts.marshaler = marshaler
	}
}

func WithClientAddr(clientAddr string) func(*options) {
	return func(opts *options) {
		opts.clientAddr = clientAddr
	}
}

func WithServiceDiscovery(sd servicediscovery.ServiceDiscovery) func(*options) {
	return func(opts *options) {
		opts.sd = sd
	}
}

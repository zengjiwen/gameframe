package gameframe

import (
	"github.com/zengjiwen/gameframe/codec"
	"github.com/zengjiwen/gameframe/marshaler"
	"github.com/zengjiwen/gameframe/servicediscovery"
)

var _opts = options{
	concurrentMode: "actor",
}

type options struct {
	serviceAddr    string
	concurrentMode string
	codec          codec.Codec
	marshaler      marshaler.Marshaler
	clientAddr     string
	sd             servicediscovery.ServiceDiscovery
	sdAddr         string
	standalone     bool
}

func WithServiceAddr(addr string) func(*options) {
	return func(opts *options) {
		opts.serviceAddr = addr
	}
}

func WithConcurrentMode(concurrentMode string) func(*options) {
	return func(opts *options) {
		opts.concurrentMode = concurrentMode
	}
}

func WithCodec(codec codec.Codec) func(*options) {
	return func(opts *options) {
		opts.codec = codec
	}
}

func WithMarshaler(marshaler marshaler.Marshaler) func(*options) {
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

func WithSDAddr(sdAddr string) func(*options) {
	return func(opts *options) {
		opts.sdAddr = sdAddr
	}
}

func WithStandalone() func(*options) {
	return func(opts *options) {
		opts.standalone = true
	}
}

package gameframe

import (
	"github.com/zengjiwen/gameframe/codecs"
	"github.com/zengjiwen/gameframe/marshalers"
)

type options struct {
	concurrentMode string
	codec          codecs.Codec
	marshaler      marshalers.Marshaler
}

var _opts options

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

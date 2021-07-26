package gameframe

type options struct {
	serviceAddr    string
	concurrentMode string
	clientAddr     string
	sdAddr         string
	standalone     bool
}

type Option interface {
	apply(opts *options)
}

type funcOption struct {
	f func(opts *options)
}

func (fo *funcOption) apply(opt *options) {
	fo.f(opt)
}

func newFuncOption(f func(opts *options)) *funcOption {
	return &funcOption{f: f}
}

func WithServiceAddr(addr string) Option {
	return newFuncOption(func(opts *options) {
		opts.serviceAddr = addr
	})
}

func WithConcurrentMode(mode string) Option {
	return newFuncOption(func(opts *options) {
		opts.concurrentMode = mode
	})
}

func WithClientAddr(addr string) Option {
	return newFuncOption(func(opts *options) {
		opts.clientAddr = addr
	})
}

func WithSDAddr(addr string) Option {
	return newFuncOption(func(opts *options) {
		opts.sdAddr = addr
	})
}

func WithStandalone() Option {
	return newFuncOption(func(opts *options) {
		opts.standalone = true
	})
}

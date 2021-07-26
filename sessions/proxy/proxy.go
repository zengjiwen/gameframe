package proxy

type Proxy interface {
	Send(route string, payload []byte) error
}

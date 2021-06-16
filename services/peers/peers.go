package peers

type Peer interface {
	Send(route string, arg interface{}) error
}

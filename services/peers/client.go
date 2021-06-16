package peers

import (
	"github.com/zengjiwen/gameframe/codecs"
	"github.com/zengjiwen/gameframe/env"
	"github.com/zengjiwen/gamenet"
)

type client struct {
	conn gamenet.Conn
}

func NewClient(conn gamenet.Conn) Peer {
	return &client{conn: conn}
}

func (c client) Send(route string, arg interface{}) error {
	payload, err := env.Marshaler.Marshal(arg)
	if err != nil {
		return err
	}

	m := codecs.NewMessage(route, payload)
	data, err := env.Codec.Encode(m)
	if err != nil {
		return err
	}

	c.conn.Send(data)
	return nil
}

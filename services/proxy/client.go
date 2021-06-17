package proxy

import (
	"github.com/zengjiwen/gameframe/codecs"
	"github.com/zengjiwen/gameframe/env"
	"github.com/zengjiwen/gamenet"
)

type client struct {
	conn gamenet.Conn
}

func NewClient(conn gamenet.Conn) Proxy {
	return &client{conn: conn}
}

func (c *client) Send(route string, payload []byte) error {
	m := codecs.NewMessage(route, payload)
	data, err := env.Codec.Encode(m)
	if err != nil {
		return err
	}

	c.conn.Send(data)
	return nil
}

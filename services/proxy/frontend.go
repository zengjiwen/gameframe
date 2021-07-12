package proxy

import (
	"github.com/zengjiwen/gameframe/codec"
	"github.com/zengjiwen/gameframe/env"
	"github.com/zengjiwen/gamenet"
)

type frontend struct {
	conn gamenet.Conn
}

func NewFrontend(conn gamenet.Conn) Proxy {
	return &frontend{conn: conn}
}

func (f *frontend) Send(route string, payload []byte) error {
	m := codec.NewMessage(route, payload)
	data, err := env.Codec.Encode(m)
	if err != nil {
		return err
	}

	f.conn.Send(data)
	return nil
}

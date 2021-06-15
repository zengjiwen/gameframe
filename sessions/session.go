package sessions

import (
	"github.com/zengjiwen/gameframe/codecs"
	"github.com/zengjiwen/gameframe/env"
	"github.com/zengjiwen/gamenet"
	"sync"
	"sync/atomic"
)

var (
	_sessionId             uint64
	_sessionMu             sync.RWMutex
	_sessions              = make(map[uint64]*Session)
	_sessionClosedCallback func(*Session)
)

type Session struct {
	ID             uint64
	uid            int
	conn           gamenet.Conn
	context        map[string]interface{}
	closedCallback func()
}

func New(c gamenet.Conn) *Session {
	s := &Session{
		ID:   genSessionId(),
		conn: c,
	}

	_sessionMu.Lock()
	_sessions[s.ID] = s
	_sessionMu.Unlock()

	return s
}

func (s *Session) Send(route string, arg interface{}) error {
	payload, err := env.Marshaler.Marshal(arg)
	if err != nil {
		return err
	}

	m := codecs.NewMessage(route, payload)
	data, err := env.Codec.Encode(m)
	if err != nil {
		return err
	}

	s.conn.Send(data)
	return nil
}

func (s *Session) OnClosed() {
	_sessionClosedCallback(s)
	s.closedCallback()

	_sessionMu.Lock()
	delete(_sessions, s.ID)
	_sessionMu.Unlock()
}

func genSessionId() uint64 {
	return atomic.AddUint64(&_sessionId, 1)
}

func (s *Session) SetClosedCallback(cb func()) {
	s.closedCallback = cb
}

func SetClosedCallback(cb func(*Session)) {
	_sessionClosedCallback = cb
}

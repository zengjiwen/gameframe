package sessions

import (
	peers2 "github.com/zengjiwen/gameframe/services/peers"
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
	peer           peers2.Peer
	context        map[string]interface{}
	closedCallback func()
	Route2ServerId map[string]string
}

func New(peer peers2.Peer) *Session {
	s := &Session{
		ID:             genSessionId(),
		peer:           peer,
		context:        make(map[string]interface{}),
		Route2ServerId: make(map[string]string),
	}

	_sessionMu.Lock()
	_sessions[s.ID] = s
	_sessionMu.Unlock()

	return s
}

func (s *Session) Send(route string, arg interface{}) error {
	return s.peer.Send(route, arg)
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

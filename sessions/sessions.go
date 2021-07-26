package sessions

import (
	"github.com/zengjiwen/gameframe/sessions/proxy"
	"sync"
	"sync/atomic"
)

var (
	_sessionId             int64
	_sessionMu             sync.RWMutex
	_sessions              = make(map[int64]*Session)
	_sessionClosedCallback func(*Session)
)

type Session struct {
	ID             int64
	uid            int
	proxy          proxy.Proxy
	context        map[string]interface{}
	closedCallback func()
	Route2ServerId map[string]string
}

func New(p proxy.Proxy) *Session {
	s := &Session{
		ID:             genSessionId(),
		proxy:          p,
		context:        make(map[string]interface{}),
		Route2ServerId: make(map[string]string),
	}

	_sessionMu.Lock()
	_sessions[s.ID] = s
	_sessionMu.Unlock()
	return s
}

func (s *Session) Send(route string, payload []byte) error {
	return s.proxy.Send(route, payload)
}

func (s *Session) OnClosed() {
	_sessionClosedCallback(s)
	s.closedCallback()

	_sessionMu.Lock()
	delete(_sessions, s.ID)
	_sessionMu.Unlock()
}

func genSessionId() int64 {
	return atomic.AddInt64(&_sessionId, 1)
}

func (s *Session) SetClosedCallback(cb func()) {
	s.closedCallback = cb
}

func SetClosedCallback(cb func(*Session)) {
	_sessionClosedCallback = cb
}

func SessionByID(sessionID int64) *Session {
	_sessionMu.RLock()
	defer _sessionMu.RUnlock()

	session, ok := _sessions[sessionID]
	if !ok {
		return nil
	}
	return session
}

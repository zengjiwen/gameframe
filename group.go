package gameframe

import (
	"errors"
	"github.com/zengjiwen/gameframe/codecs"
	"github.com/zengjiwen/gameframe/sessions"
	"sync"
)

var SessionExistErr = errors.New("session exist err")

type Group struct {
	name    string
	mu      sync.RWMutex
	members map[uint64]*sessions.Session
}

func NewGroup(n string) *Group {
	return &Group{
		name:    n,
		members: make(map[uint64]*sessions.Session),
	}
}

func (g *Group) Broadcast(route string, arg interface{}) error {
	payload, err := _app.marshaler.Marshal(arg)
	if err != nil {
		return err
	}

	m := codecs.NewMessage(route, payload)
	data, err := _app.codec.Encode(m)
	if err != nil {
		return err
	}

	g.mu.RLock()
	for _, mem := range g.members {
		mem.Send(data)
	}
	g.mu.RUnlock()

	return nil
}

func (g *Group) Multicast(route string, arg interface{}, filter func(*sessions.Session) bool) error {
	payload, err := _app.marshaler.Marshal(arg)
	if err != nil {
		return err
	}

	m := codecs.NewMessage(route, payload)
	data, err := _app.codec.Encode(m)
	if err != nil {
		return err
	}

	g.mu.RLock()
	for _, mem := range g.members {
		if !filter(mem) {
			continue
		}

		mem.Send(data)
	}
	g.mu.RUnlock()

	return nil
}

func (g *Group) Join(s *sessions.Session) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	_, ok := g.members[s.ID]
	if ok {
		return SessionExistErr
	}

	g.members[s.ID] = s
	return nil
}

func (g *Group) Leave(s *sessions.Session) {
	g.mu.Lock()
	delete(g.members, s.ID)
	g.mu.Unlock()
}

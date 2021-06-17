package gameframe

import (
	"errors"
	"github.com/zengjiwen/gameframe/env"
	"github.com/zengjiwen/gameframe/sessions"
	"sync"
)

var SessionExistErr = errors.New("session exist err")

type Group struct {
	name    string
	mu      sync.RWMutex
	members map[int64]*sessions.Session
}

func NewGroup(n string) *Group {
	return &Group{
		name:    n,
		members: make(map[int64]*sessions.Session),
	}
}

func (g *Group) Broadcast(route string, arg interface{}) {
	payload, err := env.Marshaler.Marshal(arg)
	if err != nil {
		return
	}

	g.mu.RLock()
	for _, mem := range g.members {
		mem.Send(route, payload)
	}
	g.mu.RUnlock()
}

func (g *Group) Multicast(route string, arg interface{}, filter func(*sessions.Session) bool) {
	payload, err := env.Marshaler.Marshal(arg)
	if err != nil {
		return
	}

	g.mu.RLock()
	for _, mem := range g.members {
		if !filter(mem) {
			continue
		}

		mem.Send(route, payload)
	}
	g.mu.RUnlock()
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

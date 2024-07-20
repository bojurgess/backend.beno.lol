package broker

import (
	"sync"
	"time"

	"github.com/bojurgess/backend.beno.lol/internal/config"
	"github.com/bojurgess/backend.beno.lol/internal/database"
	"github.com/bojurgess/backend.beno.lol/internal/spotify"
)

type Broker struct {
	mu         sync.Mutex
	generators map[string]*Generator
	config     *config.Config
}

type Generator struct {
	ch       chan spotify.NowPlayingResponse
	refCount int
	stopCh   chan struct{}
}

func NewBroker(c *config.Config) *Broker {
	return &Broker{
		generators: make(map[string]*Generator),
		config:     c,
	}
}

func (b *Broker) Subscribe(u *database.User) chan spotify.NowPlayingResponse {
	b.mu.Lock()
	defer b.mu.Unlock()

	if gen, exists := b.generators[u.ID]; exists {
		gen.refCount++
		return gen.ch
	}

	gen := &Generator{
		ch:       make(chan spotify.NowPlayingResponse),
		refCount: 1,
		stopCh:   make(chan struct{}),
	}
	b.generators[u.ID] = gen

	go b.GetNowPlaying(u, gen)
	return gen.ch
}

func (b *Broker) Unsubscribe(id string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if gen, exists := b.generators[id]; exists {
		gen.refCount--
		if gen.refCount == 0 {
			close(gen.stopCh)
			delete(b.generators, id)
		}
	}
}

func (b *Broker) GetNowPlaying(u *database.User, gen *Generator) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-gen.stopCh:
			close(gen.ch)
			return
		case <-ticker.C:
			np := b.WrappedNP(u)
			gen.ch <- np
		}
	}
}

func (b *Broker) WrappedNP(u *database.User) spotify.NowPlayingResponse {
	if time.Now().After(u.ExpiresAt) {
		err := spotify.RefreshAccessToken(&u.Tokens, *b.config.Env)
		if err != nil {
			panic(err)
		}
	}

	np, err := spotify.GetNowPlaying(u.Tokens)
	if err != nil {
		panic(err)
	}

	return *np
}

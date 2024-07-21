package broker

import (
	"fmt"
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
	clients  map[chan spotify.NowPlayingResponse]struct{}
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

	gen, exists := b.generators[u.ID]
	if !exists {
		println("Creating new generator")
		gen = &Generator{
			ch:       make(chan spotify.NowPlayingResponse),
			refCount: 0,
			stopCh:   make(chan struct{}),
			clients:  make(map[chan spotify.NowPlayingResponse]struct{}),
		}
		b.generators[u.ID] = gen
		go b.GetNowPlaying(u, gen)
	}

	clientCh := make(chan spotify.NowPlayingResponse)
	gen.clients[clientCh] = struct{}{}
	gen.refCount++

	return clientCh
}

func (b *Broker) Unsubscribe(id string, clientCh chan spotify.NowPlayingResponse) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if gen, exists := b.generators[id]; exists {
		delete(gen.clients, clientCh)
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
			b.mu.Lock()
			for clientCh := range gen.clients {
				select {
				case clientCh <- np:
				default:
				}
			}
			b.mu.Unlock()
		}
	}
}

func (b *Broker) WrappedNP(u *database.User) spotify.NowPlayingResponse {
	if time.Now().After(u.ExpiresAt) {
		err := spotify.RefreshAccessToken(&u.Tokens, *b.config.Env)
		if err != nil {
			fmt.Println(err)
		}
	}

	np, err := spotify.GetNowPlaying(u.Tokens)
	if err != nil {
		fmt.Println(err)
	}

	if np == nil {
		np = &spotify.NowPlayingResponse{}
	}

	return *np
}

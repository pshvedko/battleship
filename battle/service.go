package battle

import (
	"github.com/google/uuid"
	"sync"
)

type Battle struct {
	mutex sync.Mutex
	games map[uuid.UUID]*game
	sizes []int
}

func NewBattle(height, width int, sizes ...int) *Battle {
	return &Battle{
		mutex: sync.Mutex{},
		games: map[uuid.UUID]*game{},
		sizes: sizes,
	}
}

func (b *Battle) Get(id uuid.UUID) *game {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	g, ok := b.games[id]
	if !ok {
		g = New(b.sizes...)
		b.games[id] = g
	}
	return g
}

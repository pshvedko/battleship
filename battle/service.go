package battle

import (
	"github.com/google/uuid"
	"sync"
)

type Battle interface {
	Get(id uuid.UUID) *game
}

type battle struct {
	mutex sync.Mutex
	games map[uuid.UUID]*game
	sizes []int
}

func NewBattle(height, width int, sizes ...int) *battle {
	return &battle{
		mutex: sync.Mutex{},
		games: map[uuid.UUID]*game{},
		sizes: sizes,
	}
}

func (b *battle) Get(id uuid.UUID) *game {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	g, ok := b.games[id]
	if !ok {
		g = New(b.sizes...)
		b.games[id] = g
	}
	return g
}

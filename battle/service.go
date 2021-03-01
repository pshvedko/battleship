package battle

import (
	"github.com/google/uuid"
	"sync"
)

type Battle interface {
	Field(id uuid.UUID) []point
	Click(id uuid.UUID, x, y int) []point
}

type battle struct {
	mutex sync.Mutex
	games map[uuid.UUID]*game
	sizes []int
}

func NewBattle(sizes ...int) *battle {
	return &battle{
		mutex: sync.Mutex{},
		games: map[uuid.UUID]*game{},
		sizes: sizes,
	}
}

func (b *battle) get(id uuid.UUID) *game {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	g, ok := b.games[id]
	if !ok {
		g = &game{}
		g.initialize(b.sizes...)
		b.games[id] = g
	}
	return g
}

func (b *battle) Field(id uuid.UUID) []point {
	return b.get(id).Field()
}

func (b *battle) Click(id uuid.UUID, x, y int) []point {
	return b.get(id).Click(x, y)
}

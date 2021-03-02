package battle

import (
	"github.com/google/uuid"
	"sync"
)

type Battle interface {
	Begin(id uuid.UUID) []point
	Click(id uuid.UUID, x, y int) []point
	Reset(id uuid.UUID)
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

func (b *battle) remove(id uuid.UUID) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	delete(b.games, id)
}

func (b *battle) Begin(id uuid.UUID) []point {
	return b.get(id).Field()
}

func (b *battle) Reset(id uuid.UUID) {
	b.remove(id)
}

func (b *battle) Click(id uuid.UUID, x, y int) []point {
	return b.get(id).Click(x, y)
}

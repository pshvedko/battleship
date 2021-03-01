package battle

import (
	"github.com/google/uuid"
	"sync"
)

type Battle interface {
	Own(id uuid.UUID) []point
	Alien(id uuid.UUID) []point
	Shot(id uuid.UUID, x, y int) ([]point, int)
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
		g = New(b.sizes...)
		b.games[id] = g
	}
	return g
}

func (b *battle) Own(id uuid.UUID) []point {
	return b.get(id).Field(0)
}

func (b *battle) Alien(id uuid.UUID) []point {
	return b.get(id).Field(1)
}

func (b *battle) Shot(id uuid.UUID, x, y int) ([]point, int) {
	return b.get(id).Click(x, y)
}

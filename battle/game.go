package battle

import (
	"sync"
)

type game struct {
	sync.Mutex
	fields [2]field
}

func New(sizes ...int) *game {
	g := &game{}
	g.fields[0].initialize(sizes...)
	g.fields[1].initialize(sizes...)
	return g
}

func (g *game) Field0() *field {
	return &g.fields[0]
}

func (g *game) Field1() *field {
	return &g.fields[1]
}

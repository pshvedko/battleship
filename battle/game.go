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

func (g *game) Field(n int) (points []point) {
	g.Lock()
	defer g.Unlock()
	for x := range g.fields[n] {
		for y := range g.fields[n][x] {
			points = append(points, g.fields[n].point(x, y))
		}
	}
	return
}

func (g *game) Click(x int, y int) (points []point, class int) {
	g.Lock()
	defer g.Unlock()
	return nil, 0
}

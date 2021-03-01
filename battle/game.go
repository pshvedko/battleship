package battle

import (
	"sync"
)

type game struct {
	mutex  sync.Mutex
	fields [2]field
}

func (g *game) initialize(sizes ...int) {
	g.fields[0].initialize(sizes...)
	g.fields[1].initialize(sizes...)
}

func (g *game) Field(n int) (points []point) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	for x := range g.fields[n] {
		for y := range g.fields[n][x] {
			points = append(points, g.fields[n].point(x, y))
		}
	}
	return
}

func (g *game) Click(x int, y int) (points []point, class int) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	return nil, 0
}

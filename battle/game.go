package battle

import (
	"math/rand"
	"sync"
)

type shooter func() (int, int, bool)

func (f shooter) shoot() (int, int, bool) {
	return f()
}

type game struct {
	mutex  sync.Mutex
	fields [2]field
	shooter
}

func (g *game) initialize(sizes ...int) {
	g.fields[0].initialize(sizes...)
	g.fields[1].initialize(sizes...)
	g.shooter = g.randomShot
}

func (g *game) Field() (points []point) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	for n := range g.fields {
		for x := range g.fields[n] {
			for y := range g.fields[n][x] {
				if n == 1 && g.fields[n][x][y] < 2 {
					continue
				}
				points = append(points, g.fields[n].point(n, x, y))
			}
		}
	}
	return
}

func (g *game) Click(x int, y int) []point {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	points, hit := g.fields[1].shot(1, x, y)
	if !hit {
		points = append(points, g.answer()...)
	}
	return points
}

func (g *game) answer() (points []point) {
	defer func() {
		g.shooter = g.randomShot
	}()
	for {
		x, y, ok := g.shoot()
		if !ok {
			return
		}
		shots, hit := g.fields[0].shot(0, x, y)
		points = append(points, shots...)
		if !hit {
			return
		}
	}
}

func (g *game) randomShot() (x int, y int, ok bool) {
	a := g.fields[0].clean(0)
	if len(a) == 0 {
		return
	}
	p := a[rand.Int()%len(a)]
	ok = true
	x = p.X()
	y = p.Y()
	g.shooter = g.aimedShot
	return
}

func (g *game) aimedShot() (x int, y int, ok bool) {
	// TODO
	return
}

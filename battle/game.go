package battle

import (
	"sync"
)

type shooter func() (int, int, bool)

func (s shooter) shot() (int, int, bool) {
	return s()
}

type game struct {
	mutex  sync.Mutex
	fields [2]field
	hits   []point
	shooter
}

func (g *game) initialize(sizes ...int) {
	g.fields[0].initialize(sizes...)
	g.fields[1].initialize(sizes...)
	g.shooter = g.rightShot
}

func (g *game) next(x, y, n int, s shooter) (int, int, bool) {
	if g.fields[0].target(x, y) {
		return x, y, true
	}
	g.hits = g.hits[:n]
	g.shooter = s
	return g.shot()
}

func (g *game) add(shots ...point) {
	g.hits = append(g.hits, shots[len(shots)-1])
}

func (g *game) xy() (int, int) {
	return g.hits[len(g.hits)-1].XY()
}

func (g *game) Field() (points []point) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	for n := range &g.fields {
		for y := range &g.fields[n] {
			for x := range &g.fields[n][y] {
				if n == 1 && g.fields[n].raw(x, y) < 2 {
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
	for {
		x, y, ok := g.shot()
		if !ok {
			break
		}
		shots, hit := g.fields[0].shot(0, x, y)
		points = append(points, shots...)
		if !hit {
			break
		}
		g.add(shots...)
	}
	return
}

func (g *game) rightShot() (int, int, bool) {
	if len(g.hits) == 0 {
		p := g.fields[0].random(0)
		return p.X(), p.Y(), p.ok()
	}
	x, y := g.xy()
	x++
	return g.next(x, y, 1, g.leftShot)
}

func (g *game) leftShot() (int, int, bool) {
	x, y := g.xy()
	x--
	return g.next(x, y, 1, g.downShot)
}

func (g *game) downShot() (int, int, bool) {
	x, y := g.xy()
	y++
	return g.next(x, y, 1, g.upShot)
}

func (g *game) upShot() (int, int, bool) {
	x, y := g.xy()
	y--
	return g.next(x, y, 0, g.rightShot)
}

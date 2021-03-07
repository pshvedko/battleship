package battle

import (
	"log"
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
	g.shooter = g.randomShot
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

func (g *game) xy() (int, int) {
	return g.hits[len(g.hits)-1].XY()
}

func (g *game) hit(n int) {
	g.hits = g.hits[:n]
}

func (g *game) add(shots ...point) {
	g.hits = append(g.hits, shots[len(shots)-1])
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
	if len(g.hits) == 0 {
		g.shooter = g.randomShot
	}
	return
}

func (g *game) randomShot() (x int, y int, ok bool) {
	p := g.fields[0].random(0)

	log.Print(p.X(), p.Y(), p.ok())

	if !p.ok() {
		return
	}
	g.hit(0)
	g.shooter = g.rightShot
	return p.X(), p.Y(), p.ok()
}

func (g *game) rightShot() (x int, y int, ok bool) {
	x, y = g.xy()
	x++
	if ok = g.fields[0].target(x, y); ok {
		return
	}
	g.hit(1)
	g.shooter = g.leftShot
	return g.leftShot()
}

func (g *game) leftShot() (x int, y int, ok bool) {
	x, y = g.xy()
	x--
	if ok = g.fields[0].target(x, y); ok {
		return
	}
	g.hit(1)
	g.shooter = g.downShot
	return g.downShot()
}

func (g *game) downShot() (x int, y int, ok bool) {
	x, y = g.xy()
	y++
	if ok = g.fields[0].target(x, y); ok {
		return
	}
	g.hit(1)
	g.shooter = g.upShot
	return g.upShot()
}

func (g *game) upShot() (x int, y int, ok bool) {
	x, y = g.xy()
	y--
	if ok = g.fields[0].target(x, y); ok {
		return
	}
	g.hit(1)
	g.shooter = g.randomShot
	return g.randomShot()
}

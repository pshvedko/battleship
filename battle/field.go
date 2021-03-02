package battle

import (
	"math/rand"
)

const (
	fieldFree = iota
	fieldShip
	fieldMiss
	fieldShot
	fieldOpen
)

type field [10][10]int

func (f *field) initialize(sizes ...int) {
	for _, size := range sizes {
		f.add(size)
	}
}

func (f *field) add(size int) {
	var h = [4]int{0, 0, 1, -1}
	var w = [4]int{1, -1, 0, 0}
	for {
		x := rand.Int() % 10
		y := rand.Int() % 10
		z := rand.Int() % 4
		if f.try(x, y, h[z], w[z], size) {
			return
		}
	}
}

func (f *field) try(x int, y int, h int, w int, size int) bool {
	if size == 0 {
		return true
	} else if f.border(x, y) {
		return false
	} else if f.empty(x, y) {
		if f.try(x+h, y+w, h, w, size-1) {
			f.set(x, y, fieldShip)
			return true
		}
	}
	return false
}

func (f *field) empty(x, y int) bool {
	return f.zero(x, y) &&
		f.zero(x-1, y) &&
		f.zero(x+1, y) &&
		f.zero(x, y-1) &&
		f.zero(x, y+1) &&
		f.zero(x+1, y+1) &&
		f.zero(x-1, y+1) &&
		f.zero(x+1, y-1) &&
		f.zero(x-1, y-1)
}

func (f *field) zero(x, y int) bool {
	return fieldFree == f.get(x, y)
}

func (f *field) point(n, x, y int) point {
	return point(x*10*10*10 + y*10*10 + f[x][y]%fieldOpen*10 + n)
}

func (f *field) shot(n, x, y int) (points []point, hit bool) {
	if f.border(x, y) {
		return
	} else if f[x][y] < fieldMiss {
		f.inc(x, y, fieldMiss)
		if f[x][y] == fieldShot {
			points = append(points, f.around(n, x, y)...)
			hit = true
		}
	}
	points = append(points, f.point(n, x, y))
	return
}

func (f *field) around(n, x, y int) (points []point) {
	var w, h *int
	if false {
	} else if l := f.get(x-1, y); l == fieldShot {
		w = &x
		h = &y
		x--
	} else if r := f.get(x+1, y); r == fieldShot {
		w = &x
		h = &y
	} else if t := f.get(x, y-1); t == fieldShot {
		w = &y
		h = &x
		y--
	} else if b := f.get(x, y+1); b == fieldShot {
		w = &y
		h = &x
	} else if (l|r|t|b)&1 == fieldFree {
		w = &y
		h = &x
	}
	if w != nil && h != nil {
		*w--
		var a, b int
		for a = f.get(x, y); a == fieldShot; a = f.get(x, y) {
			*w--
		}
		c := *w
		*w++
		for b = f.get(x, y); b == fieldShot; b = f.get(x, y) {
			*h--
			points = append(points, f.update(fieldFree, fieldOpen, n, x, y)...)
			*h++
			*h++
			points = append(points, f.update(fieldFree, fieldOpen, n, x, y)...)
			*h--
			*w++
		}
		if (a|b)&1 == fieldFree {
			*h--
			points = append(points, f.update(fieldFree, fieldOpen, n, x, y)...)
			*h++
			points = append(points, f.update(fieldFree, fieldOpen, n, x, y)...)
			*h++
			points = append(points, f.update(fieldFree, fieldOpen, n, x, y)...)
			*w = c
			points = append(points, f.update(fieldFree, fieldOpen, n, x, y)...)
			*h--
			points = append(points, f.update(fieldFree, fieldOpen, n, x, y)...)
			*h--
			points = append(points, f.update(fieldFree, fieldOpen, n, x, y)...)
		}
	}
	return
}

func (f *field) get(x int, y int) int {
	if f.border(x, y) {
		return 0
	}
	return f[x][y] % fieldOpen
}

func (f *field) update(a, b, n, x, y int) (points []point) {
	if f.border(x, y) {
		return
	} else if f[x][y] != a {
		return
	}
	f.set(x, y, b)
	return append(points, f.point(n, x, y))
}

func (f *field) border(x int, y int) bool {
	return x < 0 || x >= len(f) || y < 0 || y >= len(f)
}

func (f *field) set(x int, y int, i int) {
	f[x][y] = i
}

func (f *field) inc(x int, y int, i int) {
	f[x][y] += i
}

func (f *field) clean(n int) (points []point) {
	for i := range f {
		for j := range f[i] {
			if f[i][j] < fieldMiss {
				points = append(points, f.point(n, i, j))
			}
		}
	}
	return
}

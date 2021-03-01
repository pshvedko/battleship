package battle

import (
	"math/rand"
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
			f[x][y] = 1
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
	return 0 == f.get(x, y)
}

func (f *field) point(n, x, y int) point {
	return point(x*10*10*10 + y*10*10 + f[x][y]%4*10 + n)
}

func (f *field) shot(n, x, y int) (points []point) {
	if f.border(x, y) {
		return
	} else if f[x][y] < 2 {
		f[x][y] += 2
		if f[x][y] == 3 {
			points = append(points, f.around(n, x, y)...)
		}
	}
	return append(points, f.point(n, x, y))
}

func (f *field) around(n, x, y int) (points []point) {
	var w, h *int
	if false {
	} else if l := f.get(x-1, y); l == 3 {
		w = &x
		h = &y
		x--
	} else if r := f.get(x+1, y); r == 3 {
		w = &x
		h = &y
	} else if t := f.get(x, y-1); t == 3 {
		w = &y
		h = &x
		y--
	} else if b := f.get(x, y+1); b == 3 {
		w = &y
		h = &x
	} else if (l|r|t|b)&1 == 0 {
		w = &y
		h = &x
	}
	if w != nil && h != nil {
		*w--
		var a, b int
		for a = f.get(x, y); a == 3; a = f.get(x, y) {
			*w--
		}
		c := *w
		*w++
		for b = f.get(x, y); b == 3; b = f.get(x, y) {
			*h--
			points = append(points, f.update(0, 4, n, x, y)...)
			*h++
			*h++
			points = append(points, f.update(0, 4, n, x, y)...)
			*h--
			*w++
		}
		if (a|b)&1 == 0 {
			*h--
			points = append(points, f.update(0, 4, n, x, y)...)
			*h++
			points = append(points, f.update(0, 4, n, x, y)...)
			*h++
			points = append(points, f.update(0, 4, n, x, y)...)
			*w = c
			points = append(points, f.update(0, 4, n, x, y)...)
			*h--
			points = append(points, f.update(0, 4, n, x, y)...)
			*h--
			points = append(points, f.update(0, 4, n, x, y)...)
		}
	}
	return
}

func (f *field) get(x int, y int) int {
	if f.border(x, y) {
		return 0
	}
	return f[x][y] % 4
}

func (f *field) update(a, b, n, x, y int) (points []point) {
	if f.border(x, y) {
		return
	} else if f[x][y] != a {
		return
	}
	f[x][y] = b
	return append(points, f.point(n, x, y))
}

func (f *field) border(x int, y int) bool {
	return x < 0 || x > 9 || y < 0 || y > 9
}

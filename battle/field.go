package battle

import "math/rand"

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
	} else if x < 0 || x > 9 || y < 0 || y > 9 {
		return false
	} else if f.around(x, y) {
		if f.try(x+h, y+w, h, w, size-1) {
			f[x][y] = 1
			return true
		}
	}
	return false
}

func (f *field) around(x, y int) bool {
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
	if x < 0 || x > 9 || y < 0 || y > 9 {
		return true
	}
	return f[x][y] == 0
}

func (f *field) point(x int, y int) point {
	return point(x*10*10 + y*10 + f[x][y])
}

func (f *field) boom(x int, y int) ([]point, int) {
	if x < 0 || x > 9 || y < 0 || y > 9 {
		return nil, 0
	}
	f[x][y] += 2
	return nil, f[x][y]
}

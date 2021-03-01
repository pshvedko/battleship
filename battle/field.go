package battle

import (
	"log"
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
	} else if x < 0 || x > 9 || y < 0 || y > 9 {
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

func (f *field) boom(n, x, y int) (points []point) {
	if x < 0 || x > 9 || y < 0 || y > 9 {
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
	a := [4]int{f.get(x+1, y), f.get(x-1, y), f.get(x, y+1), f.get(x, y-1)}
	if false {
	} else if a[0] == 3 {
		log.Print(1)
	} else if a[1] == 3 {
		log.Print(2)
	} else if a[2] == 3 {
		log.Print(3)
	} else if a[3] == 3 {
		log.Print(4)
	} else if (a[0]|a[1]|a[2]|a[3])&1 == 0 {
		log.Print(5)
		points = append(points, f.update(0, 4, n, x+1, y-0)...)
		points = append(points, f.update(0, 4, n, x-1, y-0)...)
		points = append(points, f.update(0, 4, n, x-0, y+1)...)
		points = append(points, f.update(0, 4, n, x-0, y-1)...)
		points = append(points, f.update(0, 4, n, x+1, y+1)...)
		points = append(points, f.update(0, 4, n, x-1, y-1)...)
		points = append(points, f.update(0, 4, n, x-1, y+1)...)
		points = append(points, f.update(0, 4, n, x+1, y-1)...)
	} else {
		log.Print(0)
	}
	return
}

func (f *field) get(x int, y int) int {
	if x < 0 || x > 9 || y < 0 || y > 9 {
		return 0
	}
	return f[x][y] % 4
}

func (f *field) update(a, b, n, x, y int) (points []point) {
	if x < 0 || x > 9 || y < 0 || y > 9 {
		return
	} else if f[x][y] != a {
		return
	}
	f[x][y] = b
	return append(points, f.point(n, x, y))
}

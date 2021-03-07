package battle

type point int

func (p point) X() int {
	return int(p) / 10 / 10 / 10
}

func (p point) Y() int {
	return int(p) / 10 / 10 % 10
}

func (p point) C() int {
	return int(p) / 10 % 10
}

func (p point) F() int {
	return int(p) % 10
}

func (p point) XY() (int, int) {
	return p.X(), p.Y()
}

func (p point) ok() bool {
	return p >= 0
}

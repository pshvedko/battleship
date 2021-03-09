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

func (p point) XYZ() (int, int, bool) {
	if p < 0 {
		return 0, 0, false
	}
	return p.X(), p.Y(), true
}

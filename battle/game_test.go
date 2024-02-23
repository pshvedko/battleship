package battle

import (
	"testing"
)

func benchmark(b *testing.B, a int) {
	var n, c int
	for i := 0; i < b.N; i++ {
		var g game
		g.initialize(a, 4, 3, 3, 2, 2, 2)
		for g.alive() {
			p := g.answer()
			n += len(p)
			c++
		}
	}
	b.ReportMetric(float64(n)/float64(b.N), "shots")
	b.ReportMetric(float64(c)/float64(b.N), "moves")
}

func BenchmarkGame(b *testing.B) {
	for i, n := range []string{"Random", "Weight"} {
		b.Run(n, func(b *testing.B) {
			benchmark(b, i)
		})
	}
}

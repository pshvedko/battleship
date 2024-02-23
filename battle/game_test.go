package battle

import (
	"testing"
)

func benchmark(b *testing.B, v int) {
	var n, c int
	for i := 0; i < b.N; i++ {
		var g game
		g.initialize(v, 4, 3, 3, 2, 2, 2)
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
	b.Run("0", func(b *testing.B) {
		benchmark(b, 0)
	})
	b.Run("1", func(b *testing.B) {
		benchmark(b, 1)
	})
}

package main

import (
	"fmt"
	"testing"
)

// makeScenario builds a Voronoi-painted grid of the given radius
// and returns the grid, its adjacency matrix, and a start cell in region 1.
func makeScenario(radius int) (*Grid, *AdjMatrix, Hex) {
	g := NewGrid(radius)
	seeds := []Hex{
		{Q: 0, R: 0, S: 0},
		{Q: radius / 2, R: -radius / 2, S: 0},
		{Q: -radius / 2, R: radius / 2, S: 0},
		{Q: 0, R: radius / 2, S: -radius / 2},
	}
	PaintVoronoi(g, seeds)
	m := BuildAdjMatrix(g)
	start := Hex{Q: 0, R: 0, S: 0}
	return g, m, start
}

// BenchmarkBadFill measures the O(N²) recursive approach.
// Larger grids show disproportionate slowdown compared to GoodFill.
func BenchmarkBadFill(b *testing.B) {
	for _, radius := range []int{5, 10, 15, 20} {
		b.Run(fmt.Sprintf("radius_%02d", radius), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				g, m, start := makeScenario(radius)
				target := g.Cells[start]
				b.StartTimer()

				BadFill(g, m, start, target, FilledColor)
			}
		})
	}
}

// BenchmarkGoodFill measures the O(N) iterative BFS approach.
// Run on larger radii too — it handles them without issue.
func BenchmarkGoodFill(b *testing.B) {
	for _, radius := range []int{5, 10, 15, 20, 40, 80} {
		b.Run(fmt.Sprintf("radius_%02d", radius), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				g, _, start := makeScenario(radius)
				target := g.Cells[start]
				b.StartTimer()

				GoodFill(g, start, target, FilledColor)
			}
		})
	}
}

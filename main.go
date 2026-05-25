package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	os.MkdirAll("output", 0755)

	// ── Grid setup ────────────────────────────────────────────────────────────

	const radius = 12
	const cellSize = 24.0

	// Five seeds spread across the grid to create distinct Voronoi regions.
	seeds := []Hex{
		{Q: 0, R: -8, S: 8},
		{Q: 7, R: 1, S: -8},
		{Q: -7, R: 7, S: 0},
		{Q: 4, R: 4, S: -8},
		{Q: -3, R: -5, S: 8},
	}

	base := NewGrid(radius)
	PaintVoronoi(base, seeds)

	startHex := Hex{Q: 0, R: 0, S: 0}
	targetColor := base.Cells[startHex]

	n := len(base.Cells)
	fmt.Printf("Hex grid  radius=%d  cells=%d\n\n", radius, n)

	// ── Render: before ────────────────────────────────────────────────────────

	if err := Render(base, cellSize, "output/before.png"); err != nil {
		fmt.Println("render error:", err)
		return
	}
	fmt.Println("Saved → output/before.png")

	// ── Bad fill ──────────────────────────────────────────────────────────────

	gBad := base.Clone()
	mBad := BuildAdjMatrix(gBad)

	t0 := time.Now()
	countBad := BadFill(gBad, mBad, startHex, targetColor, FilledColor)
	durBad := time.Since(t0)

	if err := Render(gBad, cellSize, "output/after_bad.png"); err != nil {
		fmt.Println("render error:", err)
		return
	}
	fmt.Println("Saved → output/after_bad.png")

	// ── Good fill ─────────────────────────────────────────────────────────────

	gGood := base.Clone()

	t0 = time.Now()
	countGood := GoodFill(gGood, startHex, targetColor, FilledColor)
	durGood := time.Since(t0)

	if err := Render(gGood, cellSize, "output/after_good.png"); err != nil {
		fmt.Println("render error:", err)
		return
	}
	fmt.Println("Saved → output/after_good.png")
	fmt.Println()

	// ── Results table ─────────────────────────────────────────────────────────

	sep := strings.Repeat("─", 60)
	fmt.Println(sep)
	fmt.Printf("%-14s  %-10s  %-12s  %s\n", "Approach", "Cells", "Time", "Complexity")
	fmt.Println(sep)
	fmt.Printf("%-14s  %-10d  %-12v  O(N²) — matrix row scan + recursion\n",
		"BadFill", countBad, durBad)
	fmt.Printf("%-14s  %-10d  %-12v  O(N)  — coordinate math + BFS queue\n",
		"GoodFill", countGood, durGood)
	fmt.Println(sep)

	if durGood > 0 {
		ratio := float64(durBad) / float64(durGood)
		fmt.Printf("\nSpeedup: %.1f×\n", ratio)
	}

	// ── Multi-size timing table ───────────────────────────────────────────────

	fmt.Printf("\n%s\n", strings.Repeat("─", 60))
	fmt.Printf("%-10s  %-8s  %-14s  %-14s  %s\n",
		"Radius", "Cells", "BadFill", "GoodFill", "Ratio")
	fmt.Println(strings.Repeat("─", 60))

	for _, r := range []int{5, 8, 12, 18, 25, 40} {
		g := NewGrid(r)
		PaintVoronoi(g, seeds[:4])
		start := Hex{Q: 0, R: 0, S: 0}
		target := g.Cells[start]
		cells := len(g.Cells)

		// bad
		gb := g.Clone()
		mb := BuildAdjMatrix(gb)
		tb0 := time.Now()
		BadFill(gb, mb, start, target, FilledColor)
		tb := time.Since(tb0)

		// good
		gg := g.Clone()
		tg0 := time.Now()
		GoodFill(gg, start, target, FilledColor)
		tg := time.Since(tg0)

		ratio := "—"
		if tg > 0 {
			ratio = fmt.Sprintf("%.1f×", float64(tb)/float64(tg))
		}
		fmt.Printf("%-10d  %-8d  %-14v  %-14v  %s\n", r, cells, tb, tg, ratio)
	}
	fmt.Println(strings.Repeat("─", 60))
}

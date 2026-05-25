package main

import "sort"

// Color identifies a region on the grid.
// 0 is reserved for "empty"; values 1–5 are painted regions; 6 is the fill highlight.
type Color int

const FilledColor Color = 6

// Hex is a cell in cube coordinate space.
// Invariant: Q + R + S == 0 always.
//
// Cube coordinates give every hex exactly 6 neighbours at the same distance,
// which is the property that makes hex grids ideal for equal-cost pathfinding
// and flood fill.
type Hex struct{ Q, R, S int }

// directions lists the six unit steps in cube space.
// Adding any one of these to a Hex yields a direct neighbour.
var directions = [6]Hex{
	{1, -1, 0}, {1, 0, -1}, {0, 1, -1},
	{-1, 1, 0}, {-1, 0, 1}, {0, -1, 1},
}

// Neighbors returns all six adjacent hexes.
// Cost: O(1) — six additions, no lookup table needed.
func (h Hex) Neighbors() [6]Hex {
	var out [6]Hex
	for i, d := range directions {
		out[i] = Hex{h.Q + d.Q, h.R + d.R, h.S + d.S}
	}
	return out
}

// Distance returns the minimum number of steps between two hexes.
// In cube space this is simply the Chebyshev distance.
func Distance(a, b Hex) int {
	dq := abs(a.Q - b.Q)
	dr := abs(a.R - b.R)
	ds := abs(a.S - b.S)
	return max3(dq, dr, ds)
}

// Grid holds all cells and their colors.
type Grid struct {
	Cells  map[Hex]Color
	Radius int
}

// NewGrid creates a filled hexagonal grid with the given radius.
// Total cells: 3*radius² + 3*radius + 1
func NewGrid(radius int) *Grid {
	g := &Grid{Cells: make(map[Hex]Color), Radius: radius}
	for q := -radius; q <= radius; q++ {
		rMin := max2(-radius, -q-radius)
		rMax := min2(radius, -q+radius)
		for r := rMin; r <= rMax; r++ {
			g.Cells[Hex{q, r, -q - r}] = 0
		}
	}
	return g
}

// Clone returns a deep copy of the grid.
func (g *Grid) Clone() *Grid {
	c := &Grid{Cells: make(map[Hex]Color, len(g.Cells)), Radius: g.Radius}
	for h, col := range g.Cells {
		c.Cells[h] = col
	}
	return c
}

// SortedHexes returns all cells in a deterministic order (Q, R, S ascending).
// Used by BuildAdjMatrix to guarantee stable row/column indices.
func (g *Grid) SortedHexes() []Hex {
	hexes := make([]Hex, 0, len(g.Cells))
	for h := range g.Cells {
		hexes = append(hexes, h)
	}
	sort.Slice(hexes, func(i, j int) bool {
		if hexes[i].Q != hexes[j].Q {
			return hexes[i].Q < hexes[j].Q
		}
		if hexes[i].R != hexes[j].R {
			return hexes[i].R < hexes[j].R
		}
		return hexes[i].S < hexes[j].S
	})
	return hexes
}

// PaintVoronoi assigns each cell the color of its nearest seed,
// producing organic-looking regions without any randomness.
func PaintVoronoi(g *Grid, seeds []Hex) {
	for h := range g.Cells {
		best := Color(1)
		minD := int(^uint(0) >> 1)
		for i, s := range seeds {
			if d := Distance(h, s); d < minD {
				minD = d
				best = Color(i + 1)
			}
		}
		g.Cells[h] = best
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func max2(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min2(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max3(a, b, c int) int {
	if a >= b && a >= c {
		return a
	}
	if b >= c {
		return b
	}
	return c
}

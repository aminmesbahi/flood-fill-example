package main

// AdjMatrix is an N×N boolean adjacency matrix where N = number of cells.
//
// matrix[i][j] == true  ⟺  hex i and hex j are direct neighbours.
//
// This data structure is the root cause of BadFill's poor performance:
// finding the (at most 6) neighbours of cell i requires reading the
// entire row i — all N entries — because we have no index into which
// columns are set.
//
// Memory cost: O(N²). For a radius-30 grid (~2800 cells) this is ~7 MB
// of booleans. For radius-100 (~30 000 cells) it becomes ~900 MB.
type AdjMatrix struct {
	Data   [][]bool    // Data[i][j] = true if hex i and hex j are neighbours
	Index  []Hex       // Index[i] → which hex is at row/column i
	Lookup map[Hex]int // Lookup[hex] → row/column index of that hex
	N      int         // number of cells = len(Index)
}

// BuildAdjMatrix constructs the full N×N adjacency matrix for grid g.
//
// Time:  O(N²) — allocating and zeroing the matrix dominates.
// Space: O(N²)
func BuildAdjMatrix(g *Grid) *AdjMatrix {
	hexes := g.SortedHexes()
	n := len(hexes)

	lookup := make(map[Hex]int, n)
	for i, h := range hexes {
		lookup[h] = i
	}

	// Allocate the full N×N matrix.
	data := make([][]bool, n)
	for i := range data {
		data[i] = make([]bool, n) // zeroed by Go runtime
	}

	// For each cell, mark its (up to 6) neighbours in both directions.
	for i, h := range hexes {
		for _, nb := range h.Neighbors() {
			if j, ok := lookup[nb]; ok {
				data[i][j] = true
			}
		}
	}

	return &AdjMatrix{Data: data, Index: hexes, Lookup: lookup, N: n}
}

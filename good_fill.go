package main

// GoodFill is an iterative BFS flood fill that uses cube coordinate arithmetic.
//
// ── Why it is fast ────────────────────────────────────────────────────────────
//
// Improvement 1 — O(1) neighbour lookup
//
//	Each hex's six neighbours are found by adding one of six fixed direction
//	vectors to the cube coordinates. This is six additions and one map lookup
//	per neighbour — O(1) total per cell, regardless of grid size.
//	A fill touching K cells costs O(6K) = O(K) time.
//
// Improvement 2 — Iterative with an explicit queue
//
//	No call stack growth. The queue lives on the heap and grows only as large
//	as the frontier, which never exceeds the total number of cells.
//	Safe for grids of any size.
//
// Complexity: O(V) time  (V = cells in the grid, each visited at most once)
//
//	O(V) space (queue + visited set)
func GoodFill(g *Grid, start Hex, target, replacement Color) int {
	if target == replacement {
		return 0
	}
	if g.Cells[start] != target {
		return 0
	}

	// Pre-allocate queue and visited set to avoid repeated allocations.
	queue := make([]Hex, 0, len(g.Cells))
	queue = append(queue, start)

	visited := make(map[Hex]bool, len(g.Cells))
	visited[start] = true

	count := 0

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		g.Cells[current] = replacement
		count++

		// ── The fast part ─────────────────────────────────────────────────────
		// Compute all 6 neighbours via cube coordinate arithmetic — O(1).
		// No matrix row to scan; no index to look up.
		for _, nb := range current.Neighbors() {
			if !visited[nb] {
				if c, exists := g.Cells[nb]; exists && c == target {
					visited[nb] = true
					queue = append(queue, nb)
				}
			}
		}
		// ─────────────────────────────────────────────────────────────────────
	}

	return count
}

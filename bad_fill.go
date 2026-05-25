package main

// BadFill is a recursive DFS flood fill that uses the adjacency matrix.
//
// ── Why it is slow ────────────────────────────────────────────────────────────
//
// Problem 1 — O(N) neighbour lookup
//
//	A hex has at most 6 neighbours. But to find them we scan the entire row of
//	the N×N matrix: O(N) work per cell even though only 6 entries are ever true.
//	For a fill that touches K cells the total scan cost is K × N.
//	In the worst case K = N, giving O(N²) time.
//
// Problem 2 — Recursive call stack
//
//	Each filled cell adds one stack frame. For a large connected region this
//	means O(N) frames on the call stack simultaneously. In Go the goroutine
//	stack grows dynamically, so we avoid a hard crash — but in C/C++/Java with
//	fixed stack sizes this would be a stack overflow. Even in Go, the O(N)
//	stack allocations slow things down.
//
// Complexity: O(N²) time,  O(N) space (call stack depth)
func BadFill(g *Grid, m *AdjMatrix, start Hex, target, replacement Color) int {
	if target == replacement {
		return 0
	}
	idx, ok := m.Lookup[start]
	if !ok {
		return 0
	}
	if g.Cells[start] != target {
		return 0
	}

	count := 0
	badRecurse(g, m, idx, target, replacement, &count)
	return count
}

func badRecurse(g *Grid, m *AdjMatrix, idx int, target, replacement Color, count *int) {
	h := m.Index[idx]

	if g.Cells[h] != target {
		return
	}

	g.Cells[h] = replacement
	*count++

	// ── The expensive part ────────────────────────────────────────────────────
	// Scan every column in row idx to find neighbours.
	// This is O(N) even though at most 6 columns will be true.
	for j := 0; j < m.N; j++ {
		if m.Data[idx][j] {
			nb := m.Index[j]
			if g.Cells[nb] == target {
				badRecurse(g, m, j, target, replacement, count) // ← recursive call
			}
		}
	}
	// ─────────────────────────────────────────────────────────────────────────
}

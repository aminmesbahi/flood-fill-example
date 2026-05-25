package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"sort"
)

// palette maps Color values to RGBA.
// Index 0 is the fallback; indices 1–5 are the Voronoi regions;
// index 6 is the flood-fill highlight.
var palette = []color.RGBA{
	{R: 40, G: 40, B: 55, A: 255},    // 0: unused
	{R: 72, G: 120, B: 200, A: 255},  // 1: indigo blue
	{R: 60, G: 175, B: 120, A: 255},  // 2: sea green
	{R: 210, G: 120, B: 60, A: 255},  // 3: terracotta
	{R: 160, G: 100, B: 210, A: 255}, // 4: lavender
	{R: 210, G: 175, B: 55, A: 255},  // 5: golden
	{R: 240, G: 80, B: 80, A: 255},   // 6: fill highlight (red)
}

func colorFor(c Color) color.RGBA {
	if int(c) < len(palette) {
		return palette[c]
	}
	return color.RGBA{R: 180, G: 180, B: 180, A: 255}
}

// darker returns a slightly dimmed version of c for borders.
func darker(c color.RGBA) color.RGBA {
	dim := func(v uint8) uint8 {
		if v < 40 {
			return 0
		}
		return v - 40
	}
	return color.RGBA{R: dim(c.R), G: dim(c.G), B: dim(c.B), A: 255}
}

// ── Hex → pixel conversion (pointy-top orientation) ──────────────────────────
//
// Pointy-top means one vertex points straight up, flat edges are vertical.
// Formula (derived from cube → axial → cartesian):
//
//	px = cellSize × ( √3 × q  +  √3/2 × r )
//	py = cellSize × ( 3/2 × r )
func hexToPixel(h Hex, size, originX, originY float64) (float64, float64) {
	px := size * (math.Sqrt(3)*float64(h.Q) + math.Sqrt(3)/2*float64(h.R))
	py := size * (1.5 * float64(h.R))
	return originX + px, originY + py
}

// hexVertices returns the 6 corner points of a pointy-top hexagon
// centred at (cx, cy) with the given outer radius.
//
// Vertex k is at angle  60k + 30  degrees (30° offset for pointy-top).
func hexVertices(cx, cy, size float64) [6][2]float64 {
	var pts [6][2]float64
	for i := 0; i < 6; i++ {
		angle := math.Pi / 180 * (60*float64(i) + 30)
		pts[i] = [2]float64{
			cx + size*math.Cos(angle),
			cy + size*math.Sin(angle),
		}
	}
	return pts
}

// fillPolygon fills the interior of an arbitrary polygon using a scanline
// algorithm. For each horizontal scanline we collect edge-crossing x values
// and flood-fill between each pair.
func fillPolygon(img *image.RGBA, pts [6][2]float64, c color.RGBA) {
	// Bounding box in y.
	minY, maxY := pts[0][1], pts[0][1]
	for _, p := range pts {
		if p[1] < minY {
			minY = p[1]
		}
		if p[1] > maxY {
			maxY = p[1]
		}
	}

	n := len(pts)
	for y := int(math.Floor(minY)); y <= int(math.Ceil(maxY)); y++ {
		fy := float64(y) + 0.5 // sample at pixel centre

		// Find x coordinates where this scanline crosses each edge.
		var xs []float64
		for i := 0; i < n; i++ {
			x1, y1 := pts[i][0], pts[i][1]
			x2, y2 := pts[(i+1)%n][0], pts[(i+1)%n][1]
			if (y1 <= fy && y2 > fy) || (y2 <= fy && y1 > fy) {
				t := (fy - y1) / (y2 - y1)
				xs = append(xs, x1+t*(x2-x1))
			}
		}
		sort.Float64s(xs)

		// Fill between each pair of crossings.
		for i := 0; i+1 < len(xs); i += 2 {
			for x := int(math.Round(xs[i])); x <= int(math.Round(xs[i+1])); x++ {
				img.SetRGBA(x, y, c)
			}
		}
	}
}

// drawBorder traces the outline of a hexagon using Bresenham's line algorithm.
func drawBorder(img *image.RGBA, cx, cy, size float64, c color.RGBA) {
	pts := hexVertices(cx, cy, size*0.94) // slightly inset
	for i := 0; i < 6; i++ {
		bresenham(img,
			pts[i][0], pts[i][1],
			pts[(i+1)%6][0], pts[(i+1)%6][1],
			c,
		)
	}
}

func bresenham(img *image.RGBA, x1, y1, x2, y2 float64, c color.RGBA) {
	dx := math.Abs(x2 - x1)
	dy := math.Abs(y2 - y1)
	steps := int(math.Max(dx, dy)) + 1
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		x := int(math.Round(x1 + t*(x2-x1)))
		y := int(math.Round(y1 + t*(y2-y1)))
		img.SetRGBA(x, y, c)
	}
}

// Render saves the grid as a PNG to filename.
// cellSize is the outer radius of each hexagon in pixels.
func Render(g *Grid, cellSize float64, filename string) error {
	// Pass 1: find the pixel bounding box (no hard-coded formula needed).
	var minPX, minPY, maxPX, maxPY float64
	first := true
	for h := range g.Cells {
		px, py := hexToPixel(h, cellSize, 0, 0)
		if first {
			minPX, minPY, maxPX, maxPY = px, py, px, py
			first = false
		} else {
			if px < minPX {
				minPX = px
			}
			if py < minPY {
				minPY = py
			}
			if px > maxPX {
				maxPX = px
			}
			if py > maxPY {
				maxPY = py
			}
		}
	}

	pad := cellSize * 1.2
	originX := -minPX + pad + cellSize
	originY := -minPY + pad + cellSize
	imgW := int(math.Ceil(maxPX - minPX + 2*pad + 2*cellSize))
	imgH := int(math.Ceil(maxPY - minPY + 2*pad + 2*cellSize))

	// Pass 2: draw.
	img := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	// Dark background.
	draw.Draw(img, img.Bounds(),
		&image.Uniform{color.RGBA{R: 18, G: 18, B: 28, A: 255}},
		image.Point{}, draw.Src)

	for hex, col := range g.Cells {
		cx, cy := hexToPixel(hex, cellSize, originX, originY)
		fill := colorFor(col)
		// Fill interior at 93% size so a thin gap shows between cells.
		pts := hexVertices(cx, cy, cellSize*0.93)
		fillPolygon(img, pts, fill)
		// Draw a slightly darker border to emphasise the cell edges.
		drawBorder(img, cx, cy, cellSize, darker(fill))
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

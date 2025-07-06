package treesectiondraw

import (
	"math"

	treesection "github.com/formdenes/structures/tree-section"
	"github.com/tdewolff/canvas"
)

const col int = 5
const row int = 7

const ringNum int = 30
const ampl float64 = 10
const finalDiameter float64 = float64(ringNum) * ampl * 3

var trees []treesection.TreeSection

func Draw(ctx *canvas.Context, w float64, h float64) {
	for i := 0; i < col; i++ {
		for j := 0; j < row; j++ {
			trees = append(trees, *treesection.NewTreeSection(ringNum, ampl, 150))
		}
	}
	ctx.SetFillColor(canvas.Transparent)

	// maxNum := math.Max(float64(col), float64(row))
	// maxLen := math.Max(float64(w), float64(h))
	scaleW := w / (float64(col) * finalDiameter)
	scaleH := h / (float64(row) * finalDiameter)
	scale := math.Max(scaleH, scaleW)
	for i := 0; i < col; i++ {
		for j := 0; j < row; j++ {
			x := (float64(i) + 0.5) * finalDiameter * scale
			y := (float64(j) + 0.5) * finalDiameter * scale
			ctx.Push()
			ctx.Translate(x, y)
			ctx.Scale(scale, scale)
			ctx.SetStrokeColor(canvas.Black)
			ctx.SetStrokeWidth(0.5)
			if i < len(trees) {
				drawTreesection(ctx, trees[i])
			}
			ctx.Pop()
		}
	}
}

func drawTreesection(ctx *canvas.Context, treesection treesection.TreeSection) {

	for _, ring := range treesection.GetRings() {
		points := ring.GetPoints()
		ctx.MoveTo(points[0].X(), points[0].Y())
		for i := 1; i < len(points); i++ {
			ctx.LineTo(points[i].X(), points[i].Y())
		}
		ctx.Close()
		ctx.Stroke()
	}
}

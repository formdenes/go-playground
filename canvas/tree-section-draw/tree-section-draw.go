package treesectiondraw

import (
	"main/utils"
	"math"
	"sync"

	"github.com/KEINOS/go-noise"
	treesection "github.com/formdenes/structures/tree-section"
	"github.com/tdewolff/canvas"
)

const col int = 5
const row int = 7

const ringNum int = 25
const ampl float64 = 15
const minDistanceRatio float64 = 0.5
const smoothness float64 = 220
const finalDiameter float64 = float64(ringNum) * ampl * 2.7

type TreeSectionDraw struct {
	treesection *treesection.TreeSection
	i           int
	j           int
}

func Draw(ctx *canvas.Context, w float64, h float64, seed int64) {
	utils.DrawBackground(ctx, canvas.White)
	ctx.SetFillColor(canvas.Transparent)
	scaleW := w / (float64(col) * finalDiameter)
	scaleH := h / (float64(row) * finalDiameter)
	scale := math.Min(scaleH, scaleW)
	drawChan := make(chan TreeSectionDraw)
	seeds := utils.GenerateRandomSequence(seed, col)
	// seeds := utils.GenerateRandomSequence(100000, col)
	var wg sync.WaitGroup
	for i := 0; i < col; i++ {
		for j := 0; j < row; j++ {
			wg.Add(1)
			go func(i int, j int) {
				defer wg.Done()
				// seed := int64(k*10000) + int64(j*10)
				// seed := int64(100000)
				seed := int64(seeds[i])
				time := float64(j) / 3.5 // + float64(j*10)
				noiseType := noise.OpenSimplex
				// curAmpl := ampl
				curAmpl := ampl - float64(j)/4
				minDistance := curAmpl * minDistanceRatio
				smoothness := smoothness - float64(i*24)
				numberOfRings := ringNum - j
				treesection := treesection.NewTreeSection(treesection.Options{
					NumberOfRings: numberOfRings,
					Ampl:          curAmpl,
					MinDistance:   minDistance,
					Smoothness:    smoothness,
					Seed:          &seed,
					NoiseType:     &noiseType,
					Time:          &time,
				})
				drawChan <- TreeSectionDraw{treesection, i, j}
			}(i, j)
		}
	}

	go func() {
		wg.Wait()
		close(drawChan)
	}()

	for draw := range drawChan {
		drawTree(ctx, *draw.treesection, draw.i, draw.j, scale, scaleW, scaleH)
	}

	// maxNum := math.Max(float64(col), float64(row))
	// maxLen := math.Max(float64(w), float64(h))
	// for i := 0; i < col; i++ {
	// 	for j := 0; j < row; j++ {
	// 		x := (float64(i) + 0.5) * finalDiameter * scale
	// 		y := (float64(j) + 0.5) * finalDiameter * scale
	// 		ctx.Push()
	// 		ctx.Translate(x, y)
	// 		ctx.Scale(scale, scale)
	// 		ctx.SetStrokeColor(canvas.Black)
	// 		ctx.SetStrokeWidth(0.5)
	// 		if i < len(trees) {
	// 			drawTreesection(ctx, trees[i])
	// 		}
	// 		ctx.Pop()
	// 	}
	// }
}

func drawTree(ctx *canvas.Context, tree treesection.TreeSection, i int, j int, scale float64, scaleW float64, scaleH float64) {
	x := (float64(i) + 0.5) * finalDiameter * scaleW
	y := ctx.Height() - ((float64(j) + 0.5) * finalDiameter * scaleH)
	ctx.Push()
	ctx.Translate(x, y)
	ctx.Scale(scale, scale)
	ctx.SetStrokeColor(canvas.Black)
	ctx.SetStrokeWidth(0.5)
	drawTreesection(ctx, tree)
	ctx.Pop()
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

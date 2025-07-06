package main

import (
	"image/color"
	"math"

	treesection "github.com/formdenes/structures/tree-section"
	"github.com/go-p5/p5"
)

const w int = 1000
const h int = 1000
const col int = 1
const row int = 1

const ringNum int = 30
const ampl float64 = 10
const finalDiameter float64 = float64(ringNum) * ampl * 3

var trees []treesection.TreeSection

// var ng noise.Generator

func main() {
	for i := 0; i < col; i++ {
		for j := 0; j < row; j++ {
			trees = append(trees, *treesection.NewTreeSection(ringNum, ampl, 150))
		}
	}
	p5.Run(setup, draw)
	// p5.NoLoop()
}

func setup() {
	p5.Canvas(w, h)
	p5.Background(color.White)
	// p5.NoLoop()
}

func draw() {
	p5.Background(color.White)
	// p5.Translate(w/2, h/2)
	maxNum := math.Max(float64(col), float64(row))
	maxLen := math.Max(float64(w), float64(h))
	scale := maxLen / ((maxNum) * finalDiameter)
	for i := 0; i < col; i++ {
		for j := 0; j < row; j++ {
			x := (float64(i) + 0.5) * finalDiameter * scale
			y := (float64(j) + 0.5) * finalDiameter * scale
			p5.Push()
			p5.Translate(x, y)
			p5.Scale(scale, scale)
			trees[i*row+j].Draw()
			p5.Pop()
		}
	}
}

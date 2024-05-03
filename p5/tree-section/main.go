package main

import (
	"image/color"

	"github.com/go-p5/p5"
	"github.com/quartercastle/vector"
)

type vec = vector.Vector

const w = 1000
const h = 1000

var tree TreeSection

// var ng noise.Generator

func main() {
	tree = *NewTreeSection(30, 10, 300)
	p5.Run(setup, draw)
}

func setup() {
	p5.Canvas(w, h)
	p5.Background(color.White)
}

func draw() {
	p5.Background(color.White)
	p5.Translate(w/2, h/2)
	// ng, err := noise.New(noise.OpenSimplex, 1000)
	// if err != nil {
	// 	return
	// }
	// smoothness := 200.0

	// ampl := 50.0
	// res := 360
	// r := 300.0
	// var points [360]vec
	// for i := range res {
	// 	radians := float64(i) * (math.Pi / float64(180))
	// 	x := (r) * math.Cos(radians)
	// 	y := (r) * math.Sin(radians)
	// 	dr := ng.Eval64(x/smoothness, y/smoothness) * ampl
	// 	x = (r + dr) * math.Cos(radians)
	// 	y = (r + dr) * math.Sin(radians)
	// 	points[i] = vec{x, y}
	// }
	// for i := range 100 {
	// 	for j := range 100 {
	// 		c := ng.Eval64(float64(i)/smoothness, float64(j)/smoothness) * 255
	// 		p5.Stroke(color.Gray{uint8(c)})
	// 		p5.Fill(color.Gray{uint8(c)})
	// 		p5.Rect(float64(i)*10, float64(j)*10, 5, 5)

	// 	}
	// }
	// for i, point := range points {
	// 	if i > 0 {
	// 		prevPoint := points[i-1]
	// 		p5.Line(prevPoint.X(), prevPoint.Y(), point.X(), point.Y())
	// 		continue
	// 	}
	// 	lastPoint := points[len(points)-1]
	// 	p5.Line(lastPoint.X(), lastPoint.Y(), point.X(), point.Y())
	// }
	// p5.Circle(0, 0, 300)
	// p5.Circle(0, 0, 400)
	tree.Draw()
}

package main

import (
	"image/color"

	"github.com/go-p5/p5"
	"github.com/quartercastle/vector"
)

type vec = vector.Vector

const w = 1000
const h = 1000
const border = 100

var points []vec
var proc p5.Proc

func main() {
	proc = p5.Proc{}
	p5.Run(setup, draw)
}

func setup() {
	p5.Canvas(w, h)
	p5.Background(color.White)
	rows := 10
	columns := 10
	for i := range rows {
		for j := range columns {
			x := (w-border)/columns*i + border
			y := (w-border)/rows*j + border
			points = append(points, vec{float64(x), float64(y)})
		}
	}
}

func draw() {
	p5.Background(color.White)
	var p p5.Path
	for _, point := range points {
		// p5.Circle(point.X(), point.Y(), 2)
		p.Vertex(point.X(), point.Y())
	}
	p.End()
}

package main

import (
	"math"

	"github.com/KEINOS/go-noise"
	"github.com/go-p5/p5"
)

const res = 360

type Ring struct {
	points     [res]vec
	smoothness float64
	ampl       float64
}

func (r *Ring) draw() {
	if r == nil {
		return
	}
	for i, p := range r.points {
		var prevP vec
		if i > 0 {
			prevP = r.points[i-1]
		} else {
			prevP = r.points[res-1]
		}
		p5.Line(prevP.X(), prevP.Y(), p.X(), p.Y())
	}
}

func NewRing(prevRing *Ring, ampl float64, smoothness float64) *Ring {
	if prevRing == nil {
		var points [res]vec
		for i := range res {
			radians := float64(i) * (math.Pi / float64(180))
			x := (50 * math.Cos(radians))
			y := (50 * math.Sin(radians))
			points[i] = vec{x, y}
		}
		return &Ring{
			points:     points,
			smoothness: smoothness,
			ampl:       ampl,
		}
	}
	var points [res]vec
	ng, err := noise.New(noise.OpenSimplex, 1000)
	if err != nil {
		return nil
	}
	for i := range res {
		prevPoint := prevRing.points[i]
		dr := (ng.Eval64(prevPoint.X()/smoothness, prevPoint.Y()/smoothness) + 1) * ampl
		step := prevPoint.Unit().Scale(dr)
		point := prevPoint.Add(step)
		points[i] = point
	}
	return &Ring{
		points:     points,
		smoothness: smoothness,
		ampl:       prevRing.ampl,
	}
}

type TreeSection struct {
	rings []*Ring
}

func (t *TreeSection) Draw() {
	if t == nil {
		return
	}
	for _, r := range t.rings {
		r.draw()
	}
}

func (t *TreeSection) AddRing() {
	lastRing := t.rings[len(t.rings)-1]
	t.rings = append(t.rings, NewRing(lastRing, lastRing.ampl, lastRing.smoothness))
}

func NewTreeSection(numberOfRings int, ampl float64, smoothness float64) *TreeSection {
	tree := &TreeSection{
		rings: []*Ring{NewRing(nil, ampl, smoothness)},
	}
	for range numberOfRings - 1 {
		tree.AddRing()
	}
	return tree
}

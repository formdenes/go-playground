package treesection

import (
	"math"

	"github.com/KEINOS/go-noise"
	"github.com/quartercastle/vector"
)

const NOISE_TYPE = noise.Perlin
const RESOLUTION int = 360
const INIT_RADIAL float64 = 5
const NOISE_SEED int64 = 1000

const PI_RAD int = 180

type vec = vector.Vector

type Ring struct {
	points     [RESOLUTION]vec
	smoothness float64
	ampl       float64
}

func (r *Ring) draw() {
	if r == nil {
		return
	}
	// for i, p := range r.points {
	// 	var prevP vec
	// 	if i > 0 {
	// 		prevP = r.points[i-1]
	// 	} else {
	// 		prevP = r.points[RESOLUTION-1]
	// 	}
	// 	p5.Line(prevP.X(), prevP.Y(), p.X(), p.Y())
	// }
}

func (r *Ring) GetPoints() [RESOLUTION]vec {
	if r == nil {
		return [RESOLUTION]vec{}
	}
	return r.points
}

func NewRing(prevRing *Ring, ampl float64, smoothness float64, count int) *Ring {
	// if prevRing == nil {
	// 	var points [RESOLUTION]vec
	// 	for i := range RESOLUTION {
	// 		radians := float64(i) * (math.Pi / float64(PI))
	// 		x := (INIT_RADIAL * math.Cos(radians))
	// 		y := (INIT_RADIAL * math.Sin(radians))
	// 		points[i] = vec{x, y}
	// 	}
	// 	return &Ring{
	// 		points:     points,
	// 		smoothness: smoothness,
	// 		ampl:       ampl,
	// 	}
	// }
	var points [RESOLUTION]vec
	ng, err := noise.New(NOISE_TYPE, NOISE_SEED*int64(count+1))
	if err != nil {
		return nil
	}
	ns, err := noise.New(NOISE_TYPE, NOISE_SEED*5000*int64(count+1))
	if err != nil {
		return nil
	}
	nss, err := noise.New(NOISE_TYPE, NOISE_SEED*15000*int64(count+1))
	if err != nil {
		return nil
	}
	var prevPoints [RESOLUTION]vec
	if prevRing == nil {
		for i := range RESOLUTION {
			radians := float64(i) * (math.Pi / float64(PI_RAD))
			x := (INIT_RADIAL * math.Cos(radians))
			y := (INIT_RADIAL * math.Sin(radians))
			prevPoints[i] = vec{x, y}
		}
	} else {
		prevPoints = prevRing.points
	}
	for i := range RESOLUTION {
		prevPoint := prevPoints[i]
		drg := (ng.Eval64(prevPoint.X()/smoothness, prevPoint.Y()/smoothness) + 1) * ampl
		drs := (ns.Eval64(prevPoint.X()/smoothness, prevPoint.Y()/smoothness) + 1) * ampl / 10
		drss := (nss.Eval64(prevPoint.X()/smoothness, prevPoint.Y()/smoothness) + 1) * ampl / 15
		noiseVal := drg + drs + drss
		step := prevPoint.Unit().Scale(noiseVal)
		point := prevPoint.Add(step)
		points[i] = point
	}
	return &Ring{
		points:     points,
		smoothness: smoothness,
		ampl:       ampl,
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

func (t *TreeSection) GetRings() []*Ring {
	if t == nil {
		return nil
	}
	return t.rings
}

func (t *TreeSection) AddRing() {
	lastRing := t.rings[len(t.rings)-1]
	t.rings = append(t.rings, NewRing(lastRing, lastRing.ampl, lastRing.smoothness, len(t.rings)))
}

func NewTreeSection(numberOfRings int, ampl float64, smoothness float64) *TreeSection {
	tree := &TreeSection{
		rings: []*Ring{NewRing(nil, ampl, smoothness, 0)},
	}
	for range numberOfRings - 1 {
		tree.AddRing()
	}
	return tree
}

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

type Options struct {
	NumberOfRings int
	Ampl          float64
	MinDistance   float64
	Smoothness    float64
	Seed          *int64
	NoiseType     *noise.Algo
	Time          *float64
}

func (o *Options) SetDefaults() {
	if o.Seed == nil {
		noiseSeed := NOISE_SEED
		o.Seed = &noiseSeed
	}
	if o.NoiseType == nil {
		noiseType := NOISE_TYPE
		o.NoiseType = &noiseType
	}
	if o.Time == nil {
		time := float64(0)
		o.Time = &time
	}
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

func NewRing(prevRing *Ring, count int, options Options) *Ring {
	var points [RESOLUTION]vec
	options.SetDefaults()
	noiseSeed := *options.Seed * int64(count+1)
	noiseType := *options.NoiseType
	smoothness := options.Smoothness
	time := *options.Time
	ampl := options.Ampl
	ng, err := noise.New(noiseType, noiseSeed)
	if err != nil {
		return nil
	}
	ns, err := noise.New(noiseType, 5000*noiseSeed)
	if err != nil {
		return nil
	}
	nss, err := noise.New(noiseType, 15000*noiseSeed)
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
		drg := (ng.Eval64(prevPoint.X()/smoothness, prevPoint.Y()/smoothness, time) + 1) * ampl
		drs := (ns.Eval64(prevPoint.X()/smoothness, prevPoint.Y()/smoothness, time) + 1) * ampl / 10
		drss := (nss.Eval64(prevPoint.X()/smoothness, prevPoint.Y()/smoothness, time) + 1) * ampl / 15
		noiseVal := math.Max((drg + drs + drss), options.MinDistance)
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
	rings   []*Ring
	options Options
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
	t.rings = append(t.rings, NewRing(lastRing, len(t.rings), t.options))
}

func NewTreeSection(options Options) *TreeSection {
	tree := &TreeSection{
		rings:   []*Ring{NewRing(nil, 0, options)},
		options: options,
	}
	for range options.NumberOfRings - 1 {
		tree.AddRing()
	}
	return tree
}

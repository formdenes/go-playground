package main

import (
	treesectiondraw "main/tree-section-draw"

	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
)

const w float64 = 1000
const h float64 = 2000

const name string = "treesection-test.svg"

func main() {

	c := canvas.New(w, h)
	ctx := canvas.NewContext(c)

	treesectiondraw.Draw(ctx, w, h)

	if err := renderers.Write(name, c); err != nil {
		panic(err)
	}
}

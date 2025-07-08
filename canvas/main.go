package main

import (
	"fmt"
	treesectiondraw "main/tree-section-draw"
	"math/rand/v2"

	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
)

const w float64 = 490
const h float64 = 390

const imageNum int = 1

const dir string = "export/treesection/"
const name string = "treesection"

func main() {

	c := canvas.New(w, h)
	for i := 0; i < imageNum; i++ {
		ctx := canvas.NewContext(c)

		seed := rand.Int64N(100000)

		treesectiondraw.Draw(ctx, w, h, seed)

		// if err := renderers.Write("treesection-test.svg", c); err != nil {
		// 	panic(err)
		// }

		if err := renderers.Write(fmt.Sprintf("%s%s-%d.svg", dir, name, seed), c); err != nil {
			panic(err)
		}
	}
}

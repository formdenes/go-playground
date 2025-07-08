package utils

import (
	"image/color"

	"github.com/tdewolff/canvas"
)

func DrawBackground(ctx *canvas.Context, color color.Color) {
	ctx.SetFillColor(color)
	bgRect := canvas.Rectangle(ctx.Width(), ctx.Height())
	ctx.DrawPath(0, 0, bgRect)
}

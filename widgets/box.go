package widgets

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Box struct {
	PColor   color.Color
	SColor   color.Color
	Inverted func() bool
}

func (b *Box) HandleDraw(screen *ebiten.Image, frame image.Rectangle) {
	// Color inversion
	pColor := b.PColor
	sColor := b.SColor
	if b.Inverted() {
		pColor = b.SColor
		sColor = b.PColor
	}

	// Circle geometry
	r := float32(math.Min(float64(frame.Dx()), float64(frame.Dy()))) / 2
	diff := frame.Bounds().Size().Sub(frame.Bounds().Size().Div(2))

	vector.DrawFilledCircle(
		screen,
		float32(frame.Min.X+diff.X),
		float32(frame.Min.Y+diff.Y),
		float32(r),
		pColor,
		true,
	)
	vector.StrokeCircle(
		screen,
		float32(frame.Min.X+diff.X),
		float32(frame.Min.Y+diff.Y),
		float32(r),
		1.2,
		sColor,
		true,
	)

}

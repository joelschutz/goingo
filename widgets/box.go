package widgets

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Box struct {
	Color color.Color
}

func (b *Box) HandleDraw(screen *ebiten.Image, frame image.Rectangle) {
	ebitenutil.DrawRect(
		screen,
		float64(frame.Min.X),
		float64(frame.Min.Y),
		float64(frame.Size().X),
		float64(frame.Size().Y),
		b.Color,
	)
}

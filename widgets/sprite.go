package widgets

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

type Sprite struct {
	Image    image.Image
	Scale    float64
	Inverted func() bool
}

func (b *Sprite) HandleDraw(screen *ebiten.Image, frame image.Rectangle) {
	img := ebiten.NewImageFromImage(b.Image)
	diff := frame.Bounds().Size().Sub(img.Bounds().Size())

	op := &colorm.DrawImageOptions{}
	op.GeoM.Translate(float64(frame.Min.X+diff.X/2), float64(frame.Min.Y+diff.Y/2))
	op.GeoM.Scale(b.Scale, b.Scale)
	c := colorm.ColorM{}
	if b.Inverted() {
		c.Scale(-1, -1, -1, 1)
		c.Translate(1, 1, 1, 0)
	}
	colorm.DrawImage(screen, img, c, op)
}

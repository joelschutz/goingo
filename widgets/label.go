package widgets

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Label struct {
	Box
	Text     string
	TextFunc func() string
	Font     font.Face
}

func (lb *Label) HandleDraw(screen *ebiten.Image, frame image.Rectangle) {
	lb.Box.HandleDraw(screen, frame)

	txt := lb.Text
	if lb.TextFunc != nil {
		txt = lb.TextFunc()
	}
	txtB := text.BoundString(lb.Font, txt).Size()
	diff := frame.Bounds().Size().Sub(txtB)

	clr := lb.SColor
	if lb.Inverted() {
		clr = lb.PColor
	}

	text.Draw(screen, txt, lb.Font, frame.Min.X+diff.X/2, frame.Min.Y+diff.Y/2+txtB.Y*3/4, clr)
}

func GetDefaultFont(size float64) font.Face {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	mplusNormalFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
	return mplusNormalFont
}

func GetFont(size float64, file []byte) font.Face {
	tt, err := opentype.Parse(file)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	mplusNormalFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
	return mplusNormalFont
}

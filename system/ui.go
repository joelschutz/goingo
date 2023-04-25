package system

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/joelschutz/goingo/component"
	"github.com/yohamta/donburi/ecs"
)

var (
	mplusNormalFont font.Face
	fontSize        = 24
)

type UIRender struct {
	board         *component.BoardState
	bounds        *image.Rectangle
	configuration *component.ConfigurationData
}

func NewUIRender(bounds *image.Rectangle) *UIRender {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(fontSize),
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}

	return &UIRender{
		bounds: bounds,
	}
}

func (r *UIRender) Draw(ecs *ecs.ECS, screen *ebiten.Image) {
	if r.board == nil {
		if entry, ok := component.Board.First(ecs.World); ok {
			r.board = component.Board.Get(entry)
		}
	}

	if r.configuration == nil {
		if entry, ok := component.Configuration.First(ecs.World); ok {
			r.configuration = component.Configuration.Get(entry)
		}
	}

	boardSize := float32(math.Min(float64(r.bounds.Dx()), float64(r.bounds.Dy())))
	cellSize := boardSize / float32(r.configuration.BoardSize+1)

	for i, amount := range r.board.Points {
		x := 1 + float32(i)
		y := float32(.5)
		switch component.Oponent(i + 1) {
		case component.BLACK:
			vector.DrawFilledCircle(screen, x*cellSize, y*cellSize, cellSize/4, color.Black, true)
			text.Draw(screen, fmt.Sprint(amount), mplusNormalFont, int(x*cellSize)-fontSize/4, int(y*cellSize)+fontSize/4, color.White)
		case component.WHITE:
			vector.DrawFilledCircle(screen, x*cellSize, y*cellSize, cellSize/4, color.White, true)
			text.Draw(screen, fmt.Sprint(amount), mplusNormalFont, int(x*cellSize)-fontSize/4, int(y*cellSize)+fontSize/4, color.Black)
		}
	}

}

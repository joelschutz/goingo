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
	"github.com/joelschutz/goingo/util"
	"github.com/yohamta/donburi/ecs"
)

const (
	ALPHABET = "ABCDEFGHJKLMNOPQRSTUVWXYZ"
)

var (
	mplusNormalFont font.Face
	fontSize        = 24
)

type BoardRender struct {
	board         *component.BoardState
	bounds        *image.Rectangle
	cursor        *component.Position
	configuration *component.ConfigurationData
}

func NewRender(bounds *image.Rectangle) *BoardRender {
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
	return &BoardRender{
		bounds: bounds,
	}
}

func (r *BoardRender) Draw(ecs *ecs.ECS, screen *ebiten.Image) {
	if r.board == nil {
		if entry, ok := component.Board.First(ecs.World); ok {
			r.board = component.Board.Get(entry)
		}
	}

	if r.cursor == nil {
		if entry, ok := component.PositionComponent.First(ecs.World); ok {
			r.cursor = component.PositionComponent.Get(entry)
		}
	}

	if r.configuration == nil {
		if entry, ok := component.Configuration.First(ecs.World); ok {
			r.configuration = component.Configuration.Get(entry)
		}
	}

	boardSize := float32(math.Min(float64(r.bounds.Dx()), float64(r.bounds.Dy())))
	cellSize := boardSize / float32(r.configuration.BoardSize+1)
	clrLines := util.DARK_GREY
	clrOutine := color.Black
	if r.configuration.DarkMode {
		clrLines = util.GREY
		clrOutine = color.White
	}

	// Draw Grid Lines
	for i := 0; i < int(r.configuration.BoardSize); i++ {
		vector.StrokeLine(
			screen,
			float32(i+1)*cellSize,
			cellSize,
			float32(i+1)*cellSize,
			boardSize-cellSize,
			2,
			clrLines,
			true,
		)
		vector.StrokeLine(
			screen,
			cellSize,
			float32(i+1)*cellSize,
			boardSize-cellSize,
			float32(i+1)*cellSize,
			2,
			clrLines,
			true,
		)
	}

	// Draw Reference Points
	{
		// Draw Center Point
		vector.DrawFilledCircle(
			screen,
			boardSize/2,
			boardSize/2,
			5,
			clrLines,
			true,
		)

		var sections int = 4
		if r.configuration.BoardSize == 19 {
			sections = 6
		}
		offset := float32((r.configuration.BoardSize - 1) / sections)

		// Draw Corner Points
		vector.DrawFilledCircle(
			screen,
			(offset+1)*cellSize,
			(offset+1)*cellSize,
			5,
			clrLines,
			true,
		)
		vector.DrawFilledCircle(
			screen,
			(offset+1)*cellSize,
			(float32(sections-1)*offset+1)*cellSize,
			5,
			clrLines,
			true,
		)
		vector.DrawFilledCircle(
			screen,
			(float32(sections-1)*offset+1)*cellSize,
			(float32(sections-1)*offset+1)*cellSize,
			5,
			clrLines,
			true,
		)
		vector.DrawFilledCircle(
			screen,
			(float32(sections-1)*offset+1)*cellSize,
			(offset+1)*cellSize,
			5,
			clrLines,
			true,
		)

		// Draw Extra Points
		if r.configuration.BoardSize == 19 {
			vector.DrawFilledCircle(
				screen,
				(offset+1)*cellSize,
				(float32(sections/2)*offset+1)*cellSize,
				5,
				clrLines,
				true,
			)
			vector.DrawFilledCircle(
				screen,
				(float32(sections/2)*offset+1)*cellSize,
				(offset+1)*cellSize,
				5,
				clrLines,
				true,
			)
			vector.DrawFilledCircle(
				screen,
				(float32(sections/2)*offset+1)*cellSize,
				(float32(sections-1)*offset+1)*cellSize,
				5,
				clrLines,
				true,
			)
			vector.DrawFilledCircle(
				screen,
				(float32(sections-1)*offset+1)*cellSize,
				(float32(sections/2)*offset+1)*cellSize,
				5,
				clrLines,
				true,
			)
		}
	}

	// Draw Cursor
	if image.Point(*r.cursor).In(image.Rect(0, 0, int(r.configuration.BoardSize), int(r.configuration.BoardSize))) {
		c := util.TRANSPARENT_BLACK
		if r.board.PlayerTurn == component.WHITE_TURN {
			c = util.TRANSPARENT_WHITE
		}
		vector.DrawFilledCircle(screen, float32(r.cursor.X+1)*cellSize, float32(r.cursor.Y+1)*cellSize, cellSize/2, c, true)
	}

	// Draw Stones
	for i, pos := range r.board.Stones {
		x := float32(i%int(r.configuration.BoardSize)) + 1
		y := float32(i/int(r.configuration.BoardSize)) + 1
		switch pos {
		case component.BLACK:
			vector.DrawFilledCircle(screen, x*cellSize, y*cellSize, cellSize/2, color.Black, true)
		case component.WHITE:
			vector.DrawFilledCircle(screen, x*cellSize, y*cellSize, cellSize/2, color.White, true)
		}
		if pos != component.EMPTY {
			vector.StrokeCircle(screen, x*cellSize, y*cellSize, cellSize/2, 2, clrOutine, true)
		}
	}

	// Draw Labels
	for i := int(r.configuration.BoardSize); i > 0; i-- {
		text.Draw(screen, fmt.Sprint((r.configuration.BoardSize+1)-i), mplusNormalFont, int(boardSize-cellSize*2/3), int(float32(i)*cellSize)+fontSize/4, clrLines)
		text.Draw(screen, fmt.Sprintf("%c", ALPHABET[i-1]), mplusNormalFont, int(float32(i)*cellSize)-fontSize/4, int(boardSize-cellSize/2), clrLines)
	}
}

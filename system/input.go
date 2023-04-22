package system

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/joelschutz/goingo/component"
	"github.com/joelschutz/goingo/event"
	"github.com/yohamta/donburi/ecs"
)

type InputManager struct {
	bounds        *image.Rectangle
	cursor        *component.Position
	configuration *component.ConfigurationData
	board         *component.BoardState
}

func NewInputManager(bounds *image.Rectangle) *InputManager {
	return &InputManager{
		bounds: bounds,
	}
}

func (im *InputManager) Update(ecs *ecs.ECS) {
	if im.cursor == nil {
		cursor := ecs.World.Create(component.PositionComponent)
		im.cursor = component.PositionComponent.Get(ecs.World.Entry(cursor))
	}

	if im.board == nil {
		if entry, ok := component.Board.First(ecs.World); ok {
			im.board = component.Board.Get(entry)
		}
	}

	if im.configuration == nil {
		if entry, ok := component.Configuration.First(ecs.World); ok {
			im.configuration = component.Configuration.Get(entry)
		}
	}

	boardSize := float32(math.Min(float64(im.bounds.Dx()), float64(im.bounds.Dy())))
	cellSize := boardSize / float32(im.configuration.BoardSize+1)

	cX, cY := ebiten.CursorPosition()
	im.cursor.X = ((cX + int(cellSize/2)) / int(cellSize)) - 1
	im.cursor.Y = ((cY + int(cellSize/2)) / int(cellSize)) - 1

	if image.Point(*im.cursor).In(image.Rect(0, 0, int(im.configuration.BoardSize), int(im.configuration.BoardSize))) {
		for btn := ebiten.MouseButton0; btn < ebiten.MouseButtonMax; btn++ {
			if inpututil.IsMouseButtonJustReleased(btn) {
				event.InteractionEvent.Publish(ecs.World, &event.Interaction{
					Position: component.Position{
						X: im.cursor.X,
						Y: im.cursor.Y},
					Button: btn,
				})
			}
		}
	}
}

package event

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/goingo/component"
	"github.com/yohamta/donburi/features/events"
)

type Interaction struct {
	Position component.Position
	Button   ebiten.MouseButton
}

var InteractionEvent = events.NewEventType[*Interaction]()

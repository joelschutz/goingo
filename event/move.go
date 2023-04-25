package event

import (
	"github.com/joelschutz/goingo/component"
	"github.com/yohamta/donburi/features/events"
)

type Move struct {
	Position component.Position
	Player   component.Turn
}

var MoveEvent = events.NewEventType[*Move]()

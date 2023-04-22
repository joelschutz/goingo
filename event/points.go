package event

import (
	"github.com/joelschutz/goingo/component"
	"github.com/yohamta/donburi/features/events"
)

type Points struct {
	Amount float32
	Target component.Oponent
}

var PointsEvent = events.NewEventType[*Points]()

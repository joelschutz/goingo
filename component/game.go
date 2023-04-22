package component

import (
	"image"

	"github.com/yohamta/donburi"
)

type Oponent int

const (
	EMPTY Oponent = iota
	BLACK
	WHITE
)

type BoardState struct {
	Pieces     []Oponent
	Points     [2]float32
	PlayerTurn Turn
}

var Board = donburi.NewComponentType[BoardState]()

type Turn bool

const (
	BLACK_TURN Turn = false
	WHITE_TURN Turn = true
)

func (pt *Turn) Toggle() {
	t := !*pt
	pt = &t
}

func (pt Turn) String() string {
	if pt {
		return "White"
	}
	return "Black"
}

type Position image.Point

var PositionComponent = donburi.NewComponentType[Position]()

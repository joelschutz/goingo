package system

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/goingo/component"
	"github.com/joelschutz/goingo/event"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type RuleSystem struct {
	board         *component.BoardState
	configuration *component.ConfigurationData
}

func NewRuleSystem(ecs *ecs.ECS) *RuleSystem {
	rs := RuleSystem{}
	if entry, ok := component.Board.First(ecs.World); ok {
		rs.board = component.Board.Get(entry)
	}

	if entry, ok := component.Configuration.First(ecs.World); ok {
		rs.configuration = component.Configuration.Get(entry)
	}
	return &rs
}

func (rs *RuleSystem) HandleInput(w donburi.World, e *event.Interaction) {
	if e.Button == ebiten.MouseButtonRight {
		// Pass turn
		event.MoveEvent.Publish(w, &event.Move{
			Position: component.Position{
				X: -1,
				Y: -1},
			Player: rs.board.PlayerTurn,
		})
		rs.board.PlayerTurn = !rs.board.PlayerTurn
	} else {
		cell := &rs.board.Stones[e.Position.X+e.Position.Y*int(rs.configuration.BoardSize)]
		if *cell == component.EMPTY {
			// Place stone
			if rs.board.PlayerTurn == component.BLACK_TURN {
				*cell = component.BLACK
			} else {
				*cell = component.WHITE
			}
			// Check for capture
			_, _, aEnemies := rs.getArmy(e.Position, nil, nil)
			for _, p := range aEnemies {
				eLib, eConn, _ := rs.getArmy(p, nil, nil)
				if len(eLib) == 0 {
					rs.captureArmy(eConn)
					event.PointsEvent.Publish(w, &event.Points{
						Amount: float32(len(eConn)),
						Target: *cell,
					})
					rs.board.Points[*cell-1] += float32(len(eConn))
				}
			}
			// Check for suicide
			aLib, _, aEnemies := rs.getArmy(e.Position, nil, nil)
			if len(aLib) != 0 {
				// End turn
				event.MoveEvent.Publish(w, &event.Move{
					Position: e.Position,
					Player:   rs.board.PlayerTurn,
				})
				rs.board.PlayerTurn = !rs.board.PlayerTurn
			} else {
				*cell = component.EMPTY
			}
		} else {
			fmt.Println(rs.getArmy(e.Position, nil, nil))
		}
	}

}

func (rs *RuleSystem) captureArmy(army []component.Position) {
	for _, p := range army {
		rs.board.Stones[p.X+p.Y*int(rs.configuration.BoardSize)] = component.EMPTY
	}
}

func (rs *RuleSystem) getArmy(pos component.Position, liberties, connections []component.Position) (lib, conn, enemies []component.Position) {
	targetColor := rs.cellState(pos)

	if liberties == nil {
		liberties = make([]component.Position, 0)
	}
	if connections == nil {
		connections = make([]component.Position, 0)
	}
	connections = append(connections, pos)
	enemies = make([]component.Position, 0)

	if targetColor > component.EMPTY {
		// Cell names
		// |  |v1|  |
		// |v2|v0|v3|
		// |  |v4|  |

		p1 := component.Position{pos.X, pos.Y - 1}
		if v1 := rs.cellState(p1); v1 == component.EMPTY {
			if !rs.inArray(p1, liberties) {
				liberties = append(liberties, p1)
			}
		} else if v1 == targetColor && !rs.inArray(p1, connections) {
			liberties, connections, _ = rs.getArmy(p1, liberties, connections)
		} else if v1 == rs.oponentStone(targetColor) {
			enemies = append(enemies, p1)
		}

		p2 := component.Position{pos.X - 1, pos.Y}
		if v2 := rs.cellState(p2); v2 == component.EMPTY {
			if !rs.inArray(p2, liberties) {
				liberties = append(liberties, p2)
			}
		} else if v2 == targetColor && !rs.inArray(p2, connections) {
			liberties, connections, _ = rs.getArmy(p2, liberties, connections)
		} else if v2 == rs.oponentStone(targetColor) {
			enemies = append(enemies, p2)
		}

		p3 := component.Position{pos.X + 1, pos.Y}
		if v3 := rs.cellState(p3); v3 == component.EMPTY {
			if !rs.inArray(p3, liberties) {
				liberties = append(liberties, p3)
			}
		} else if v3 == targetColor && !rs.inArray(p3, connections) {
			liberties, connections, _ = rs.getArmy(p3, liberties, connections)
		} else if v3 == rs.oponentStone(targetColor) {
			enemies = append(enemies, p3)
		}

		p4 := component.Position{pos.X, pos.Y + 1}
		if v4 := rs.cellState(p4); v4 == component.EMPTY {
			if !rs.inArray(p4, liberties) {
				liberties = append(liberties, p4)
			}
		} else if v4 == targetColor && !rs.inArray(p4, connections) {
			liberties, connections, _ = rs.getArmy(p4, liberties, connections)
		} else if v4 == rs.oponentStone(targetColor) {
			enemies = append(enemies, p4)
		}
		return liberties, connections, enemies
	}
	return nil, nil, nil
}

func (rs *RuleSystem) cellState(pos component.Position) component.Oponent {
	xOut := pos.X < 0 || pos.X > int(rs.configuration.BoardSize-1)
	yOut := pos.Y < 0 || pos.Y > int(rs.configuration.BoardSize-1)

	if xOut || yOut {
		return -1
	}
	return rs.board.Stones[pos.X+pos.Y*int(rs.configuration.BoardSize)]
}

func (rs *RuleSystem) inArray(pos component.Position, arr []component.Position) bool {
	for _, v := range arr {
		if pos == v {
			return true
		}
	}
	return false
}

func (rs *RuleSystem) oponentStone(stone component.Oponent) component.Oponent {
	if stone == component.BLACK {
		return component.WHITE
	}
	return component.BLACK
}

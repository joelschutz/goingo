package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/joelschutz/goingo/component"
	"github.com/joelschutz/goingo/event"
	"github.com/joelschutz/goingo/layers"
	"github.com/joelschutz/goingo/system"
	"github.com/joelschutz/goingo/util"
	"github.com/joelschutz/goingo/widgets"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/events"
	"github.com/yohamta/furex/v2"
)

type GameScene struct {
	sm     *util.SceneManager[donburi.World]
	ecs    *ecs.ECS
	bounds image.Rectangle
}

func (g *GameScene) Unload() donburi.World {
	return g.ecs.World
}

func (g *GameScene) Load(w donburi.World, sm *util.SceneManager[donburi.World]) {
	g.sm = sm
	g.bounds = image.Rectangle{}

	entry, _ := component.Configuration.First(w)
	conf := component.Configuration.Get(entry)
	boardSize := conf.BoardSize
	assets := conf.Assets
	w, board := spawnWorld(w, boardSize)

	g.ecs = ecs.NewECS(w)

	rules := system.NewRuleSystem(g.ecs)
	event.InteractionEvent.Subscribe(g.ecs.World, rules.HandleInput)

	// Load Sounds
	sound := system.NewAudioPlayer(assets, audio.NewContext(48000))
	event.MoveEvent.Subscribe(g.ecs.World, sound.HandleMove)
	event.PointsEvent.Subscribe(g.ecs.World, sound.HandleCapture)

	// Load Images
	iBulb, _, err := util.LoadImage("assets/icons8-idea-60.png", conf.Assets)
	if err != nil {
		panic(err)
	}

	// Load AI
	ai, err := system.NewAIConn(g.ecs, boardSize)
	if err == nil {
		event.MoveEvent.Subscribe(g.ecs.World, ai.HandleMove)
	}

	ui := system.NewMenuRender(&g.bounds,
		func(r *system.MenuRender) {
			// Create UI
			r.GameUI = &furex.View{
				Width:     r.Bounds.Dx(),
				Height:    r.Bounds.Dy(),
				Direction: furex.Column,
				Justify:   furex.JustifyStart,
				Wrap:      furex.Wrap,
			}

			// Add Header
			header := &furex.View{
				Grow:      1,
				Justify:   furex.JustifyStart,
				MarginTop: 10,
			}
			header.AddChild(&furex.View{
				Width: 90,
				Handler: &widgets.Label{
					Box: widgets.Box{
						PColor:   color.Black,
						SColor:   color.White,
						Inverted: func() bool { return false },
					},
					Font: widgets.GetDefaultFont(18),
					TextFunc: func() string {
						return fmt.Sprint(board.Points[0])
					},
				},
			})
			header.AddChild(&furex.View{
				Width: 90,
				Handler: &widgets.Label{
					Box: widgets.Box{
						PColor:   color.White,
						SColor:   color.Black,
						Inverted: func() bool { return false },
					},
					Font: widgets.GetDefaultFont(18),
					TextFunc: func() string {
						return fmt.Sprint(board.Points[1])
					},
				},
			})
			header.AddChild(&furex.View{Grow: 1})
			header.AddChild(&furex.View{
				Width: 90,
				Handler: &widgets.SpriteButton{
					Action: func() { conf.DarkMode = !conf.DarkMode },
					Sprite: widgets.Sprite{
						Inverted: func() bool { return conf.DarkMode },
						Image:    iBulb,
						Scale:    1,
					},
				},
			})
			r.GameUI.AddChild(header)

			// Add Board
			r.GameUI.AddChild(&furex.View{
				Grow: float64(conf.BoardSize + 5),
			})
		})
	g.ecs.
		AddSystem(ui.Update).
		AddSystem(system.NewInputManager(&g.bounds).Update).
		AddRenderer(layers.LayerBackground, system.Background.DrawBackground).
		AddRenderer(layers.LayerBoard, system.NewRender(&g.bounds).Draw).
		AddRenderer(layers.LayerUI, ui.Draw)
}

func (g *GameScene) Update() error {
	events.ProcessAllEvents(g.ecs.World)
	g.ecs.Update()
	return nil
}

func (g *GameScene) Draw(screen *ebiten.Image) {
	screen.Clear()
	g.ecs.DrawLayer(layers.LayerBackground, screen)
	g.ecs.DrawLayer(layers.LayerBoard, screen)
	g.ecs.DrawLayer(layers.LayerUI, screen)
}

func (g *GameScene) Layout(width, height int) (int, int) {
	g.bounds = image.Rect(0, 0, width, height)
	return width, height
}

func spawnWorld(world donburi.World, boardSize int) (donburi.World, *component.BoardState) {
	board := world.Entry(world.Create(component.Board))
	donburi.SetValue(
		board,
		component.Board,
		component.BoardState{
			Stones: make([]component.Oponent, boardSize*boardSize),
		})
	world.Create(component.PositionComponent)

	return world, component.Board.Get(board)
}

package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/goingo/component"
	"github.com/joelschutz/goingo/layers"
	"github.com/joelschutz/goingo/system"
	"github.com/joelschutz/goingo/util"
	"github.com/joelschutz/goingo/widgets"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/furex/v2"
)

type MenuScene struct {
	sm         *util.SceneManager[donburi.World]
	ecs        *ecs.ECS
	bounds     image.Rectangle
	AINotFound bool
}

func (g *MenuScene) Unload() donburi.World {
	return g.ecs.World
}

func (g *MenuScene) Load(w donburi.World, sm *util.SceneManager[donburi.World]) {
	g.sm = sm
	g.bounds = image.Rectangle{}

	entry, _ := component.Configuration.First(w)
	conf := component.Configuration.Get(entry)
	boardSize := conf.BoardSize

	g.ecs = ecs.NewECS(spawnWorld2(w, boardSize))

	// Check if GNUGo is present
	_, err := system.NewAIConn(g.ecs, boardSize)
	if err != nil {
		g.AINotFound = true
		conf.AIEnabled = false
	}
	// event.MoveEvent.Subscribe(g.ecs.World, ai.HandleMove)
	fontS := widgets.GetDefaultFont(18)
	fontL := widgets.GetDefaultFont(32)

	ui := system.NewMenuRender(&g.bounds,
		func(r *system.MenuRender) {
			colors := []color.Color{
				color.RGBA{0x59, 0x98, 0x1a, 0xff},
				color.RGBA{0x81, 0xb6, 0x22, 0xff},
				util.GREY,
			}
			actions := []func(){
				func() {
					var sel int
					for i, v := range SizeOptions {
						if v == conf.BoardSize {
							sel = i + 1
							break
						}
					}
					if sel == len(SizeOptions) {
						sel = 0
					}
					conf.BoardSize = SizeOptions[sel]
				},
				func() {
					if !g.AINotFound {
						conf.AIEnabled = !conf.AIEnabled
					}
				},
				func() { g.sm.Load(&GameScene{}) },
			}
			texts := []func() string{
				func() string { return fmt.Sprintf("Size x%d", conf.BoardSize) },
				func() string {
					opp := "Player"
					if conf.AIEnabled {
						opp = "AI"
					}
					return fmt.Sprintf("vs%s", opp)
				},
				func() string { return "Play" },
			}
			r.GameUI = &furex.View{
				Width:        r.Bounds.Dx(),
				Height:       r.Bounds.Dy(),
				Direction:    furex.Column,
				Justify:      furex.JustifySpaceAround,
				AlignItems:   furex.AlignItemCenter,
				AlignContent: furex.AlignContentCenter,
				Wrap:         furex.Wrap,
			}
			r.GameUI.AddChild(&furex.View{
				Width:        500,
				Height:       100,
				MarginBottom: 100,
				MarginTop:    100,
				Handler: &widgets.Label{
					Box:  widgets.Box{colors[len(colors)-1]},
					Text: "GoinGo",
					Font: fontL,
				},
			})
			menu := &furex.View{
				Direction:    furex.Row,
				AlignItems:   furex.AlignItemCenter,
				Justify:      furex.JustifyCenter,
				Width:        600,
				Grow:         1,
				Shrink:       1,
				MarginBottom: 100,
			}
			r.GameUI.AddChild(menu)

			for i := 0; i < 3; i++ {
				menu.AddChild(&furex.View{
					Height: 200,
					Shrink: 1,
					Grow:   1,
					Handler: &widgets.Button{
						Action: actions[i],
						Label: widgets.Label{
							Box:      widgets.Box{colors[i%len(colors)]},
							TextFunc: texts[i],
							Font:     fontS,
						},
					},
				})
			}
		},
	)
	g.ecs.
		AddSystem(ui.Update).
		AddRenderer(layers.LayerBackground, system.DrawBackground).
		AddRenderer(layers.LayerUI, ui.Draw)
}

func (g *MenuScene) Update() error {
	// events.ProcessAllEvents(g.ecs.World)
	g.ecs.Update()
	return nil
}

func (g *MenuScene) Draw(screen *ebiten.Image) {
	screen.Clear()
	g.ecs.DrawLayer(layers.LayerBackground, screen)
	// g.ecs.DrawLayer(layers.LayerBoard, screen)
	g.ecs.DrawLayer(layers.LayerUI, screen)
}

func (g *MenuScene) Layout(width, height int) (int, int) {
	g.bounds = image.Rect(0, 0, width, height)
	return width, height
}

func spawnWorld2(world donburi.World, boardSize int) donburi.World {
	// board := world.Entry(world.Create(component.Board))
	// donburi.SetValue(
	// 	board,
	// 	component.Board,
	// 	component.BoardState{
	// 		Pieces: make([]component.Oponent, boardSize*boardSize),
	// 	})

	return world
}

package main

import (
	"image"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/goingo/component"
	"github.com/joelschutz/goingo/event"
	"github.com/joelschutz/goingo/layers"
	"github.com/joelschutz/goingo/system"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/events"
)

type Game struct {
	ecs    *ecs.ECS
	bounds image.Rectangle
}

func (g *Game) configure() {
	g.bounds = image.Rectangle{}
	g.ecs = ecs.NewECS(spawnWorld(13))

	rules := system.NewRuleSystem(g.ecs)
	event.InteractionEvent.Subscribe(g.ecs.World, rules.HandleInput)

	g.ecs.
		AddSystem(system.NewInputManager(&g.bounds).Update).
		AddRenderer(layers.LayerBackground, system.DrawBackground).
		AddRenderer(layers.LayerBoard, system.NewRender(&g.bounds).Draw).
		AddRenderer(layers.LayerUI, system.NewUIRender(&g.bounds).Draw)
}

func (g *Game) Update() error {
	events.ProcessAllEvents(g.ecs.World)
	g.ecs.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	g.ecs.DrawLayer(layers.LayerBackground, screen)
	g.ecs.DrawLayer(layers.LayerBoard, screen)
	g.ecs.DrawLayer(layers.LayerUI, screen)
	// g.ecs.DrawLayer(layers.LayerMetrics, screen)
}

func (g *Game) Layout(width, height int) (int, int) {
	g.bounds = image.Rect(0, 0, width, height)
	return width, height
}

func spawnWorld(boardSize uint) donburi.World {
	world := donburi.NewWorld()
	configuration := world.Entry(world.Create(component.Configuration))
	donburi.SetValue(
		configuration,
		component.Configuration,
		component.ConfigurationData{
			BoardSize: boardSize,
		},
	)
	// for i := 0; i < int(boardSize*boardSize); i++ {
	// 	b[i] = rand.Intn(3)
	// }
	board := world.Entry(world.Create(component.Board))
	donburi.SetValue(
		board,
		component.Board,
		component.BoardState{
			Pieces: make([]component.Oponent, boardSize*boardSize),
		})

	return world
}

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowSizeLimits(300, 200, -1, -1)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	rand.Seed(time.Now().UTC().UnixNano())
	g := &Game{}
	g.configure()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

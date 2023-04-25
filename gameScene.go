package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/joelschutz/goingo/component"
	"github.com/joelschutz/goingo/event"
	"github.com/joelschutz/goingo/layers"
	"github.com/joelschutz/goingo/system"
	"github.com/joelschutz/goingo/util"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/events"
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

	g.ecs = ecs.NewECS(spawnWorld(w, boardSize))

	rules := system.NewRuleSystem(g.ecs)
	event.InteractionEvent.Subscribe(g.ecs.World, rules.HandleInput)

	sound := system.NewAudioPlayer(assets, audio.NewContext(48000))
	event.MoveEvent.Subscribe(g.ecs.World, sound.HandleMove)
	event.PointsEvent.Subscribe(g.ecs.World, sound.HandleCapture)

	ai, _ := system.NewAIConn(g.ecs, boardSize)
	event.MoveEvent.Subscribe(g.ecs.World, ai.HandleMove)
	g.ecs.
		AddSystem(system.NewInputManager(&g.bounds).Update).
		AddRenderer(layers.LayerBackground, system.DrawBackground).
		AddRenderer(layers.LayerBoard, system.NewRender(&g.bounds).Draw).
		AddRenderer(layers.LayerUI, system.NewUIRender(&g.bounds).Draw)
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

func spawnWorld(world donburi.World, boardSize int) donburi.World {
	board := world.Entry(world.Create(component.Board))
	donburi.SetValue(
		board,
		component.Board,
		component.BoardState{
			Stones: make([]component.Oponent, boardSize*boardSize),
		})
	world.Create(component.PositionComponent)

	return world
}

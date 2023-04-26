package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"io"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/goingo/component"
	"github.com/joelschutz/goingo/layers"
	"github.com/joelschutz/goingo/system"
	"github.com/joelschutz/goingo/util"
	"github.com/joelschutz/goingo/widgets"
	"github.com/pkg/browser"
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
	// furex.Debug = true
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

	// Load Fonts
	ff, err := conf.Assets.Open("assets/DejaVuSansCondensed.ttf")
	if err != nil {
		panic(err)
	}
	fFile, err := io.ReadAll(ff)
	if err != nil {
		panic(err)
	}
	fontFace := widgets.GetFont(24, fFile)

	// Load Images
	iLogo, _, err := util.LoadImage("assets/logo.png", conf.Assets)
	if err != nil {
		panic(err)
	}
	iBulb, _, err := util.LoadImage("assets/icons8-idea-60.png", conf.Assets)
	if err != nil {
		panic(err)
	}
	iGithub, _, err := util.LoadImage("assets/icons8-github-60.png", conf.Assets)
	if err != nil {
		panic(err)
	}
	iItch, _, err := util.LoadImage("assets/icons8-itch-io-60.png", conf.Assets)

	ui := system.NewMenuRender(&g.bounds,
		func(r *system.MenuRender) {
			// Create UI
			r.GameUI = &furex.View{
				Width:        r.Bounds.Dx(),
				Height:       r.Bounds.Dy(),
				Direction:    furex.Column,
				Justify:      furex.JustifySpaceAround,
				AlignItems:   furex.AlignItemCenter,
				AlignContent: furex.AlignContentCenter,
				Wrap:         furex.Wrap,
			}

			// Add Logo
			r.GameUI.AddChild(&furex.View{
				Width:        500,
				Height:       100,
				MarginBottom: 100,
				MarginTop:    100,
				Handler: &widgets.Sprite{
					Image:    iLogo,
					Scale:    1,
					Inverted: func() bool { return conf.DarkMode },
				},
			})

			// Add Menu
			{
				pColors := []color.Color{
					color.White,
					color.White,

					color.Black,
				}
				sColors := []color.Color{
					color.Black,
					color.Black,

					color.White,
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
						Height: 100,
						Shrink: 1,
						Grow:   1,
						Handler: &widgets.LabelButton{
							Action: actions[i],
							Label: widgets.Label{
								Box: widgets.Box{
									PColor:   pColors[i%len(pColors)],
									SColor:   sColors[i%len(sColors)],
									Inverted: func() bool { return conf.DarkMode },
								},
								TextFunc: texts[i],
								Font:     fontFace,
							},
						},
					})
				}
			}

			// Add Footer
			{
				actions := []func(){
					func() { browser.OpenURL("https://github.com/joelschutz/goingo") },
					func() { browser.OpenURL("https://kam1sama.itch.io/goingo") },
					func() { conf.DarkMode = !conf.DarkMode },
				}
				imgs := []image.Image{
					iGithub,
					iItch,
					iBulb,
				}
				footer := &furex.View{
					Direction:    furex.Row,
					AlignContent: furex.AlignContentEnd,
					Justify:      furex.JustifySpaceAround,
					Height:       50,
					Width:        r.Bounds.Dx(),
					MarginBottom: 50,
				}
				r.GameUI.AddChild(footer)

				for i := 0; i < 3; i++ {

					footer.AddChild(&furex.View{
						Width: 100,
						Handler: &widgets.SpriteButton{
							Action: actions[i],
							Sprite: widgets.Sprite{
								Inverted: func() bool { return conf.DarkMode },
								Image:    imgs[i],
								Scale:    1,
							},
						},
					})

				}
			}

		},
	)
	g.ecs.
		AddSystem(ui.Update).
		AddRenderer(layers.LayerBackground, system.Background.DrawBackground).
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

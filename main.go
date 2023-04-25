package main

import (
	"embed"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/goingo/component"
	"github.com/joelschutz/goingo/util"
	"github.com/yohamta/donburi"
)

var (
	//go:embed assets
	_assets     embed.FS
	_boardSize  = SizeOptions[0]
	SizeOptions = [3]int{9, 13, 19}
)

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowSizeLimits(300, 200, -1, -1)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	rand.Seed(time.Now().UTC().UnixNano())

	world := donburi.NewWorld()

	configuration := world.Entry(world.Create(component.Configuration))
	donburi.SetValue(
		configuration,
		component.Configuration,
		component.ConfigurationData{
			BoardSize: _boardSize,
			AIEnabled: true,
			Assets:    _assets,
		},
	)
	// sm := util.NewSceneManager[donburi.World](&GameScene{}, world)
	sm := util.NewSceneManager[donburi.World](&MenuScene{}, world)

	if err := ebiten.RunGame(sm); err != nil {
		log.Fatal(err)
	}
}

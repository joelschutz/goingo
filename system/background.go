package system

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/goingo/component"
	"github.com/joelschutz/goingo/util"
	"github.com/yohamta/donburi/ecs"
)

type BackgroundData struct {
	clr           color.Color
	configuration *component.ConfigurationData
}

var Background = &BackgroundData{}

func (b *BackgroundData) DrawBackground(ecs *ecs.ECS, screen *ebiten.Image) {
	// screen.Fill(color.RGBA{235, 198, 111, 255})
	if b.configuration == nil {
		if entry, ok := component.Configuration.First(ecs.World); ok {
			b.configuration = component.Configuration.Get(entry)

		}
	}
	b.clr = util.GREY
	if b.configuration.DarkMode {
		b.clr = util.DARK_GREY
	}
	screen.Fill(b.clr)
}

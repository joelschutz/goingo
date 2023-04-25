package system

import (
	"image"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/joelschutz/goingo/component"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/furex/v2"
)

type MenuRender struct {
	initOnce      sync.Once
	Bounds        *image.Rectangle
	configuration *component.ConfigurationData
	GameUI        *furex.View
	setup         func(r *MenuRender)
}

func NewMenuRender(bounds *image.Rectangle, setup func(r *MenuRender)) *MenuRender {
	return &MenuRender{
		Bounds: bounds,
		setup:  setup,
	}
}
func (r *MenuRender) Update(ecs *ecs.ECS) {
	r.initOnce.Do(func() {
		r.setup(r)
	})
	r.GameUI.UpdateWithSize(r.Bounds.Dx(), r.Bounds.Dy())
}

func (r *MenuRender) Draw(ecs *ecs.ECS, screen *ebiten.Image) {

	if r.configuration == nil {
		if entry, ok := component.Configuration.First(ecs.World); ok {
			r.configuration = component.Configuration.Get(entry)
		}
	}

	r.GameUI.Draw(screen)
}

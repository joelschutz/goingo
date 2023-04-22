package system

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi/ecs"
)

func DrawBackground(ecs *ecs.ECS, screen *ebiten.Image) {
	screen.Fill(color.RGBA{235, 198, 111, 255})
}

package widgets

type LabelButton struct {
	Label
	Action func()
}

func (b *LabelButton) HandleJustPressedMouseButtonLeft(x, y int) bool {

	return true
}
func (b *LabelButton) HandleJustReleasedMouseButtonLeft(x, y int) {
	if b.Action != nil {
		b.Action()
	}
}

type SpriteButton struct {
	Sprite
	Action func()
}

func (b *SpriteButton) HandleJustPressedMouseButtonLeft(x, y int) bool {

	return true
}
func (b *SpriteButton) HandleJustReleasedMouseButtonLeft(x, y int) {
	if b.Action != nil {
		b.Action()
	}
}

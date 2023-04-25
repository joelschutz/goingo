package widgets

type Button struct {
	Label
	Action func()
}

func (b *Button) HandleJustPressedMouseButtonLeft(x, y int) bool {

	return true
}
func (b *Button) HandleJustReleasedMouseButtonLeft(x, y int) {
	if b.Action != nil {
		b.Action()
	}
}

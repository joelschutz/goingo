package util

import "github.com/hajimehoshi/ebiten/v2"

type Scene[T any] interface {
	ebiten.Game
	Load(T, *SceneManager[T])
	Unload() T
}

type SceneManager[T any] struct {
	current Scene[T]
}

func NewSceneManager[T any](scene Scene[T], state T) *SceneManager[T] {
	s := &SceneManager[T]{current: scene}
	scene.Load(state, s)
	return s
}

func (s *SceneManager[T]) Update() error {
	return s.current.Update()
}

func (s *SceneManager[T]) Draw(screen *ebiten.Image) {
	s.current.Draw(screen)
}

func (s *SceneManager[T]) Load(scene Scene[T]) {
	scene.Load(s.current.Unload(), s)
	s.current = scene
}

func (s *SceneManager[T]) Layout(w, h int) (int, int) {
	return s.current.Layout(w, h)
}

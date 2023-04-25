package system

import (
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/joelschutz/goingo/event"
	"github.com/yohamta/donburi"
)

type AudioPlayer struct {
	pMove    *audio.Player
	pCapture *audio.Player
	pPass    *audio.Player
	context  *audio.Context
	assets   fs.FS
}

func NewAudioPlayer(assets fs.FS, context *audio.Context) *AudioPlayer {
	rs := AudioPlayer{
		assets:  assets,
		context: context,
	}

	fMove, err := rs.assets.Open("assets/move0.ogg")
	if err != nil {
		panic(err)
	}
	sMove, err := vorbis.DecodeWithoutResampling(fMove)
	if err != nil {
		panic(err)
	}
	rs.pMove, err = rs.context.NewPlayer(sMove)
	if err != nil {
		panic(err)
	}

	fCapture, err := rs.assets.Open("assets/capture.ogg")
	if err != nil {
		panic(err)
	}
	sCapture, err := vorbis.DecodeWithoutResampling(fCapture)
	if err != nil {
		panic(err)
	}
	rs.pCapture, err = rs.context.NewPlayer(sCapture)
	if err != nil {
		panic(err)
	}

	fPass, err := rs.assets.Open("assets/pass.ogg")
	if err != nil {
		panic(err)
	}
	sPass, err := vorbis.DecodeWithoutResampling(fPass)
	if err != nil {
		panic(err)
	}
	rs.pPass, err = rs.context.NewPlayer(sPass)
	if err != nil {
		panic(err)
	}

	return &rs
}

func (rs *AudioPlayer) HandleMove(w donburi.World, e *event.Move) {
	if e.Position.X < 0 {
		rs.pPass.Rewind()
		rs.pPass.Play()
	} else {
		rs.pMove.Rewind()
		rs.pMove.Play()
	}
}

func (rs *AudioPlayer) HandleCapture(w donburi.World, e *event.Points) {
	rs.pCapture.Rewind()
	rs.pCapture.Play()
}

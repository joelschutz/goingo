package system

import (
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/goingo/component"
	"github.com/joelschutz/goingo/event"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type AIConn struct {
	output        io.Reader
	input         io.Writer
	configuration *component.ConfigurationData
}

func NewAIConn(ecs *ecs.ECS, boardSize int) (*AIConn, error) {

	if exec.Command("gnugo", "-v").Run() != nil {
		return nil, fmt.Errorf("gnugo not found")
	}

	rs := AIConn{}

	if entry, ok := component.Configuration.First(ecs.World); ok {
		rs.configuration = component.Configuration.Get(entry)
	}
	cmd := exec.Command("gnugo", "--boardsize", strconv.Itoa(boardSize), "--mode", "gtp", "--level", "1")

	rs.output, _ = cmd.StdoutPipe()
	rs.input, _ = cmd.StdinPipe()
	cmd.Start()
	return &rs, nil
}

func (rs *AIConn) HandleMove(w donburi.World, e *event.Move) {
	if e.Player == component.BLACK_TURN && rs.configuration.AIEnabled {
		r := make([]byte, 20)

		rs.input.Write([]byte(fmt.Sprintf("play b %s\n", rs.convertPosition(e.Position))))
		rs.input.Write([]byte(fmt.Sprintf("genmove w\n")))
		time.Sleep(time.Millisecond * 500)
		rs.output.Read(r)
		m := cleanString(string(r), "= \n")
		mv := rs.convertString(m)
		fmt.Println("got move:", m, mv)
		btn := ebiten.MouseButtonLeft
		if mv.X < 0 {
			btn = ebiten.MouseButtonRight
		}

		event.InteractionEvent.Publish(w, &event.Interaction{
			Position: mv,
			Button:   btn,
		})
	}
}

func (rs *AIConn) convertPosition(pos component.Position) string {
	if pos.X < 0 {
		return "pass"
	}
	return fmt.Sprintf("%c%v", ALPHABET[pos.X], rs.configuration.BoardSize-pos.Y)
}

func (rs *AIConn) convertString(pos string) component.Position {
	if pos == "pass" || pos == "PASS" {
		return component.Position{
			X: -1,
			Y: -1,
		}
	}

	v, err := strconv.Atoi(string(pos[1:3]))
	if err != nil {
		v, _ = strconv.Atoi(string(pos[1:2]))
	}
	return component.Position{
		X: strings.Index(ALPHABET, string(pos[0])),
		Y: rs.configuration.BoardSize - v,
	}
}

func cleanString(s, target string) string {
	for _, v := range target {
		s = strings.ReplaceAll(s, string(v), "")
	}
	return strings.TrimSpace(s)
}

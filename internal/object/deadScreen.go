package object

import (
	"github.com/3auris/snakery/pkg/grafio"
	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"strconv"
)

// DeadScreen object of game which can be painted
type DeadScreen struct {
	Score  *Score
	Font   ttf.Font
	Screen GameScreen

	toMenu bool
}

func (ds *DeadScreen) HandleEvent(event sdl.Event) {
	switch ev := event.(type) {
	case *sdl.KeyboardEvent:
		if ev.State != sdl.PRESSED {
			break
		}

		switch event.(*sdl.KeyboardEvent).Keysym.Sym {
		case sdl.K_RETURN:
			ds.toMenu = true
		}
	}
}

func (ds *DeadScreen) Update() GameState {
	if ds.toMenu {
		ds.toMenu = false
		return MenuScreen
	}

	return DeadSnake
}

// Paint paints text and Score to renderer
func (ds DeadScreen) Paint(d grafio.Drawer) error {
	if err := d.Background(0, 0, 0, 0); err != nil {
		return errors.Wrap(err, "could not set background")
	}

	opts := grafio.TextOpts{Size: 20, XCof: .05, YCof: .15, Color: grafio.RGBA{R: 255, G: 255, B: 255, A: 255}}

	if err := d.Text("Final score: "+strconv.Itoa(ds.Score.amount), opts); err != nil {
		return errors.Wrap(err, "could not draw amount score text")
	}

	opts.Size = 10
	opts.YCof = .30

	if err := d.Text("Press (Enter) to restart the game", opts); err != nil {
		return errors.Wrap(err, "could not draw how to restart game text")
	}
	return nil
}

package object

import (
	"github.com/3auris/snakery/pkg/grafio"
	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type WelcomeText struct {
	Font   ttf.Font
	Screen GameScreen
	//Snake  *Snake

	changeState bool
}

func (wt *WelcomeText) HandleEvent(event sdl.Event) {
	switch ev := event.(type) {
	case *sdl.KeyboardEvent:
		if ev.State != sdl.PRESSED {
			break
		}

		switch event.(*sdl.KeyboardEvent).Keysym.Sym {
		case sdl.K_RETURN:
			wt.changeState = true
		}
	}
}

func (wt *WelcomeText) Update() GameState {
	if wt.changeState {
		//wt.Snake.reset()

		wt.changeState = false
		return SnakeRunning
	}

	return MenuScreen
}

// Paint paints text and Score to renderer
func (wt WelcomeText) Paint(d grafio.Drawer) error {
	opts := grafio.TextOpts{Size: 10, XCof: .05, YCof: .15, Color: grafio.RGBA{R: 34, G: 139, B: 34, A: 10}}

	if err := d.Text("Welcome to the snake game", opts); err != nil {
		return errors.Wrap(err, "failed to draw the text")
	}

	opts.YCof = .30

	if err := d.Text("Press (Enter) to start the game", opts); err != nil {
		return errors.Wrap(err, "failed to draw the text")
	}

	return nil
}

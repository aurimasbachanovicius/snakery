package object

import (
	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"

	"github.com/3auris/snakery/pkg/grafio"
)

// WelcomeText object of showing welcome text
type WelcomeText struct {
	Screen GameScreen
	Snake  *Snake

	changeState bool
}

// HandleEvent handles events from input devices
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

// Update updates snake and chooses game state to return
func (wt *WelcomeText) Update() GameState {
	if wt.changeState {
		wt.Snake.reset()

		wt.changeState = false
		return SnakeRunning
	}

	return MenuScreen
}

// Paint paints text and Score to renderer
func (wt WelcomeText) Paint(d grafio.Drawer) error {
	opts := grafio.TextOpts{Size: 17, XCof: .05, YCof: .15, Color: grafio.ColorGreen}

	if err := d.Text("Welcome to the snake game", opts); err != nil {
		return errors.Wrap(err, "failed to draw the text")
	}

	opts.Size = 15
	opts.YCof = .30

	if err := d.Text("Press (Enter) to start the game", opts); err != nil {
		return errors.Wrap(err, "failed to draw the text")
	}

	return nil
}

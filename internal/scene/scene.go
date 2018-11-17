package scene

import (
	"fmt"
	"github.com/3auris/snakery/internal/object"
	"github.com/3auris/snakery/pkg/grafio"
	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"os"
	"time"
)

// Scene holds paints and state of the current game
type Scene struct {
	r *sdl.Renderer
	w *sdl.Window

	drawer grafio.Drawer
	state  object.GameState
	paints map[object.GameState][]object.Paintable
}

// New create new Scene with given parameters
func New(fontPath string, screenWidth, screenHeight int32) (*Scene, error) {
	w, r, err := prepareSdl2(int32(screenWidth), int32(screenHeight))
	if err != nil {
		return nil, fmt.Errorf("could not prepare sdl2: %v", err)
	}

	drawer, err := grafio.NewSdl2Draw(r, "res/ubuntu.ttf", screenWidth, screenHeight)
	if err != nil {
		return nil, errors.Wrap(err, "could not create sdl2drawer")
	}

	font, err := ttf.OpenFont(fontPath, 100)
	if err != nil {
		return nil, fmt.Errorf("could not load font: %v", err)
	}

	scrn := object.GameScreen{W: screenWidth, H: screenHeight}

	apple, err := object.NewApple(r)
	if err != nil {
		return nil, fmt.Errorf("could not create apple: %v", err)
	}

	score := object.NewScore(font)
	snake := object.NewSnake(apple, score, font, scrn)
	//deadScreen := &object.DeadScreen{Score: score, Font: *font, Screen: scrn}
	menuScreen := &object.WelcomeText{Font: *font, Screen: scrn, Snake: snake}

	return &Scene{
		r:      r,
		w:      w,
		drawer: drawer,

		state: object.MenuScreen,
		paints: map[object.GameState][]object.Paintable{
			object.MenuScreen:   {menuScreen},
			object.SnakeRunning: {snake, apple, score},
			//object.DeadSnake:    {deadScreen},
		},
	}, nil
}

// Run runs goroutine and updates all paints and listening of events
func (s *Scene) Run(events <-chan sdl.Event) <-chan error {
	errc := make(chan error)

	go func() {
		ticker := time.Tick(55 * time.Millisecond)

		for {
			select {
			case e := <-events:
				for _, paint := range s.paints[s.state] {
					switch p := paint.(type) {
					case object.Handleable:
						p.HandleEvent(e)
					}
				}

				if done := s.handleExit(e); done {
					os.Exit(0)
				}
			case <-ticker:
				s.state = s.update()

				if err := s.paint(); err != nil {
					errc <- err
				}
			}
		}

		os.Exit(0)
	}()

	return errc
}

func (s *Scene) handleExit(event sdl.Event) bool {
	switch ev := event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.KeyboardEvent:
		if ev.State != sdl.PRESSED {
			break
		}

		switch event.(*sdl.KeyboardEvent).Keysym.Sym {
		case sdl.K_ESCAPE:
			return true
		}
	}
	return false
}

func (s Scene) update() object.GameState {
	for _, paint := range s.paints[s.state] {
		switch p := paint.(type) {
		case object.Updateable:
			state := p.Update()
			if state != s.state {
				return state
			}
		}
	}

	return s.state
}

func (s Scene) paint() error {
	err := s.drawer.Presentation(func() error {
		for _, paint := range s.paints[s.state] {
			if err := paint.Paint(s.drawer); err != nil {
				return errors.Wrap(err, "failed to paint")
			}
		}
		return nil
	})

	if err != nil {
		return errors.Wrap(err, "failed to present")
	}

	return nil
}

// Clear removes or destroys all allocated objects to free memory
func (s Scene) Clear() {
	defer s.clearSdl2()

	for _, objects := range s.paints {
		for _, paint := range objects {
			switch p := paint.(type) {
			case object.Destroyable:
				p.Destroy()
			}
		}
	}
}

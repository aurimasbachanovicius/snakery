package scene

import (
	"fmt"
	"github.com/3auris/snakery/internal/object"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"os"
	"time"
)

// Scene holds paints and state of the current game
type Scene struct {
	r *sdl.Renderer
	w *sdl.Window

	state  object.GameState
	paints map[object.GameState][]object.Paintable
}

// New create new Scene with given parameters
func New(fontPath string, screenWidth, screenHeight int32) (*Scene, error) {
	w, r, err := prepareSdl2(int32(screenWidth), int32(screenHeight))
	if err != nil {
		return nil, fmt.Errorf("could not prepare sdl2: %v", err)
	}

	font, err := ttf.OpenFont(fontPath, 100)
	if err != nil {
		return nil, fmt.Errorf("could not load font: %v", err)
	}

	scrn := object.GameScreen{W: screenWidth, H: screenHeight}

	apple := object.NewApple()
	score := object.NewScore(font)
	snake := object.NewSnake(apple, score, font, scrn)
	deadScreen := object.DeadScreen{Score: score, Font: *font, Screen: scrn}

	return &Scene{
		r: r,
		w: w,

		state: object.SnakeRunning,
		paints: map[object.GameState][]object.Paintable{
			object.SnakeRunning: {snake, apple, score},
			object.DeadSnake:    {deadScreen},
		},
	}, nil
}

// Run runs goroutine and updates all paints and listening of events
func (s Scene) Run(events <-chan sdl.Event) <-chan error {
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

				if done := s.handleEvent(e); done {
					os.Exit(0)
					return
				}
			case <-ticker:
				s.update()

				if err := s.paint(); err != nil {
					errc <- err
				}
			}
		}

		os.Exit(0)
	}()

	return errc
}

func (s *Scene) handleEvent(event sdl.Event) bool {
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

func (s *Scene) update() {
	for _, paint := range s.paints[s.state] {
		switch p := paint.(type) {
		case object.Updateable:
			state := p.Update()
			if state != s.state {
				s.state = state
				return
			}
		}
	}
}

func (s Scene) paint() error {
	s.r.Clear()

	s.r.SetDrawColor(255, 255, 255, 255)
	s.r.FillRect(nil)

	for _, paint := range s.paints[s.state] {
		if err := paint.Paint(s.r); err != nil {
			return fmt.Errorf("failed to paint: %v", err)
		}
	}

	s.r.Present()

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

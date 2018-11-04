package main

import (
	"fmt"
	"github.com/3auris/snakery/apple"
	"github.com/3auris/snakery/score"
	"github.com/3auris/snakery/snake"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"os"
	"strconv"
	"time"
)

type scene struct {
	r *sdl.Renderer
	f *ttf.Font

	snake *snake.Snake
	apple *apple.Apple
	score *score.Score
}

func newScene(r *sdl.Renderer) (*scene, error) {
	font, err := ttf.OpenFont("assets/ubuntu.ttf", 100)
	if err != nil {
		return nil, fmt.Errorf("could not load font: %v", err)
	}

	return &scene{
		r: r,
		f: font,

		snake: snake.New(),
		apple: apple.New(),
		score: score.New(font),
	}, nil
}

func (s scene) run(events <-chan sdl.Event) <-chan error {
	errc := make(chan error)
	go func() {
		frame := 0
		tick := time.Tick(time.Millisecond)
		start := time.Now()

		for {
			select {
			case e := <-events:
				if done := s.handleEvent(e, start); done {
					os.Exit(0)
					return
				}
			case <-tick:
				frame++
				if !s.snake.IsTimeUpdate(frame) {
					continue
				}

				if s.snake.IsDead() {
					s.paintDead(s.score.Amount())
					return
				}

				s.update()

				if err := s.paint(); err != nil {
					errc <- err
				}
			}
		}

		fmt.Printf("%s\n", time.Since(start))
		os.Exit(0)
	}()

	return errc
}

func (s scene) update() {
	s.apple.Update()
	s.snake.Update()
	s.snake.Eat(s.apple, s.score)
	s.snake.TouchDeadZone()
}

func (s *scene) handleEvent(event sdl.Event, t time.Time) bool {
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
		case sdl.K_LEFT, sdl.K_a:
			s.snake.ChangeVector(snake.Left)
		case sdl.K_RIGHT, sdl.K_d:
			s.snake.ChangeVector(snake.Right)
		case sdl.K_UP, sdl.K_w:
			s.snake.ChangeVector(snake.Up)
		case sdl.K_DOWN, sdl.K_s:
			s.snake.ChangeVector(snake.Down)
		}
	}
	return false
}

func (s scene) paint() error {
	s.r.Clear()

	s.r.SetDrawColor(255, 255, 255, 255)
	s.r.FillRect(nil)

	if err := s.snake.Paint(s.r); err != nil {
		return fmt.Errorf("failed to paint snake: %v", err)
	}

	if err := s.apple.Paint(s.r); err != nil {
		return fmt.Errorf("failed to paint apple: %v", err)
	}

	if err := s.score.Paint(s.r); err != nil {
		return fmt.Errorf("failed to paint score: %v", err)
	}

	s.r.Present()

	return nil
}

func (s scene) destroy() {
	s.score.Destroy()
}

func (s scene) paintDead(score int) error {
	s.r.Clear()

	s.r.SetDrawColor(0, 0, 0, 0)
	s.r.FillRect(nil)

	sAmount := strconv.Itoa(score)

	c := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	sf, err := s.f.RenderUTF8Solid("Final score: "+sAmount, c)
	if err != nil {
		return fmt.Errorf("could not render title: %v", err)
	}
	defer sf.Free()

	t, err := s.r.CreateTextureFromSurface(sf)
	if err != nil {
		return fmt.Errorf("could not create texture: %v", err)
	}
	defer t.Destroy()

	rect := &sdl.Rect{X: 10, Y: 100, W: 490, H: 60}

	if err := s.r.Copy(t, nil, rect); err != nil {
		return fmt.Errorf("could not copy texture: %v", err)
	}

	s.r.Present()

	return nil
}

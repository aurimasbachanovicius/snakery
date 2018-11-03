package main

import (
	"fmt"
	"github.com/3auris/snakery/apple"
	"github.com/3auris/snakery/snake"
	"github.com/veandco/go-sdl2/sdl"
	"os"
	"time"
)

type scene struct {
	r     *sdl.Renderer
	snake *snake.Snake
	apple *apple.Apple
}

func newScene(r *sdl.Renderer) *scene {
	return &scene{
		r:     r,
		snake: snake.New(),
		apple: apple.New(),
	}
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
					fmt.Println("You're dead")
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
	s.snake.Eat(s.apple)
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

	s.r.Present()

	return nil
}

package main

import (
	"fmt"
	"github.com/3auris/snakery/scene"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"os"
	"runtime"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}
}

func run() error {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return fmt.Errorf("could not init sdl: %v", err)
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		return fmt.Errorf("could not init ttf: %v", err)
	}
	defer ttf.Quit()

	w, r, err := sdl.CreateWindowAndRenderer(500, 500, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("could not create window: %v", err)
	}
	defer w.Destroy()

	if err != nil {
		return fmt.Errorf("failed to get surface: %v", err)
	}

	s, err := scene.New(r, "res/ubuntu.ttf")
	if err != nil {
		return fmt.Errorf("could not create scene, %v", err)
	}
	defer s.Clear()

	events := make(chan sdl.Event)
	errc := s.Run(events)

	runtime.LockOSThread()
	for {
		select {
		case events <- sdl.WaitEvent():
		case err := <-errc:
			return err
		}

	}

	return nil
}

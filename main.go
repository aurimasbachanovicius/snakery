package main

import (
	"fmt"
	"github.com/3auris/snakery/internal/scene"
	"github.com/veandco/go-sdl2/sdl"
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
	s, err := scene.New("res/ubuntu.ttf", 500, 500)
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

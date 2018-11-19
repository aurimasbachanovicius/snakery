package main

import (
	"fmt"
	"github.com/3auris/snakery/internal/scene"
	"github.com/3auris/snakery/pkg/grafio"
	"github.com/pkg/errors"
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
	r, destroy, err := scene.PrepareSdl2(500, 500)
	if err != nil {
		return errors.Wrap(err, "could not prepare sdl2 engine")
	}
	defer destroy()

	drawer, err := grafio.NewSdl2Draw(r, 500, 500)
	free, err := drawer.LoadResources("res/fonts", "res/textures")
	if err != nil {
		return errors.Wrap(err, "could not load resources")
	}
	defer free()

	s, err := scene.New(drawer)
	if err != nil {
		return fmt.Errorf("could not create scene: %v", err)
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

package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"

	"github.com/3auris/snakery/internal/object"
	"github.com/3auris/snakery/internal/scene"
	"github.com/3auris/snakery/pkg/grafio"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}
}

func run() (erro error) {
	r, destroy, err := scene.PrepareSdl2(500, 500)
	if err != nil {
		return errors.Wrap(err, "could not prepare sdl2 engine")
	}
	defer destroy()

	drawer, err := grafio.NewSdl2Draw(r, object.FontUbuntu, 500, 500)

	free, err := drawer.LoadResources("res/fonts", "res/textures")
	if err != nil {
		return errors.Wrap(err, "could not load resources")
	}
	defer func() {
		if err := free(); err != nil {
			erro = errors.Wrap(err, "could not free resources")
		}
	}()

	s, err := scene.New(drawer)
	if err != nil {
		return fmt.Errorf("could not create scene: %v", err)
	}

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
}

package grafio

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func prepareSdl2(width, height int32) (r *sdl.Renderer, destroy func() error, err error) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return nil, nil, fmt.Errorf("could not init sdl: %v", err)
	}

	if err := ttf.Init(); err != nil {
		return nil, nil, fmt.Errorf("could not init ttf: %v", err)
	}

	w, r, err := sdl.CreateWindowAndRenderer(width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create window: %v", err)
	}

	if err != nil {
		return nil, nil, fmt.Errorf("failed to get surface: %v", err)
	}

	destroy = func() error {
		if err := r.Destroy(); err != nil {
			return errors.Wrap(err, "could not destroy renderer")
		}

		if err := w.Destroy(); err != nil {
			return errors.Wrap(err, "could not destroy window")
		}

		ttf.Quit()
		sdl.Quit()

		return nil
	}

	return r, destroy, nil
}

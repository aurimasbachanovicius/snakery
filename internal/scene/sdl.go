package scene

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// PrepareSdl2 prepares the sdl2 window.
func PrepareSdl2(width, height int32) (renderer *sdl.Renderer, destroy func(), erro error) {
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

	destroy = func() {
		if err := r.Destroy(); err != nil {
			erro = errors.Wrap(err, "could not destroy renderer")
			return
		}

		if err := w.Destroy(); err != nil {
			erro = errors.Wrap(err, "could not destroy window")
			return
		}

		ttf.Quit()
		sdl.Quit()
	}

	return r, destroy, nil
}

package scene

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func PrepareSdl2(width, height int32) (*sdl.Renderer, func(), error) {
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

	destroy := func() {
		r.Destroy()
		w.Destroy()
		ttf.Quit()
		sdl.Quit()
	}

	return r, destroy, nil
}

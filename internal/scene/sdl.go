package scene

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func prepareSdl2(width, height int32) (*sdl.Window, *sdl.Renderer, error) {
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

	return w, r, nil
}

func (s Scene) clearSdl2() {
	s.r.Destroy()
	s.w.Destroy()
	ttf.Quit()
	sdl.Quit()
}

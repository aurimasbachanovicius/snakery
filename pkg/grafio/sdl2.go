package grafio

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"io/ioutil"
)

type Sdl2Draw struct {
	r *sdl.Renderer
	f *ttf.Font

	textures []*sdl.Texture

	w, h int32
}

func (s Sdl2Draw) ScreenHeight() int32 {
	return s.h
}

func (s Sdl2Draw) ScreenWidth() int32 {
	return s.w
}

func NewSdl2Draw(r *sdl.Renderer, w, h int32) (*Sdl2Draw, error) {
	return &Sdl2Draw{
		r: r,
		//f: font,

		w: w,
		h: h,
	}, nil
}

func (s *Sdl2Draw) Background(r, g, b, a uint8) error {
	if err := s.r.SetDrawColor(r, g, b, a); err != nil {
		return errors.Wrap(err, "couldn't set color")
	}

	if err := s.r.FillRect(nil); err != nil {
		return errors.Wrap(err, "couldn't fill rect")
	}

	return nil
}

func (s *Sdl2Draw) Text(txt string, opts TextOpts) error {
	c := sdl.Color{R: opts.R, G: opts.G, B: opts.B, A: opts.A}
	surface, err := s.f.RenderUTF8Solid(txt, c)
	if err != nil {
		return errors.Wrap(err, "could not render title")
	}
	defer surface.Free()

	texture, err := s.r.CreateTextureFromSurface(surface)
	if err != nil {
		return errors.Wrap(err, "could not create texture")
	}
	defer texture.Destroy()

	rect := &sdl.Rect{
		X: sizeCal(s.w, opts.XCof),
		Y: sizeCal(s.h, opts.YCof),
		W: sizeCal(s.w, .90), // todo calculate size
		H: sizeCal(s.h, .10), // todo calculate size
	}

	if err := s.r.Copy(texture, nil, rect); err != nil {
		return errors.Wrap(err, "could not copy texture")
	}

	return nil
}

func (s *Sdl2Draw) Present(f func() error) error {
	if err := s.r.Clear(); err != nil {
		return errors.Wrap(err, "could not clear the renderer")
	}

	if err := s.Background(255, 255, 255, 255); err != nil {
		return errors.Wrap(err, "could not draw background")
	}

	if err := f(); err != nil {
		return errors.Wrap(err, "could not execute f (from parameter) function")
	}

	s.r.Present()

	return nil
}

func (s *Sdl2Draw) LoadResources(fontsPath, texturesPath string) (func() error, error) {
	files, err := ioutil.ReadDir(texturesPath)
	if err != nil {
		return nil, errors.Wrap(err, "could not read dir")
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		image, err := img.Load(texturesPath + "/" + f.Name())
		if err != nil {
			return nil, fmt.Errorf("Failed to create texture: %v\n", err)
		}

		t, err := s.r.CreateTextureFromSurface(image)
		if err != nil {
			return nil, fmt.Errorf("Failed to create texture: %v\n", err)
		}
		image.Free()

		s.textures = append(s.textures, t)
	}

	return func() error {
		for _, texture := range s.textures {
			if err := texture.Destroy(); err != nil {
				return errors.Wrap(err, "could not destroy texture")
			}
		}
		return nil
	}, nil
}

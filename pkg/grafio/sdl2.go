package grafio

import (
	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Sdl2Draw struct {
	r *sdl.Renderer
	f *ttf.Font

	w, h int32
}

func NewSdl2Draw(r *sdl.Renderer, fontPath string, w, h int32) (*Sdl2Draw, error) {
	font, err := ttf.OpenFont(fontPath, 100)
	if err != nil {
		return nil, errors.Wrap(err, "could not load font")
	}

	return &Sdl2Draw{
		r: r,
		f: font,

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

func (s *Sdl2Draw) Presentation(f func() error) error {
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

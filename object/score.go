package object

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"strconv"
	"sync"
)

type score struct {
	mu sync.RWMutex

	amount int
	font   *ttf.Font
}

func NewScore(f *ttf.Font) *score {
	return &score{amount: 0, font: f}
}

func (s score) Paint(r *sdl.Renderer) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sAmount := strconv.Itoa(s.amount)

	c := sdl.Color{R: 0, G: 0, B: 0, A: 0}
	sf, err := s.font.RenderUTF8Solid(sAmount, c)
	if err != nil {
		return fmt.Errorf("could not render title: %v", err)
	}
	defer sf.Free()

	t, err := r.CreateTextureFromSurface(sf)
	if err != nil {
		return fmt.Errorf("could not create texture: %v", err)
	}
	defer t.Destroy()

	width := int32(len(sAmount)) * 16
	rect := &sdl.Rect{X: 500 - width - 10, Y: 10, W: 0 + width, H: 20}

	if err := r.Copy(t, nil, rect); err != nil {
		return fmt.Errorf("could not copy texture: %v", err)
	}

	return nil
}

func (s *score) Increase() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.amount++
}

func (s score) Destroy() {
	s.font.Close()
}

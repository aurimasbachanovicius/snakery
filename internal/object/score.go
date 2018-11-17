package object

import (
	"fmt"
	"github.com/3auris/snakery/pkg/grafio"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"strconv"
	"sync"
)

// Score the game pbject
type Score struct {
	mu sync.RWMutex

	amount int
	font   *ttf.Font
}

// NewScore creates Score with default and given values
func NewScore(f *ttf.Font) *Score {
	return &Score{amount: 0, font: f}
}

// Paint the score number to renderer to the corner
func (s Score) Paint(r grafio.Drawer) error {
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

// Destroy destroys the opened font
func (s Score) Destroy() {
	s.font.Close()
}

// Increase increased the amount of Score
func (s *Score) Increase() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.amount++
}

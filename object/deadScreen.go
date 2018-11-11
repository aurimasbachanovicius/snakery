package object

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"strconv"
)

// DeadScreen object of game which can be painted
type DeadScreen struct {
	Score  *Score
	Font   ttf.Font
	Screen GameScreen
}

// Paint paints text and Score to renderer
func (ds DeadScreen) Paint(r *sdl.Renderer) error {
	r.SetDrawColor(0, 0, 0, 0)
	r.FillRect(nil)

	sAmount := strconv.Itoa(ds.Score.amount)

	c := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	sf, err := ds.Font.RenderUTF8Solid("Final score: "+sAmount, c)
	if err != nil {
		return fmt.Errorf("could not render title: %v", err)
	}
	defer sf.Free()

	t, err := r.CreateTextureFromSurface(sf)
	if err != nil {
		return fmt.Errorf("could not create texture: %v", err)
	}
	defer t.Destroy()

	rect := &sdl.Rect{X: 10, Y: int32(ds.Screen.H / 5), W: ds.Screen.W - 10, H: int32(ds.Screen.H / 5)}

	if err := r.Copy(t, nil, rect); err != nil {
		return fmt.Errorf("could not copy texture: %v", err)
	}

	return nil
}

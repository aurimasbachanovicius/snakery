package object

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type WelcomeText struct {
	Font   ttf.Font
	Screen GameScreen

	changeState bool
}

func (wt *WelcomeText) HandleEvent(event sdl.Event) {
	switch ev := event.(type) {
	case *sdl.KeyboardEvent:
		if ev.State != sdl.PRESSED {
			break
		}

		switch event.(*sdl.KeyboardEvent).Keysym.Sym {
		case sdl.K_RETURN:
			wt.changeState = true
		}
	}
}

func (wt *WelcomeText) Update() GameState {
	if wt.changeState {
		wt.changeState = false
		return SnakeRunning
	}

	return MenuScreen
}

// Paint paints text and Score to renderer
func (wt WelcomeText) Paint(r *sdl.Renderer) error {
	r.SetDrawColor(255, 255, 255, 0)
	r.FillRect(nil)

	c := sdl.Color{R: 34, G: 139, B: 34, A: 10}
	welcomeSf, err := wt.Font.RenderUTF8Solid("Welcome to the snake game", c)
	if err != nil {
		return fmt.Errorf("could not render title: %v", err)
	}
	defer welcomeSf.Free()

	startSf, err := wt.Font.RenderUTF8Solid("Press (Enter) to start the game", c)
	if err != nil {
		return fmt.Errorf("could not render title: %v", err)
	}
	defer startSf.Free()

	welcomeT, err := r.CreateTextureFromSurface(welcomeSf)
	if err != nil {
		return fmt.Errorf("could not create texture: %v", err)
	}
	defer welcomeT.Destroy()

	startT, err := r.CreateTextureFromSurface(startSf)
	if err != nil {
		return fmt.Errorf("could not create texture: %v", err)
	}
	defer startT.Destroy()

	welcomeRect := &sdl.Rect{
		X: size(wt.Screen.W, .05),
		Y: size(wt.Screen.H, .15),
		W: size(wt.Screen.W, .90),
		H: size(wt.Screen.H, .10),
	}

	startRect := &sdl.Rect{
		X: size(wt.Screen.W, .05),
		Y: size(wt.Screen.H, .30),
		W: size(wt.Screen.W, .90),
		H: size(wt.Screen.H, .10),
	}

	if err := r.Copy(welcomeT, nil, welcomeRect); err != nil {
		return fmt.Errorf("could not copy texture: %v", err)
	}

	if err := r.Copy(startT, nil, startRect); err != nil {
		return fmt.Errorf("could not copy texture: %v", err)
	}

	return nil
}

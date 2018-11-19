package object

import (
	"fmt"
	"github.com/3auris/snakery/pkg/grafio"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"strconv"
)

// DeadScreen object of game which can be painted
type DeadScreen struct {
	Score  *Score
	Font   ttf.Font
	Screen GameScreen

	toMenu bool
}

func (ds *DeadScreen) HandleEvent(event sdl.Event) {
	switch ev := event.(type) {
	case *sdl.KeyboardEvent:
		if ev.State != sdl.PRESSED {
			break
		}

		switch event.(*sdl.KeyboardEvent).Keysym.Sym {
		case sdl.K_RETURN:
			ds.toMenu = true
		}
	}
}

func (ds *DeadScreen) Update() GameState {
	if ds.toMenu {
		ds.toMenu = false
		return MenuScreen
	}

	return DeadSnake
}

// Paint paints text and Score to renderer
func (ds DeadScreen) Paint(d grafio.Drawer) error {
	//r.SetDrawColor(0, 0, 0, 0)
	//r.FillRect(nil)

	sAmount := strconv.Itoa(ds.Score.amount)

	c := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	scoreSf, err := ds.Font.RenderUTF8Solid("Final score: "+sAmount, c)
	if err != nil {
		return fmt.Errorf("could not render title: %v", err)
	}
	defer scoreSf.Free()

	restartSf, err := ds.Font.RenderUTF8Solid("Press (Enter) to go to menu", c)
	if err != nil {
		return fmt.Errorf("could not render title: %v", err)
	}
	defer restartSf.Free()

	//scoreT, err := r.CreateTextureFromSurface(scoreSf)
	if err != nil {
		return fmt.Errorf("could not create texture: %v", err)
	}
	//defer scoreT.Destroy()

	//restartT, err := r.CreateTextureFromSurface(restartSf)
	if err != nil {
		return fmt.Errorf("could not create texture: %v", err)
	}
	//defer restartT.Destroy()

	//scoreRect := &sdl.Rect{
	//	X: size(ds.Screen.W, .05),
	//	Y: size(ds.Screen.H, .15),
	//	W: size(ds.Screen.W, .90),
	//	H: size(ds.Screen.H, .20),
	//}
	//
	//restartRect := &sdl.Rect{
	//	X: size(ds.Screen.W, .05),
	//	Y: size(ds.Screen.H, .40),
	//	W: size(ds.Screen.W, .90),
	//	H: size(ds.Screen.H, .10),
	//}

	//if err := r.Copy(scoreT, nil, scoreRect); err != nil {
	//	return fmt.Errorf("could not copy texture: %v", err)
	//}
	//
	//if err := r.Copy(restartT, nil, restartRect); err != nil {
	//	return fmt.Errorf("could not copy texture: %v", err)
	//}
	return nil
}

package object

import (
	"github.com/veandco/go-sdl2/sdl"
)

// GameScreen have struct of width and height of screen
type GameScreen struct {
	W, H int32
}

// GameState game current action
type GameState int

const (
	// SnakeRunning state when the Snake is running/moving and it's not the end yet
	SnakeRunning GameState = 1

	// DeadSnake state when Snake touches something and dies and need to show dead screen
	DeadSnake GameState = 2

	// MenuScreen state when There's menu shown for setting and entering into game
	MenuScreen GameState = 3
)

// Paintable paints something to sdl renderer
type Paintable interface {
	Paint(r *sdl.Renderer) error
}

// Updateable object data can or should be updated every each frame with certain information in function
type Updateable interface {
	Update() GameState
}

// Destroyable something created in Update() and can be Destroyed after exit of program to save memory
type Destroyable interface {
	Destroy()
}

// Handleable it can handle input from sdl events
type Handleable interface {
	HandleEvent(event sdl.Event)
}

func size(size int32, cof float32) int32 {
	return int32(float32(size) * (float32(cof)))
}

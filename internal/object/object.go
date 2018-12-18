package object

import (
	"github.com/veandco/go-sdl2/sdl"

	"github.com/3auris/snakery/pkg/grafio"
)

const (
	textureApple string = "apple.png"

	// FontUbuntu filename of ubuntu font in the resources
	FontUbuntu string = "ubuntu.ttf"
)

// GameScreen have struct of width and height of screen
type GameScreen struct {
	W, H int32
}

// GameState game current action
type GameState int

const (
	// SnakeRunning state when the Snake is running/moving and it's not the end yet
	SnakeRunning GameState = iota

	// DeadSnake state when Snake touches something and dies and need to show dead screen
	DeadSnake

	// MenuScreen state when There's menu shown for setting and entering into game
	MenuScreen
)

// Paintable paints something to sdl renderer
type Paintable interface {
	Paint(d grafio.Drawer) error
}

// Updateable object data can or should be updated every each frame with certain information in function
type Updateable interface {
	Update() GameState
}

// Handleable it can handle input from sdl events
type Handleable interface {
	HandleEvent(event sdl.Event)
}

package object

import (
	"github.com/veandco/go-sdl2/sdl"
)

type GameState int

const (
	SnakeRunning GameState = 1
	DeadSnake    GameState = 2
)

type Paintable interface {
	Paint(r *sdl.Renderer) error
}

type Updateable interface {
	Update() GameState
}

type Destroyable interface {
	Destroy()
}

type Handleable interface {
	HandleEvent(event sdl.Event)
}

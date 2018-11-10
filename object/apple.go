package object

import (
	"github.com/3auris/snakery/pkg/geometrio"
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
	"sync"
	"time"
)

type apple struct {
	mu sync.RWMutex

	x, y  int32
	size  int32
	eaten bool
}

func NewApple() *apple {
	return &apple{eaten: true, size: 10}
}

func (a *apple) Update() GameState {
	a.mu.Lock()
	defer a.mu.Unlock()

	if ! a.eaten {
		return SnakeRunning
	}

	a.eaten = false

	rand.Seed(time.Now().UnixNano())
	rX := rand.Intn(460-1) + 1

	rand.Seed(time.Now().UnixNano())
	rY := rand.Intn(460-1) + 1

	a.x = int32(rX)
	a.y = int32(rY)

	return SnakeRunning
}

func (a apple) Paint(r *sdl.Renderer) error {
	a.mu.RLock()
	defer a.mu.RUnlock()

	r.SetDrawColor(255, 0, 0, 0)

	if err := r.FillRect(&sdl.Rect{X: a.x, Y: a.y, W: a.size, H: a.size}); err != nil {
		return err
	}

	return nil
}

func (a apple) ExistsIn(pl, pr geometrio.Cord) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.eaten == true {
		return false
	}

	l := geometrio.Cord{X: a.x, Y: a.y}
	r := geometrio.Cord{
		X: a.x + a.size,
		Y: a.y + a.size,
	}

	return geometrio.IsOverlapping(l, r, pl, pr)
}

func (a *apple) EatApple() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.eaten = true
}

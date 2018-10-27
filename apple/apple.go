package apple

import (
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
	"sync"
	"time"
)

type Apple struct {
	mu sync.RWMutex

	x, y  int32
	size  int32
	eaten bool
}

func New() *Apple {
	return &Apple{eaten: true, size: 10}
}

func (a *Apple) Update() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.eaten == true {
		rand.Seed(time.Now().Unix())
		rX := rand.Intn(460-1) + 1

		rand.Seed(time.Now().Unix())
		rY := rand.Intn(460-1) + 1

		a.x = int32(rX)
		a.y = int32(rY)

		a.eaten = false
	}

}

func (a Apple) Paint(r *sdl.Renderer) error {
	a.mu.RLock()
	defer a.mu.RUnlock()

	r.SetDrawColor(255, 0, 0, 0)

	if err := r.FillRect(&sdl.Rect{X: a.x, Y: a.y, W: a.size, H: a.size}); err != nil {
		return err
	}

	return nil
}

func (a Apple) ExistsIn(x, y, w, h int32) bool {
	if a.x < x+w && a.x+a.size > x &&
		a.y < y+h && a.y+a.size > y {
		return true
	}

	return false
}

func (a *Apple) EatApple() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.eaten = true
}

package object

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/3auris/snakery/pkg/geometrio"
	"github.com/3auris/snakery/pkg/grafio"
)

// Apple is the game object
type Apple struct {
	mu *sync.RWMutex

	x, y  int32
	size  int32
	eaten bool
}

// NewApple creates new Apple with default values
func NewApple() *Apple {
	x, y := getAppleRandXY()

	return &Apple{mu: &sync.RWMutex{}, eaten: false, size: 16, x: int32(x), y: int32(y)}
}

func getAppleRandXY() (x, y int) {
	rand.Seed(time.Now().UnixNano())

	x = rand.Intn(460-1) + 1
	y = rand.Intn(460-1) + 1

	return x, y
}

// Update check is Apple is eaten and changes the state of Apple coordinates
func (a *Apple) Update() GameState {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.eaten {
		return SnakeRunning
	}

	a.eaten = false

	x, y := getAppleRandXY()

	a.x = int32(x)
	a.y = int32(y)

	return SnakeRunning
}

// Paint paints Apple to the given renderer
func (a Apple) Paint(d grafio.Drawer) error {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if err := d.TextureRect(a.x, a.y, a.size, a.size, textureApple); err != nil {
		return fmt.Errorf("could not paint Apple: %v", err)
	}

	return nil
}

func (a Apple) existsIn(pl, pr geometrio.Cord) bool {
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

func (a *Apple) eatApple() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.eaten = true
}

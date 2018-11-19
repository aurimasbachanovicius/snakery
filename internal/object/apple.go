package object

import (
	"fmt"
	"github.com/3auris/snakery/pkg/geometrio"
	"github.com/3auris/snakery/pkg/grafio"
	"math/rand"
	"sync"
	"time"
)

// Apple is the game object
type Apple struct {
	mu sync.RWMutex

	x, y  int32
	size  int32
	eaten bool
}

// NewApple creates new Apple with default values
func NewApple() *Apple {
	return &Apple{eaten: true, size: 16}
}

// Update check is apple is eaten and changes the state of apple coordinates
func (a *Apple) Update() GameState {
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

// Paint paints apple to the given renderer
func (a Apple) Paint(d grafio.Drawer) error {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if err := d.TextureRect(a.x, a.y, a.size, a.size, TextureApple); err != nil {
		return fmt.Errorf("could not paint apple: %v", err)
	}

	return nil
}

// ExistsIn check is the Apple exists in rectangle between pl and pr coordinates
func (a Apple) ExistsIn(pl, pr geometrio.Cord) bool {
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

// EatApple set state of eaten to true
func (a *Apple) EatApple() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.eaten = true
}

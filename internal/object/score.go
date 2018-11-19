package object

import (
	"github.com/3auris/snakery/pkg/grafio"
	"github.com/pkg/errors"
	"strconv"
	"sync"
)

// Score the game pbject
type Score struct {
	mu sync.RWMutex

	amount int
}

// NewScore creates Score with default and given values
func NewScore() *Score {
	return &Score{amount: 0}
}

// Paint the score number to renderer to the corner
func (s Score) Paint(d grafio.Drawer) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	opts := grafio.TextOpts{Size: 2, XCof: .05, YCof: .90, Color: grafio.RGBA{R: 0, G: 0, B: 0, A: 0}}
	if err := d.Text(strconv.Itoa(s.amount), opts); err != nil {
		return errors.Wrap(err, "failed to draw the score")
	}

	return nil
}

// Increase increased the amount of Score
func (s *Score) Increase() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.amount++
}

package object

import (
	"github.com/3auris/snakery/pkg/grafio"
	"github.com/pkg/errors"
	"strconv"
	"sync"
)

// Score the game object
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

	opts := grafio.TextOpts{Size: 13, XCof: .94, YCof: .01, Color: grafio.ColorBlack, Align: grafio.Right}
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

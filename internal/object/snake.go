package object

import (
	"fmt"
	"github.com/3auris/snakery/pkg/geometrio"
	"github.com/3auris/snakery/pkg/grafio"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"math"
	"sync"
)

type snakeVector int

const (
	up    snakeVector = 1
	down  snakeVector = 2
	left  snakeVector = 3
	right snakeVector = 4

	stepSize int32 = 10
)

type part struct {
	x, y   int32
	w, h   int32
	vector snakeVector
}

// Snake game object
type Snake struct {
	mu sync.RWMutex

	parts []*part
	size  int

	lockMove bool

	apple  *Apple
	score  *Score
	font   *ttf.Font
	screen GameScreen
}

// NewSnake create Snake struct with default and given values
func NewSnake(a *Apple, s *Score, f *ttf.Font, scr GameScreen) *Snake {
	return &Snake{
		parts:    []*part{{x: 50, y: 50, w: 120, h: stepSize, vector: right}},
		size:     120,
		lockMove: false,

		apple:  a,
		score:  s,
		font:   f,
		screen: scr,
	}
}

// HandleEvent handles the movement of snake
func (s *Snake) HandleEvent(event sdl.Event) {
	switch ev := event.(type) {
	case *sdl.KeyboardEvent:
		if ev.State != sdl.PRESSED {
			break
		}

		switch event.(*sdl.KeyboardEvent).Keysym.Sym {
		case sdl.K_LEFT, sdl.K_a:
			s.changeVector(left)
		case sdl.K_RIGHT, sdl.K_d:
			s.changeVector(right)
		case sdl.K_UP, sdl.K_w:
			s.changeVector(up)
		case sdl.K_DOWN, sdl.K_s:
			s.changeVector(down)
		}
	}
}

func (s *Snake) reset() {
	n := NewSnake(s.apple, s.score, s.font, s.screen)
	s.score.amount = 0
	s.score = n.score
	s.size = n.size
	s.parts = n.parts
}

// Update updates the snake ar gives the GameState
func (s *Snake) Update() GameState {
	if s.touchDeadZone() {
		return DeadSnake
	}

	s.eat()
	s.move()

	return SnakeRunning
}

// Paint paints snake to renderer
func (s *Snake) Paint(d grafio.Drawer) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	r.SetDrawColor(34, 139, 34, 10)

	for _, part := range s.parts {
		if err := r.FillRect(&sdl.Rect{X: part.x, Y: part.y, W: part.w, H: part.h}); err != nil {
			return fmt.Errorf("failed to fill rect with the Snake part: %v", err)
		}
	}

	return nil
}

func (s *Snake) getWholeSize() int {
	sum := 0

	for _, part := range s.parts {
		switch part.vector {
		case up, down:
			sum += int(math.Abs(float64(part.h)))
		case left, right:
			sum += int(math.Abs(float64(part.w)))
		}
	}

	return sum
}

func (s *Snake) eat() {
	s.mu.RLock()
	latest := s.latestPart()
	exists := s.apple.ExistsIn(latest.getCords())
	s.mu.RUnlock()

	if exists {
		s.apple.EatApple()
		s.score.Increase()

		s.mu.Lock()
		s.size += 50
		s.mu.Unlock()
	}
}

func (s *Snake) move() {
	s.lockMove = false

	first := s.parts[0]
	if s.getWholeSize() > s.size && len(s.parts) > 1 {
		if first.h == 0 || first.w == 0 {
			s.parts = append(s.parts[:0], s.parts[1:]...)

			first = s.parts[0]
		}

		switch first.vector {
		case up:
			first.h += stepSize
			first.y -= stepSize
		case down:
			first.h -= stepSize
			first.y += stepSize
		case left:
			first.x -= stepSize
			first.w += stepSize
		case right:
			first.x += stepSize
			first.w -= stepSize
		}
	}

	latest := s.latestPart()
	switch latest.vector {
	case up:
		if s.getWholeSize() > s.size {
			latest.y -= stepSize
		} else {
			latest.h -= stepSize
		}
	case down:
		if s.getWholeSize() > s.size {
			latest.y += stepSize
		} else {
			latest.h += stepSize
		}
	case left:
		if s.getWholeSize() > s.size {
			latest.x -= stepSize
		} else {
			latest.w -= stepSize
		}
	case right:
		if s.getWholeSize() > s.size {
			latest.x += stepSize
		} else {
			latest.w += stepSize
		}
	}
}

func (s *Snake) touchDeadZone() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	head := s.latestPart()
	sl, sr := head.getCords()

	if sl.X > s.screen.W ||
		sl.Y > s.screen.H ||
		sr.X > s.screen.W ||
		sr.Y > s.screen.H ||
		sl.X < 0 ||
		sl.Y < 0 ||
		sr.X < 0 ||
		sr.Y < 0 {
		return true
	}

	if len(s.parts) <= 3 {
		return false
	}

	parts := s.parts[:len(s.parts)-3]

	for _, part := range parts {
		pl, pr := part.getCords()

		if geometrio.IsOverlapping(pl, pr, sl, sr) {
			return true
		}
	}

	return false
}

func (s *Snake) latestPart() *part {
	return s.parts[len(s.parts)-1]
}

func (p part) getCords() (geometrio.Cord, geometrio.Cord) {
	return geometrio.Cord{
		X: p.x,
		Y: p.y,
	}, geometrio.Cord{
		X: p.x + p.w,
		Y: p.y + p.h,
	}
}

func (s Snake) canGo(v snakeVector) bool {
	lv := s.latestPart().vector

	if s.lockMove || v == lv {
		return false
	}

	if (v == left || v == right) && (lv == left || lv == right) {
		return false
	}

	if (v == up || v == down) && (lv == up || lv == down) {
		return false
	}

	return true
}

func (s *Snake) changeVector(v snakeVector) {
	s.mu.Lock()
	defer s.mu.Unlock()

	latest := s.latestPart()

	if !s.canGo(v) {
		return
	}

	p := part{x: latest.x + latest.w, y: latest.y + latest.h, vector: v}

	switch v {
	case up:
		p.x -= stepSize
		p.y -= stepSize
		if latest.vector == left {
			p.x += stepSize
		}
		p.w = stepSize
	case down:
		if latest.vector == right {
			p.x -= stepSize
		}

		p.w = stepSize
	case left:
		if latest.vector == down {
			p.y -= stepSize
		}
		p.x -= stepSize
		p.h = stepSize
	case right:
		if latest.vector == down {
			p.y -= stepSize
		}
		p.h = stepSize
	}

	s.parts = append(s.parts, &p)
	s.lockMove = true
}

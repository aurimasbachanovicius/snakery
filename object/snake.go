package object

import (
	"fmt"
	"github.com/3auris/snakery/pkg/geometrio"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"math"
	"sync"
)

type snakeVector int

const (
	Up    snakeVector = 1
	Down  snakeVector = 2
	Left  snakeVector = 3
	Right snakeVector = 4

	StepSize int32 = 10
)

type part struct {
	x, y   int32
	w, h   int32
	vector snakeVector
}

type snake struct {
	mu sync.RWMutex

	parts []*part
	size  int

	lockMove bool

	apple  *apple
	score  *score
	font   *ttf.Font
	screen GameScreen
}

func NewSnake(a *apple, s *score, f *ttf.Font, scr GameScreen) *snake {
	return &snake{
		parts:    []*part{{x: 50, y: 50, w: 120, h: StepSize, vector: Right}},
		size:     120,
		lockMove: false,

		apple:  a,
		score:  s,
		font:   f,
		screen: scr,
	}
}

func (s *snake) HandleEvent(event sdl.Event) {
	switch ev := event.(type) {
	case *sdl.KeyboardEvent:
		if ev.State != sdl.PRESSED {
			break
		}

		switch event.(*sdl.KeyboardEvent).Keysym.Sym {
		case sdl.K_LEFT, sdl.K_a:
			s.changeVector(Left)
		case sdl.K_RIGHT, sdl.K_d:
			s.changeVector(Right)
		case sdl.K_UP, sdl.K_w:
			s.changeVector(Up)
		case sdl.K_DOWN, sdl.K_s:
			s.changeVector(Down)
		}
	}
}

func (s *snake) Update() GameState {
	if s.touchDeadZone() {
		return DeadSnake
	}

	s.eat()
	s.move()

	return SnakeRunning
}

func (s *snake) Paint(r *sdl.Renderer) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	r.SetDrawColor(34, 139, 34, 10)

	for _, part := range s.parts {
		if err := r.FillRect(&sdl.Rect{X: part.x, Y: part.y, W: part.w, H: part.h}); err != nil {
			return fmt.Errorf("failed to fill rect with the snake part: %v", err)
		}
	}

	return nil
}

func (s *snake) getWholeSize() int {
	sum := 0

	for _, part := range s.parts {
		switch part.vector {
		case Up, Down:
			sum += int(math.Abs(float64(part.h)))
		case Left, Right:
			sum += int(math.Abs(float64(part.w)))
		}
	}

	return sum
}

func (s *snake) eat() {
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

func (s *snake) move() {
	s.lockMove = false

	first := s.parts[0]
	if s.getWholeSize() > s.size && len(s.parts) > 1 {
		if first.h == 0 || first.w == 0 {
			s.parts = append(s.parts[:0], s.parts[1:]...)

			first = s.parts[0]
		}

		switch first.vector {
		case Up:
			first.h += StepSize
			first.y -= StepSize
		case Down:
			first.h -= StepSize
			first.y += StepSize
		case Left:
			first.x -= StepSize
			first.w += StepSize
		case Right:
			first.x += StepSize
			first.w -= StepSize
		}
	}

	latest := s.latestPart()
	switch latest.vector {
	case Up:
		if s.getWholeSize() > s.size {
			latest.y -= StepSize
		} else {
			latest.h -= StepSize
		}
	case Down:
		if s.getWholeSize() > s.size {
			latest.y += StepSize
		} else {
			latest.h += StepSize
		}
	case Left:
		if s.getWholeSize() > s.size {
			latest.x -= StepSize
		} else {
			latest.w -= StepSize
		}
	case Right:
		if s.getWholeSize() > s.size {
			latest.x += StepSize
		} else {
			latest.w += StepSize
		}
	}
}

func (s *snake) touchDeadZone() bool {
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

func (s *snake) latestPart() *part {
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

func (s snake) canGo(v snakeVector) bool {
	lv := s.latestPart().vector

	if s.lockMove || v == lv {
		return false
	}

	if (v == Left || v == Right) && (lv == Left || lv == Right) {
		return false
	}

	if (v == Up || v == Down) && (lv == Up || lv == Down) {
		return false
	}

	return true
}

func (s *snake) changeVector(v snakeVector) {
	s.mu.Lock()
	defer s.mu.Unlock()

	latest := s.latestPart()

	if !s.canGo(v) {
		return
	}

	p := part{x: latest.x + latest.w, y: latest.y + latest.h, vector: v}

	switch v {
	case Up:
		p.x -= StepSize
		p.y -= StepSize
		if latest.vector == Left {
			p.x += StepSize
		}
		p.w = StepSize
	case Down:
		if latest.vector == Right {
			p.x -= StepSize
		}

		p.w = StepSize
	case Left:
		if latest.vector == Down {
			p.y -= StepSize
		}
		p.x -= StepSize
		p.h = StepSize
	case Right:
		if latest.vector == Down {
			p.y -= StepSize
		}
		p.h = StepSize
	}

	s.parts = append(s.parts, &p)
	s.lockMove = true
}

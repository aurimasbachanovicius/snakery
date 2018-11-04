package snake

import (
	"github.com/3auris/snakery/apple"
	"github.com/3auris/snakery/math"
	"github.com/3auris/snakery/score"
	"github.com/veandco/go-sdl2/sdl"
	"sync"
)

type Vector int

const (
	Up    Vector = 1
	Down  Vector = 2
	Left  Vector = 3
	Right Vector = 4

	Fat int32 = 10
)

type part struct {
	x, y   int32
	w, h   int32
	vector Vector
}

type Snake struct {
	mu sync.RWMutex

	parts []*part
	size  int
	speed float32
	step  int32

	lockMove bool

	dead bool
}

func New() *Snake {
	return &Snake{
		parts:    []*part{{x: 50, y: 50, w: 120, h: Fat, vector: Right}},
		size:     120,
		step:     10,
		speed:    2,
		lockMove: false,
		dead:     false,
	}
}

func (s Snake) canGo(v Vector, latestV Vector) bool {
	if s.lockMove {
		return false
	}

	if v == latestV {
		return false
	}

	if (v == Left || v == Right) && (latestV == Left || latestV == Right) {
		return false
	}

	if (v == Up || v == Down) && (latestV == Up || latestV == Down) {
		return false
	}

	return true
}

func (s *Snake) ChangeVector(v Vector) {
	s.mu.Lock()
	defer s.mu.Unlock()

	latest := s.parts[len(s.parts)-1]

	if !s.canGo(v, latest.vector) {
		return
	}

	p := part{x: latest.x + latest.w, y: latest.y + latest.h, vector: v}

	switch v {
	case Up:
		p.x -= Fat
		p.y -= Fat
		if latest.vector == Left {
			p.x += Fat
		}
		p.w = Fat
	case Down:
		if latest.vector == Right {
			p.x -= Fat
		}

		p.w = Fat
	case Left:
		if latest.vector == Down {
			p.y -= Fat
		}
		p.x -= Fat
		p.h = Fat
	case Right:
		if latest.vector == Down {
			p.y -= Fat
		}
		p.h = Fat
	}

	s.parts = append(s.parts, &p)
	s.lockMove = true
}

func (s *Snake) Update() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.speed > 50 {
		s.speed = 50
	}

	s.lockMove = false

	first := s.parts[0]
	if s.getWholeSize() > s.size && len(s.parts) > 1 {
		if first.h == 0 || first.w == 0 {
			s.parts = append(s.parts[:0], s.parts[1:]...)

			first = s.parts[0]
		}

		switch first.vector {
		case Up:
			first.h += s.step
			first.y -= s.step
		case Down:
			first.h -= s.step
			first.y += s.step
		case Left:
			first.x -= s.step
			first.w += s.step
		case Right:
			first.x += s.step
			first.w -= s.step
		}
	}

	latest := s.parts[len(s.parts)-1]
	switch latest.vector {
	case Up:
		if s.getWholeSize() > s.size {
			latest.y -= s.step
		} else {
			latest.h -= s.step
		}
	case Down:
		if s.getWholeSize() > s.size {
			latest.y += s.step
		} else {
			latest.h += s.step
		}
	case Left:
		if s.getWholeSize() > s.size {
			latest.x -= s.step
		} else {
			latest.w -= s.step
		}
	case Right:
		if s.getWholeSize() > s.size {
			latest.x += s.step
		} else {
			latest.w += s.step
		}
	}
}

func (s *Snake) Paint(r *sdl.Renderer) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	r.SetDrawColor(34, 139, 34, 10)

	for _, part := range s.parts {
		if err := r.FillRect(&sdl.Rect{X: part.x, Y: part.y, W: part.w, H: part.h}); err != nil {
			return err
		}
	}

	return nil
}

func (s *Snake) getWholeSize() int {
	sum := 0

	for _, part := range s.parts {
		switch part.vector {
		case Up, Down:
			sum += int(math.Abs(part.h))
		case Left, Right:
			sum += int(math.Abs(part.w))
		}
	}

	return sum
}

func (s *Snake) Eat(a *apple.Apple, score *score.Score) {
	s.mu.RLock()
	latest := s.parts[len(s.parts)-1]
	exists := a.ExistsIn(latest.getCord())
	s.mu.RUnlock()

	if exists {
		a.EatApple()
		score.Increase()

		s.mu.Lock()

		s.size += 50
		s.speed += 1

		s.mu.Unlock()
	}
}

func (s *Snake) TouchDeadZone() {
	s.mu.Lock()
	defer s.mu.Unlock()

	head := s.parts[len(s.parts)-1]
	sl, sr := head.getCord()

	if sl.X > 500 || sl.Y > 500 || sr.X > 500 || sr.Y > 500 {
		s.dead = true
		return
	}

	if sl.X < 0 || sl.Y < 0 || sr.X < 0 || sr.Y < 0 {
		s.dead = true
		return
	}

	if len(s.parts) <= 3 {
		return
	}

	parts := s.parts[:len(s.parts)-3]

	for _, part := range parts {
		pl, pr := part.getCord()

		if math.IsOverlapping(pl, pr, sl, sr) {
			s.dead = true
			return
		}
	}
}

func (p part) getCord() (math.Cord, math.Cord) {
	return math.Cord{
		X: p.x,
		Y: p.y,
	}, math.Cord{
		X: p.x + p.w,
		Y: p.y + p.h,
	}
}

func (s Snake) IsDead() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.dead
}

func (s Snake) IsTimeUpdate(f int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if f == 1 {
		return true
	}

	return f%(100-int(s.speed)) == 0
}

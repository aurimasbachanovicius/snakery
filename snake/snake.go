package snake

import (
	"github.com/3auris/snakery/apple"
	"github.com/3auris/snakery/math"
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
	size  int32
	speed float32
	step  int32
}

func New() *Snake {
	return &Snake{
		parts: []*part{{x: 50, y: 50, w: 120, h: Fat, vector: Right}},
		size:  120,
		step:  2,
	}
}

func (Snake) canGo(v Vector, latestV Vector) bool {
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
}

func (s *Snake) Update(frame int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.speed < 2 {
		s.speed = 2
	}

	if frame%20-int(s.speed) != 0 {
		return
	}

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

func (s *Snake) getWholeSize() int32 {
	var sum int32 = 0

	for _, part := range s.parts {
		switch part.vector {
		case Up, Down:
			sum += math.Abs(part.h)
		case Left, Right:
			sum += math.Abs(part.w)
		}
	}

	return sum
}

func (s *Snake) Touch(a *apple.Apple) {
	latest := s.parts[len(s.parts)-1]

	if a.ExistsIn(latest.x, latest.y, latest.w, latest.h) {
		s.mu.Lock()
		s.size += 50
		s.speed += 0.2
		s.mu.Unlock()

		a.EatApple()
	}
}

package geometrio

// IsOverlapping checks is one object is overlapping with another object
func IsOverlapping(l1 Cord, r1 Cord, l2 Cord, r2 Cord) bool {
	if r2.Y >= r1.Y && r2.Y >= l1.Y && l2.Y >= r1.Y && l2.Y >= l1.Y {
		return false
	}

	if r2.Y <= r1.Y && r2.Y <= l1.Y && l2.Y <= r1.Y && l2.Y <= l1.Y {
		return false
	}

	if l2.X <= l1.X && l2.X <= r1.X && r2.X <= l1.X && r2.X <= r1.X {
		return false
	}

	if l2.X >= l1.X && l2.X >= r1.X && r2.X >= l1.X && r2.X >= r1.X {
		return false
	}

	return true;
}

// Cord is coordinates of one point of object
type Cord struct {
	X, Y int32
}

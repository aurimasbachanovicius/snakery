package geometrio

import "testing"

func TestIsOverlapping(t *testing.T) {
	type args struct {
		l1 Cord
		r1 Cord
		l2 Cord
		r2 Cord
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "two rect. touching the sides of each other",
			args: args{
				l1: Cord{X: 10, Y: 10},
				r1: Cord{X: 40, Y: 20},
				l2: Cord{X: 30, Y: 20},
				r2: Cord{X: 80, Y: 30},
			},
			want: false,
		}, {
			name: "one rect. coordinates ar too far away from another rect.",
			args: args{
				l1: Cord{X: 10, Y: 10},
				r1: Cord{X: 40, Y: 20},
				l2: Cord{X: 30, Y: 30},
				r2: Cord{X: 80, Y: 40},
			},
			want: false,
		}, {
			name: "one rect. coordinates overlaps another rect.",
			args: args{
				l1: Cord{X: 10, Y: 10},
				r1: Cord{X: 40, Y: 20},
				l2: Cord{X: 20, Y: 19},
				r2: Cord{X: 80, Y: 30},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsOverlapping(tt.args.l1, tt.args.r1, tt.args.l2, tt.args.r2); got != tt.want {
				t.Errorf("IsOverlapping() = %v, want %v", got, tt.want)
			}
		})
	}
}

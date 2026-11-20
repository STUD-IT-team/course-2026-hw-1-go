package summary

import "testing"

func TestSummarize(t *testing.T) {
	cases := []struct {
		name       string
		amounts    []int
		total      int
		max        int
		largeCount int
	}{
		{"no orders", nil, 0, 0, 0},
		{"zero amount", []int{0}, 0, 0, 0},
		{"single order", []int{1500}, 1500, 1500, 1},
		{"threshold", []int{999, 1000}, 1999, 1000, 1},
		{"maximum first", []int{2500, 100, 900}, 3500, 2500, 1},
		{"maximum in the middle", []int{700, 1300, 1000, 0}, 3000, 1300, 2},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			total, max, largeCount := Summarize(tc.amounts)
			if total != tc.total || max != tc.max || largeCount != tc.largeCount {
				t.Errorf("Summarize(%v) = (%d, %d, %d), want (%d, %d, %d)",
					tc.amounts, total, max, largeCount, tc.total, tc.max, tc.largeCount)
			}
		})
	}
}

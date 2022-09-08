package quadtree

import (
	"math"
	"reflect"
	"testing"
)

func TestContainsPoint(t *testing.T) {
	t.Parallel()

	type testCase struct {
		Rect rect
		X, Y float64
		Want bool
	}

	cases := []testCase{
		{rect{0, 0, 10, 10}, 2, 2, true},
		{rect{5, 5, 10, 10}, 1, 1, false},
		{rect{5, 5, 10, 10}, 15, 15, false},
		{rect{5, 5, 10, 10}, 6, 15, false},
		{rect{5, 5, 10, 10}, 15, 6, false},
	}

	for i, tc := range cases {
		if res := tc.Rect.ContainsPoint(tc.X, tc.Y); res != tc.Want {
			t.Fatalf("case[%d] failed, want %t got: %t", i, tc.Want, res)
		}
	}
}

func TestContainsRect(t *testing.T) {
	t.Parallel()

	type testCase struct {
		A, B rect
		Want bool
	}

	cases := []testCase{
		{rect{0, 0, 10, 10}, rect{2, 2, 8, 8}, true},
		{rect{5, 5, 10, 10}, rect{1, 1, 6, 6}, false},
		{rect{5, 5, 10, 10}, rect{6, 6, 12, 12}, false},
	}

	for i, tc := range cases {
		if res := tc.A.ContainsRect(tc.B); res != tc.Want {
			t.Fatalf("case[%d] failed, want %t got: %t", i, tc.Want, res)
		}
	}
}

func TestOverlaps(t *testing.T) {
	t.Parallel()

	type testCase struct {
		A, B rect
		Want bool
	}

	cases := []testCase{
		{rect{0, 0, 10, 10}, rect{2, 2, 8, 8}, true},
		{rect{5, 5, 10, 10}, rect{1, 1, 6, 6}, true},
		{rect{5, 5, 10, 10}, rect{6, 6, 12, 12}, true},
		{rect{5, 5, 10, 10}, rect{15, 0, 4, 4}, false},
	}

	for i, tc := range cases {
		if res := tc.A.Overlaps(tc.B); res != tc.Want {
			t.Fatalf("case[%d] failed, want %t got: %t", i, tc.Want, res)
		}
	}
}

func TestSplit(t *testing.T) {
	t.Parallel()

	type testCase struct {
		Rect rect
		Want [4]rect
	}

	cases := []testCase{
		{
			rect{0, 0, 8, 8},
			[4]rect{
				{0, 0, 4, 4},
				{4, 0, 8, 4},
				{0, 4, 4, 8},
				{4, 4, 8, 8},
			},
		},
		{
			rect{0, 0, 20, 20},
			[4]rect{
				{0, 0, 10, 10},
				{10, 0, 20, 10},
				{0, 10, 10, 20},
				{10, 10, 20, 20},
			},
		},
	}

	for i, tc := range cases {
		if res := tc.Rect.Split(); !reflect.DeepEqual(res, tc.Want) {
			t.Fatalf("case[%d] failed, want %v got: %v", i, tc.Want, res)
		}
	}
}

func TestDim(t *testing.T) {
	t.Parallel()

	const (
		e = 0.0000001
		w = 8.0
		h = 9.0
	)

	rc := rc(0, 0, w, h)

	if math.Abs(rc.Width()-w) > e {
		t.Fatal("width")
	}

	if math.Abs(rc.Heigth()-h) > e {
		t.Fatal("heigth")
	}
}

func TestPad(t *testing.T) {
	t.Parallel()

	const (
		e = 20.0000001
		w = 8.0
		h = 9.0
		d = 10.0
	)

	rc := rc(0, 0, w, h).Pad(d)

	if math.Abs(rc.Width()-w) > e {
		t.Fatal("width")
	}

	if math.Abs(rc.Heigth()-h) > e {
		t.Fatal("heigth")
	}
}

func TestClip(t *testing.T) {
	t.Parallel()

	const (
		w = 20.0
		h = 15.0
	)

	var area = rc(0, 0, w, h)

	var cases = []rect{
		rc(-1, -1, 2, 4),
		rc(-1, 4, 4, 5),
		rc(4, -2, 10, 10),
		rc(19, 5, 6, 6),
		rc(6, 14, 5, 5),
		rc(25, 16, 5, 5),
	}

	for i, c := range cases {
		tc := c.Clip(area)

		if tc.X0 < area.X0 {
			t.Fatalf("case[%d] fail fot X0: %f", i, tc.X0)
		}

		if tc.X1 > area.X1 {
			t.Fatalf("case[%d] fail fot X1: %f", i, tc.X1)
		}

		if tc.Y0 < area.Y0 {
			t.Fatalf("case[%d] fail fot Y0", i)
		}

		if tc.Y1 > area.Y1 {
			t.Fatalf("case[%d] fail fot Y1", i)
		}
	}
}

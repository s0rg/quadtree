package quadtree

import (
	"math/rand"
	"testing"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	q := New[int](100, 100, 4)

	if q.Size() != 0 {
		t.Fatal("non-empty")
	}

	q.Add(10, 10, 5, 5, 1)
	q.Add(15, 20, 10, 10, 2)
	q.Add(40, 10, 4, 4, 3)
	q.Add(90, 90, 5, 8, 4)

	if q.Size() != 4 {
		t.Fatal("empty")
	}

	var (
		val int
		ok  bool
	)

	if val, ok = q.Get(9, 9, 11, 11); !ok {
		t.Fatal("got nothing at (9, 9, 11, 11)")
	}

	if val != 1 {
		t.Fatalf("unexpected value: %d", val)
	}

	if _, ok = q.Get(101, 101, 111, 111); ok {
		t.Fatal("got something at (101, 101, 111, 111)")
	}

	if val, ok = q.Get(39, 9, 42, 12); !ok {
		t.Fatal("got nothing at (39, 9, 42, 12)")
	}

	if val != 3 {
		t.Fatalf("invalid value: %d", val)
	}
}

func TestDel(t *testing.T) {
	t.Parallel()

	q := New[int](100, 100, 4)

	q.Add(5, 5, 2, 2, 1)
	q.Add(10, 10, 3, 3, 2)
	q.Add(11, 11, 4, 4, 3)

	q.Del(12, 12)

	if s := q.Size(); s != 2 {
		t.Fatalf("unexpected len = %d", s)
	}
}

func TestForEach(t *testing.T) {
	t.Parallel()

	q := New[int](100, 100, 4)

	q.Add(5, 5, 2, 3, 1)
	q.Add(10, 10, 5, 5, 2)
	q.Add(11, 11, 3, 3, 3)
	q.Add(25, 25, 10, 10, 100)

	var sum int

	q.ForEach(1, 1, 20, 20, func(_, _, _, _ float64, val int) {
		sum += val
	})

	const sumExpect = 6

	if sum != sumExpect {
		t.Fatalf("unexpected sum = %d", sum)
	}
}

func TestMove(t *testing.T) {
	t.Parallel()

	q := New[int](100, 100, 4)

	q.Add(5, 5, 10, 10, 1)

	var (
		val int
		ok  bool
	)

	if val, ok = q.Get(4, 4, 6, 6); !ok || val != 1 {
		t.Fatal("step 1 unexpected")
	}

	if q.Move(55, 55, 0, 0) {
		t.Fatal("unexpected move success")
	}

	if !q.Move(5, 5, 50, 50) {
		t.Fatal("cannot move")
	}

	if val, ok = q.Get(45, 45, 55, 55); !ok || val != 1 {
		t.Fatal("step 2 unexpected")
	}
}

func TestKNearest(t *testing.T) {
	t.Parallel()

	q := New[int](50, 50, 4)

	q.Add(1, 1, 2, 2, 1)
	q.Add(5, 1, 2, 2, 2)
	q.Add(1, 5, 2, 2, 3)
	q.Add(10, 10, 2, 2, 10)
	q.Add(1, 10, 2, 2, 10)

	var sum int

	q.KNearest(3, 3, 5, 3, func(_, _, _, _ float64, val int) {
		sum += val
	})

	const sumExpect = 6

	if sum != sumExpect {
		t.Fatalf("unexpected sum = %d", sum)
	}
}

func BenchmarkTreeKNearest10(b *testing.B)  { benchmarkKNearestN(b, 10) }
func BenchmarkTreeKNearest100(b *testing.B) { benchmarkKNearestN(b, 100) }

func benchmarkKNearestN(b *testing.B, count int) {
	b.Helper()

	const maxSize = 10.0

	rand.Seed(benchmarkSeed)

	q := New[int](benchmarkSide, benchmarkSide, nodeTestDepth)

	for i := 0; i < count; i++ {
		rc := randRect(benchmarkSide-maxSize, maxSize)
		q.Add(rc.X0, rc.Y0, rc.Width(), rc.Heigth(), i)
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		for i := 0; i < count; i++ {
			b.StopTimer()

			x, y := randPoint(maxSize)
			dist := randFloat(1, maxSize)
			k := int(randFloat(1, maxSize))

			b.StartTimer()

			q.KNearest(x, y, dist, k, func(_, _, _, _ float64, val int) {})
		}
	}
}

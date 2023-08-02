package quadtree

import (
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

func BenchmarkTree(b *testing.B) {
	const (
		benchmarkSide = 1000.0
		maxSize       = 10.0
	)

	q := New[int](benchmarkSide, benchmarkSide, nodeTestDepth)

	x, y := 1.0, 1.0
	x2, y2 := 800.0, 800.0

	b.ResetTimer()

	b.Run("Add", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			q.Add(x, y, maxSize, maxSize, n)
		}
	})

	b.Run("Get", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, _ = q.Get(x, y, maxSize, maxSize)
		}
	})

	b.Run("Move", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = q.Move(x, y, x2, y2)
			x, y, x2, y2 = x2, y2, x, y
		}
	})

	b.Run("ForEach", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			q.ForEach(x, y, maxSize, maxSize, func(_, _, _, _ float64, val int) {})
		}
	})

	b.Run("KNearest", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			q.KNearest(x, y, maxSize, int(maxSize), func(_, _, _, _ float64, val int) {})
		}
	})

	b.Run("Del", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = q.Del(x, y)
		}
	})
}

package quadtree

import (
	"math/rand"
	"testing"
	"time"
)

const nodeTestDepth = 4

func TestNodeInsertSearchOK(t *testing.T) {
	t.Parallel()

	n := makeNode[int](rect{0, 0, 80, 80})
	n.Grow(nodeTestDepth, 0)

	if s := n.Size(); s != 0 {
		t.Fatalf("empty but size is: %d", s)
	}

	n.Insert(rect{1, 1, 5, 5}, 1)
	n.Insert(rect{20, 23, 25, 28}, 2)
	n.Insert(rect{30, 31, 35, 55}, 3)
	n.Insert(rect{36, 32, 39, 58}, 4)

	if s := n.Size(); s != 4 {
		t.Fatalf("expected size 3 but size is: %d", s)
	}

	var found int

	n.Search(rect{29, 30, 40, 60}, func(it *item[int]) bool {
		switch it.Value {
		case 3, 4:
			found++
		default:
			t.Fatalf("wrong value found: %d at (%f, %f)", it.Value, it.Rect.X0, it.Rect.Y0)
		}

		return (found < 2)
	})

	if found != 2 {
		t.Fatalf("wrong found count: %d", found)
	}

	n.Search(rect{18, 18, 30, 30}, func(it *item[int]) bool {
		if it.Value != 2 {
			t.Fatalf("wrong value found: %d at (%f, %f)", it.Value, it.Rect.X0, it.Rect.Y0)
		}

		return false
	})
}

func TestNodeInsertSearchOOB(t *testing.T) {
	t.Parallel()

	n := makeNode[int](rect{10, 10, 80, 80})
	n.Grow(nodeTestDepth, 0)

	if n.Insert(rect{1, 1, 5, 5}, 1) {
		t.Fatal("out-of-bound insert 1")
	}

	if n.Insert(rect{120, 23, 125, 28}, 2) {
		t.Fatal("out-of-bound insert 2")
	}

	if n.Size() != 0 {
		t.Fatal("out-of-bound insert size")
	}

	n.Search(rect{90, 10, 100, 5}, func(it *item[int]) bool {
		t.Fatalf("found out-of-bound: %d at (%f, %f)", it.Value, it.Rect.X0, it.Rect.Y0)

		return false
	})
}

func TestNodeDel(t *testing.T) {
	t.Parallel()

	n := makeNode[int](rect{0, 0, 80, 80})
	n.Grow(nodeTestDepth/2, 0)

	n.Insert(rect{1, 1, 5, 5}, 1)
	n.Insert(rect{20, 23, 25, 28}, 2)
	n.Insert(rect{30, 31, 35, 55}, 3)
	n.Insert(rect{26, 26, 28, 28}, 4)

	n.Del(100, 100)

	if n.Size() != 4 {
		t.Fatal("del remove out-of-bound")
	}

	n.Del(70, 70)

	if n.Size() != 4 {
		t.Fatal("del remove unexisted")
	}

	n.Del(21, 24)

	if n.Size() != 3 {
		t.Fatal("del doesnt remove 2")
	}

	n.Search(rect{18, 18, 30, 30}, func(it *item[int]) bool {
		if it.Value == 2 {
			t.Fatalf("2 still exists at (%f, %f)", it.Rect.X0, it.Rect.Y0)
		}

		return false
	})

	n.Del(32, 42)

	if n.Size() != 2 {
		t.Fatal("del doesnt remove 3")
	}

	n.Search(rect{31, 32, 31, 32}, func(it *item[int]) bool {
		if it.Value == 3 {
			t.Fatalf("3 still exists %d at (%f, %f)", it.Value, it.Rect.X0, it.Rect.Y0)
		}

		return false
	})
}

const benchmarkSide = 1000.0

var benchmarkSeed = time.Now().UnixNano()

func BenchmarkNodeInsert1(b *testing.B)   { benchmarkInsertN(b, 1) }
func BenchmarkNodeInsert10(b *testing.B)  { benchmarkInsertN(b, 10) }
func BenchmarkNodeInsert100(b *testing.B) { benchmarkInsertN(b, 100) }

func BenchmarkNodeDel10(b *testing.B)  { benchmarkDelN(b, 10) }
func BenchmarkNodeDel100(b *testing.B) { benchmarkDelN(b, 100) }

func BenchmarkNodeSearch10(b *testing.B)  { benchmarkSearchN(b, 10) }
func BenchmarkNodeSearch100(b *testing.B) { benchmarkSearchN(b, 100) }

func benchmarkInsertN(b *testing.B, count int) {
	b.Helper()

	const maxSize = 10.0

	rand.Seed(benchmarkSeed)

	node := makeNode[int](rc(0, 0, benchmarkSide, benchmarkSide))
	node.Grow(nodeTestDepth, 0)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		for i := 0; i < count; i++ {
			b.StopTimer()

			rc := randRect(benchmarkSide-maxSize, maxSize)

			b.StartTimer()

			node.Insert(rc, i)
		}
	}
}

func benchmarkDelN(b *testing.B, count int) {
	b.Helper()

	const maxSize = 10.0

	rand.Seed(benchmarkSeed)

	node := makeNode[int](rc(0, 0, benchmarkSide, benchmarkSide))
	node.Grow(nodeTestDepth, 0)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		b.StopTimer()

		for i := 0; i < count; i++ {
			rc := randRect(benchmarkSide-maxSize, maxSize)
			node.Insert(rc, i)
		}

		for i := 0; i < count; i++ {
			x, y := randPoint(benchmarkSide - maxSize)

			b.StartTimer()
			node.Del(x, y)
			b.StopTimer()
		}
	}
}

func benchmarkSearchN(b *testing.B, count int) {
	b.Helper()

	const maxSize = 10.0

	rand.Seed(benchmarkSeed)

	node := makeNode[int](rc(0, 0, benchmarkSide, benchmarkSide))
	node.Grow(nodeTestDepth, 0)

	for i := 0; i < count; i++ {
		rc := randRect(benchmarkSide-maxSize, maxSize)
		node.Insert(rc, i)
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		for i := 0; i < count; i++ {
			b.StopTimer()

			rc := randRect(benchmarkSide-maxSize, maxSize)

			b.StartTimer()

			node.Search(rc, func(_ *item[int]) bool { return false })
		}
	}
}

func randFloat(min, max float64) (rv float64) {
	return min + (rand.Float64() * (max - min))
}

func randRect(maxPos, maxSide float64) (r rect) {
	x, y := randFloat(1, maxPos-1), randFloat(1, maxPos-1)
	w, h := randFloat(1, maxSide), randFloat(1, maxSide)

	return rc(x, y, w, h)
}

func randPoint(maxPos float64) (x, y float64) {
	x, y = randFloat(1, maxPos-1), randFloat(1, maxPos-1)

	return
}

package quadtree

import "math"

// Tree is a quad-tree container.
type Tree[T any] struct {
	root *node[T]
}

// Iter is a shorthand for iterator callback type.
type Iter[T any] func(x, y, w, h float64, value T)

// New creates new, empty quad tree instance with given area width and height and maximum depth value.
func New[T any](w, h float64, depth int) (rv *Tree[T]) {
	root := makeNode[T](rc(0, 0, w, h))
	root.Grow(depth, 0)

	return &Tree[T]{
		root: root,
	}
}

// Size returns current elements count.
func (t *Tree[T]) Size() int {
	return t.root.Size()
}

// Add adds rectangle, with top-left corner at (x, y) and given width and height.
func (t *Tree[T]) Add(x, y, w, h float64, value T) (ok bool) {
	return t.root.Insert(rc(x, y, w, h), value)
}

// Get returns item located at given coordinates, if any.
func (t *Tree[T]) Get(x, y, w, h float64) (value T, ok bool) {
	area := rc(x, y, w, h).Clip(t.root.Rect)

	t.root.Search(area, func(it *item[T]) (next bool) {
		value, ok = it.Value, true

		return false
	})

	return value, ok
}

// Del removes any items located at given coordinates, returns true if succeed.
func (t *Tree[T]) Del(x, y float64) (ok bool) {
	_, ok = t.root.Del(x, y)

	return
}

// Move moves item located at (x, y) to (newX, newY), returns true if succeed.
func (t *Tree[T]) Move(x, y, newX, newY float64) (ok bool) {
	var it item[T]

	if it, ok = t.root.Del(x, y); !ok {
		return
	}

	return t.Add(newX, newY, it.Rect.Width(), it.Rect.Heigth(), it.Value)
}

// ForEach iterates over items in given region.
func (t *Tree[T]) ForEach(x, y, w, h float64, iter Iter[T]) {
	area := rc(x, y, w, h).Clip(t.root.Rect)

	t.root.Search(area, func(it *item[T]) (next bool) {
		iter(it.Rect.X0, it.Rect.Y0, it.Rect.Width(), it.Rect.Heigth(), it.Value)

		return true
	})
}

// KNearest iterates over k nearest for given coordinates within given distance.
func (t *Tree[T]) KNearest(x, y, distance float64, k int, iter Iter[T]) {
	var (
		found int
		area  = rc(x, y, x, y).Pad(distance).Clip(t.root.Rect)
	)

	t.root.Search(area, func(it *item[T]) (next bool) {
		cx, cy := it.Rect.Center()

		if dist2d(x, y, cx, cy) <= distance {
			iter(it.Rect.X0, it.Rect.Y0, it.Rect.Width(), it.Rect.Heigth(), it.Value)

			found++
		}

		return found < k
	})
}

func dist2d(x1, y1, x2, y2 float64) (rv float64) {
	dx, dy := x2-x1, y2-y1

	return math.Sqrt(dx*dx + dy*dy)
}

package quadtree

const childCount = 4

type item[T any] struct {
	Value T
	Rect  rect
}

type node[T any] struct {
	Childs []*node[T]
	Items  []item[T]
	Rect   rect
}

func makeNode[T any](r rect) *node[T] {
	return &node[T]{
		Rect:   r,
		Childs: []*node[T]{},
	}
}

func (n *node[T]) Grow(max, cur int) {
	if cur > max {
		return
	}

	n.Childs = make([]*node[T], childCount)

	for i, r := range n.Rect.Split() {
		n.Childs[i] = makeNode[T](r)
		n.Childs[i].Grow(max, cur+1)
	}
}

func (n *node[T]) Insert(r rect, value T) (ok bool) {
	if !n.Rect.ContainsRect(r) {
		return
	}

	for i := 0; i < len(n.Childs); i++ {
		if c := n.Childs[i]; c.Rect.ContainsRect(r) {
			return c.Insert(r, value)
		}
	}

	n.Items = append(n.Items, item[T]{Rect: r, Value: value})

	return true
}

func (n *node[T]) Search(r rect, iter func(*item[T]) bool) {
	if !n.Rect.Overlaps(r) {
		return
	}

	var it *item[T]

	for i := 0; i < len(n.Items); i++ {
		if it = &n.Items[i]; !it.Rect.Overlaps(r) {
			continue
		}

		if !iter(it) {
			return
		}
	}

	for i := 0; i < len(n.Childs); i++ {
		switch {
		case r.ContainsRect(n.Childs[i].Rect):
			if !n.Childs[i].ForEach(iter) {
				return
			}
		case n.Childs[i].Rect.Overlaps(r):
			n.Childs[i].Search(r, iter)
		}
	}
}

func (n *node[T]) ForEach(iter func(*item[T]) bool) (next bool) {
	for i := 0; i < len(n.Items); i++ {
		if !iter(&n.Items[i]) {
			return false
		}
	}

	for i := 0; i < len(n.Childs); i++ {
		if !n.Childs[i].ForEach(iter) {
			return false
		}
	}

	return true
}

func (n *node[T]) Size() (total int) {
	total = len(n.Items)

	for i := 0; i < len(n.Childs); i++ {
		total += n.Childs[i].Size()
	}

	return total
}

func (n *node[T]) Del(x, y float64) (it item[T], ok bool) {
	if !n.Rect.ContainsPoint(x, y) {
		return
	}

	for i := 0; i < len(n.Items); i++ {
		if it = n.Items[i]; it.Rect.ContainsPoint(x, y) {
			last := len(n.Items) - 1

			if i != last {
				n.Items[i] = n.Items[last]
			}

			n.Items = n.Items[:last]

			return it, true
		}
	}

	for i := 0; i < len(n.Childs); i++ {
		if it, ok = n.Childs[i].Del(x, y); ok {
			return it, ok
		}
	}

	return
}

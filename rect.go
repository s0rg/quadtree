package quadtree

const half = 2.0

type rect struct {
	X0, Y0, X1, Y1 float64
}

func rc(x, y, w, h float64) (r rect) {
	return rect{
		X0: x,
		Y0: y,
		X1: x + w,
		Y1: y + h,
	}
}

func (r rect) Pad(d float64) (rv rect) {
	return rect{
		X0: r.X0 - d,
		Y0: r.Y0 - d,
		X1: r.X1 + d,
		Y1: r.Y1 + d,
	}
}

func (r rect) Clip(b rect) (rv rect) {
	if r.X0 < b.X0 {
		r.X0 = b.X0
	}

	if r.Y0 < b.Y0 {
		r.Y0 = b.Y0
	}

	if r.X1 > b.X1 {
		r.X1 = b.X1
	}

	if r.Y1 > b.Y1 {
		r.Y1 = b.Y1
	}

	return r
}

func (r *rect) Center() (x, y float64) {
	dx, dy := r.Width()/half, r.Heigth()/half

	return r.X0 + dx, r.Y0 + dy
}

func (r *rect) Width() float64 {
	return r.X1 - r.X0
}

func (r *rect) Heigth() float64 {
	return r.Y1 - r.Y0
}

func (r *rect) ContainsPoint(x, y float64) (yes bool) {
	if x < r.X0 || y < r.Y0 {
		return
	}

	if x > r.X1 || y > r.Y1 {
		return
	}

	return true
}

func (r *rect) ContainsRect(v rect) (yes bool) {
	if !r.ContainsPoint(v.X0, v.Y0) {
		return
	}

	if !r.ContainsPoint(v.X1, v.Y1) {
		return
	}

	return true
}

func (r *rect) Overlaps(v rect) (yes bool) {
	return r.X0 < v.X1 && r.X1 >= v.X0 && r.Y0 < v.Y1 && r.Y1 >= v.Y0
}

func (r *rect) Split() (rv [4]rect) {
	dx, dy := r.Width()/half, r.Heigth()/half

	return [4]rect{
		{r.X0, r.Y0, r.X0 + dx, r.Y0 + dy}, // upper left
		{r.X0 + dx, r.Y0, r.X1, r.Y1 - dy}, // upper right
		{r.X0, r.Y0 + dy, r.X0 + dx, r.Y1}, // lower left
		{r.X0 + dx, r.Y0 + dy, r.X1, r.Y1}, // lower right
	}
}

package point

type Point struct {
	X, Y int
}

type Vec struct {
	Dx, Dy int
}

func New(x, y int) Point {
	return Point{x, y}
}

func (p Point) Sub(q Point) Vec {
	return Vec{p.X - q.X, p.Y - q.Y}
}

func (p Point) Move(v Vec) Point {
	return Point{p.X + v.Dx, p.Y + v.Dy}
}

func (v Vec) Neg() Vec {
	return Vec{-v.Dx, -v.Dy}
}

func (v Vec) Scale(s int) Vec {
	return Vec{s * v.Dx, s * v.Dy}
}

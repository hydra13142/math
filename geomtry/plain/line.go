package plain

import "math"

// 表示一个直线
type Line Segment

// 两点确定一条直线
func (s *Line) Set(A, B Spot) {
	s.O, s.K = A, AimTo(A, B)
}

// 点p相对直线s的位置
func (s *Line) Location(p Spot) (i int) {
	switch k := OuterProduct(AimTo(s.O, p), s.K); {
	case k < 0:
		return WL
	case k > 0:
		return WR
	default:
		return WM
	}
}

// 点p到直线s所在直线的距离
func (s *Line) Distance(p Spot) float64 {
	ap := AimTo(s.O, p)
	return math.Abs(OuterProduct(ap, s.K)) / s.K.Abs()
}

// 点p到直线s所在直线的垂点
func (s *Line) Vertical(p Spot) Spot {
	ab := s.K.Unit()
	ds := OuterProduct(AimTo(s.O, p), ab)
	return p.Move(ab.Spin(math.Pi / 2).Mul(ds))
}

// 两个直线的交集，交点个数、具体的点。
func (s *Line) Cross(l *Line) (int, Spot) {
	sl := AimTo(s.O, l.O)
	th := OuterProduct(s.K, l.K)
	if th != 0 {
		k := OuterProduct(sl, l.K) / th
		return Intersect, Spot{s.O.X + s.K.I*k, s.O.Y + s.K.J*k}
	}
	if OuterProduct(s.K, sl) != 0 {
		return Paraellel, Spot{}
	}
	return Collinear, Spot{}
}

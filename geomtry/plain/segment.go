package plain

import "math"

const (
	WM = 1 << iota // 横向：中间
	WL             // 横向：左侧
	WR             // 横向：右侧
	HM             // 纵向：中间
	HT             // 纵向：上方
	HB             // 纵向：下方
)

const (
	Intersect = 0 // 相交
	Paraellel = 1 // 平行
	Collinear = 3 // 共线
)

// 表示一个线段
type Segment struct {
	O Spot
	K Vector
}

// 两点确定一条线段
func (s *Segment) Set(A, B Spot) {
	s.O, s.K = A, AimTo(A, B)
}

// 返回线段的长度。
func (s *Segment) Length() float64 {
	return s.K.Abs()
}

// 点p相对线段s的位置
func (s *Segment) Location(p Spot) (i int) {
	ap := AimTo(s.O, p)
	bp := ap.Dec(s.K)
	switch k := OuterProduct(ap, s.K); {
	case k < 0:
		i = WL
	case k > 0:
		i = WR
	default:
		i = WM
	}
	if InnerProduct(s.K, bp) > 0 {
		return i | HT
	}
	if InnerProduct(s.K, ap) < 0 {
		return i | HB
	}
	return i | HM
}

// 两个线段的交集，返回两个线段所在直线的位置关系、交点个数、具体的点。
func (s *Segment) Cross(l *Segment) (int, []Spot) {
	var A, B Spot
	sl := AimTo(s.O, l.O)
	th := OuterProduct(s.K, l.K)
	if th != 0 {
		k := OuterProduct(sl, l.K) / th
		A = Spot{s.O.X + s.K.I*k, s.O.Y + s.K.J*k}
		p := AimTo(A, l.O)
		if InnerProduct(p, p.Add(l.K)) > 0 {
			return Intersect, nil
		}
		p = AimTo(A, s.O)
		if InnerProduct(p, p.Add(s.K)) > 0 {
			return Intersect, nil
		}
		return Intersect, []Spot{A}
	}
	if OuterProduct(s.K, sl) != 0 {
		return Paraellel, nil
	}
	if InnerProduct(s.K, l.K) < 0 {
		A, B = l.O.Move(l.K), l.O
	} else {
		A, B = l.O, l.O.Move(l.K)
	}
	switch both(s, &A, &B) {
	case 0:
		return Collinear, nil
	case 1:
		return Collinear, []Spot{A}
	default:
		return Collinear, []Spot{A, B}
	}
}

// 已知A、B在s所在直线上，A=>B与s方向相同
// 返回线段AB与线段s的公共部分（0：无，1：点，2：线段）
func both(s *Segment, A, B *Spot) int {
	var X, Y float64
	if math.Abs(s.K.I) > math.Abs(s.K.J) {
		X = (A.X - s.O.X) / s.K.I
		Y = (B.X - s.O.X) / s.K.I
	} else {
		X = (A.Y - s.O.Y) / s.K.J
		Y = (B.Y - s.O.Y) / s.K.J
	}
	if X > 1 || Y < 0 {
		return 0
	}
	if X == 1 {
		return 1
	}
	if Y == 0 {
		*A = s.O
		return 1
	}
	if X <= 0 {
		*A = s.O
	}
	if Y >= 1 {
		*B = s.O.Move(s.K)
	}
	return 2
}

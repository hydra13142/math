package solid

const (
	Uncoplane = 1  // 异面
	Intersect = 2  // 相交
	Paraellel = 4  // 平行
	Collinear = 12 // 共线
)

// 表示一条直线，隐藏内部的变量，以便在创建时检查合法性
type Line struct {
	O Spot
	K Vector
}

// 两点确定一条直线
func (l *Line) Set(p, q Spot) {
	l.O = p
	l.K = AimTo(p, q)
}

// 返回点到直线的距离
func (l *Line) Distance(p Spot) float64 {
	return OuterProduct(AimTo(p, l.O), l.K).Abs() / l.K.Abs()
}

// 求点到直线的垂点
func (l *Line) Vertical(p Spot) Spot {
	return l.O.Move(l.K.Mul(InnerProduct(AimTo(l.O, p), l.K) / InnerProduct(l.K, l.K)))
}

// 两条直线的位置关系
func (l *Line) Cross(s *Line) (int, []Spot) {
	ref := AimTo(l.O, s.O)
	if OuterProduct(l.K, s.K).Abs() == 0 {
		if OuterProduct(l.K, ref).Abs() == 0 {
			return Collinear, nil
		}
		return Paraellel, nil
	}
	a := InnerProduct(l.K, s.K)
	b := InnerProduct(l.K, l.K)
	c := InnerProduct(s.K, s.K)
	d := InnerProduct(ref, l.K)
	e := InnerProduct(ref, s.K)
	k := a*a - b*c
	x := (a*e - c*d) / k
	y := (b*e - a*d) / k
	p0 := l.O.Move(l.K.Mul(x))
	p1 := s.O.Move(s.K.Mul(y))
	if p0 == p1 {
		return Intersect, []Spot{p0}
	}
	return Uncoplane, []Spot{p0, p1}
}

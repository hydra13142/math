package solid

// 一个平面
type Plain struct {
	O Spot   // 基准点
	K Vector // 法向量
}

// 三点确定一个平面
func (p *Plain) Set(a, b, c Spot) {
	p.O = Spot{(a.X + b.X + c.X) / 3, (a.Y + b.Y + c.Y) / 3, (a.Z + b.Z + c.Z) / 3}
	p.K = OuterProduct(AimTo(a, c), AimTo(a, b))
}

// 求点到面的距离
func (p *Plain) Distance(s Spot) float64 {
	return InnerProduct(AimTo(p.O, s), p.K) / p.K.Abs()
}

// 求点到面的垂点
func (p *Plain) Vertical(s Spot) Spot {
	return s.Move(p.K.Mul(InnerProduct(AimTo(s, p.O), p.K) / InnerProduct(p.K, p.K)))
}

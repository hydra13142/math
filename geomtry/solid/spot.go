package solid

// 表示一个点。
type Spot struct {
	X, Y, Z float64
}

// 移动一个点
func (p Spot) Move(v Vector) Spot {
	return Spot{p.X + v.I, p.Y + v.J, p.Z + v.K}
}

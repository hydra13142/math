package plain

// 表示一个点。
type Spot struct {
	X, Y float64
}

// 返回一个点移动后的新位置，不移动自身
func (p Spot) Move(v Vector) Spot {
	return Spot{p.X + v.I, p.Y + v.J}
}

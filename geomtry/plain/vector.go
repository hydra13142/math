package plain

import "math"

// 表示一个向量。
type Vector struct {
	I, J float64
}

// 返回向量a->b。
func AimTo(a, b Spot) Vector {
	return Vector{b.X - a.X, b.Y - a.Y}
}

// 返回向量的模。
func (v Vector) Abs() float64 {
	return math.Hypot(v.I, v.J)
}

// 返回同向的单位向量
func (v Vector) Unit() Vector {
	r := v.Abs()
	return Vector{v.I / r, v.J / r}
}

// 返回旋转该向量一定弧度的向量
func (v Vector) Spin(r float64) Vector {
	c, s := math.Cos(r), math.Sin(r)
	return Vector{v.I*c - v.J*s, v.J*c + v.I*s}
}

// 返回两个向量的和
func (v Vector) Add(u Vector) Vector {
	return Vector{v.I + u.I, v.J + u.J}
}

// 返回两个向量的差
func (v Vector) Dec(u Vector) Vector {
	return Vector{v.I - u.I, v.J - u.J}
}

// 返回系数和向量的乘积
func (v Vector) Mul(k float64) Vector {
	return Vector{v.I * k, v.J * k}
}

// 返回两个向量的内积。
func InnerProduct(p, q Vector) float64 {
	return p.I*q.I + p.J*q.J
}

// 返回两个向量的外积；因为是平面，返回浮点数。
func OuterProduct(p, q Vector) float64 {
	return p.I*q.J - q.I*p.J
}

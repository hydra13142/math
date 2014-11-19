package solid

import "math"

type Vector struct {
	I, J, K float64
}

// 返回向量a->b
func AimTo(a, b Spot) Vector {
	return Vector{b.X - a.X, b.Y - a.Y, b.Z - a.Z}
}

// 返回向量的模。
func (v Vector) Abs() float64 {
	return math.Sqrt(v.I*v.I + v.J*v.J + v.K*v.K)
}

// 返回同向的单位向量；如果向量的模为零，返回该向量。
func (v Vector) Unit() Vector {
	r := v.Abs()
	return Vector{v.I / r, v.J / r, v.K / r}
}

// 返回两个向量的和
func (v Vector) Add(u Vector) Vector {
	return Vector{v.I + u.I, v.J + u.J, v.K + u.K}
}

// 返回两个向量的差
func (v Vector) Dec(u Vector) Vector {
	return Vector{v.I - u.I, v.J - u.J, v.K - u.K}
}

// 返回系数和向量的乘积
func (v Vector) Mul(k float64) Vector {
	return Vector{v.I * k, v.J * k, v.K * k}
}

// 返回两个向量的内积
func InnerProduct(p, q Vector) float64 {
	return p.I*q.I + p.J*q.J + p.K*q.K
}

// 返回两个向量的外积
func OuterProduct(p, q Vector) Vector {
	return Vector{p.J*q.K - q.J*p.K, q.I*p.K - p.I*q.K, p.I*q.J - q.I*p.J}
}

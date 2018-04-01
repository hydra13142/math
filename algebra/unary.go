package algebra

import (
	"math"
	"sort"
)

// 一元非负整次数多项式
type Unary []float64

// 返回多项式的次数
func (this Unary) Order() int {
	return len(this) - 1
}

// 计算多项式函数的值, a[0]+a[1]*x+a[2]*x^2+...+a[n]*x^n
func (this Unary) Compute(x float64) (y float64) {
	for i := len(this) - 1; i >= 0; i-- {
		y = y*x + this[i]
	}
	return y
}

// 求微分多项式
func (this Unary) Reduce() Unary {
	if len(this) <= 1 {
		return Unary{0}
	}
	that := make(Unary, len(this)-1)
	for i, e := range this[1:] {
		that[i] = float64(i+1) * e
	}
	return that
}

// 求积分多项式
func (this Unary) Integral() Unary {
	if len(this) == 0 {
		return Unary{0}
	}
	that := make(Unary, len(this)+1)
	for i, e := range this {
		that[i+1] = e / float64(i+1)
	}
	return that
}

// 多项式相加
func (this Unary) Add(that Unary) Unary {
	x, y := len(this), len(that)
	if x < y {
		this, x, that, y = that, y, this, x
	}
	i, r := 0, make(Unary, x)
	for ; i < y; i++ {
		r[i] = this[i] + that[i]
	}
	for ; i < x; i++ {
		r[i] = this[i]
	}
	return r
}

// 多项式相减
func (this Unary) Sub(that Unary) Unary {
	x, y := len(this), len(that)
	if x >= y {
		i, r := 0, make(Unary, x)
		for ; i < y; i++ {
			r[i] = this[i] - that[i]
		}
		for ; i < x; i++ {
			r[i] = this[i]
		}
		return r
	} else {
		i, r := 0, make(Unary, y)
		for ; i < x; i++ {
			r[i] = this[i] - that[i]
		}
		for ; i < y; i++ {
			r[i] = -that[i]
		}
		return r
	}
}

// 多项式相乘
func (this Unary) Mul(that Unary) Unary {
	x, y := len(this), len(that)
	r := make(Unary, x+y-1)
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			r[i+j] += this[i] * that[j]
		}
	}
	return r
}

// 多项式相除
func (this Unary) Div(that Unary) Unary {
	x, y := len(this), len(that)
	if y > x {
		return Unary{0}
	}
	u := make(Unary, x)
	r := make(Unary, x-y+1)
	copy(u, this)
	for i, j := x-y, x-1; i >= 0; i, j = i-1, j-1 {
		k := u[j] / that[y-1]
		for a, b := j, y-1; b >= 0; a, b = a-1, b-1 {
			u[a] -= k * that[b]
		}
		r[i] = k
	}
	return r
}

// 多项式取余
func (this Unary) Mod(that Unary) Unary {
	x, y := len(this), len(that)
	if y > x {
		return this
	}
	u := make(Unary, x)
	copy(u, this)
	for i := x - 1; i >= y-1; i-- {
		k := u[i] / that[y-1]
		for a, b := i, y-1; b >= 0; a, b = a-1, b-1 {
			u[a] -= k * that[b]
		}
	}
	for i := y - 2; i >= 0; i-- {
		if u[i] != 0 {
			return u[:i+1]
		}
	}
	return Unary{0}
}

// 多项式相除并取余
func (this Unary) DivMod(that Unary) (Unary, Unary) {
	x, y := len(this), len(that)
	if y > x {
		return Unary{0}, this
	}
	u := make(Unary, x)
	r := make(Unary, x-y+1)
	copy(u, this)
	for i, j := x-y, x-1; i >= 0; i, j = i-1, j-1 {
		k := u[j] / that[y-1]
		for a, b := j, y-1; b >= 0; a, b = a-1, b-1 {
			u[a] -= k * that[b]
		}
		r[i] = k
	}
	for i := y - 2; i >= 0; i-- {
		if u[i] != 0 {
			return r, u[:i+1]
		}
	}
	return r, Unary{0}
}

// 乘以一个系数
func (this Unary) ScalarMul(k float64) Unary {
	x := len(this)
	r := make(Unary, x)
	for i := 0; i < x; i++ {
		r[i] = this[i] * k
	}
	return r
}

// n个p相乘
func (this Unary) Pow(n uint) Unary {
	if n == 0 {
		return Unary{1}
	}
	r := make(Unary, len(this)*int(n))
	copy(r, this)
	for ; n > 1; n-- {
		r = r.Mul(this)
	}
	return r
}

// 移动多项式曲线：x>0向右移动，x<0向左移动；y>0向上移动，y<0向下移动。
func (this Unary) Move(x, y float64) Unary {
	r := make(Unary, len(this))
	if x != 0 {
		T := make(Unary, len(this))
		k := make(Unary, len(this))
		for i := 0; i < len(this); i++ {
			T[i] = 1
			k[i] = 1
			for j := i; j > 0; j-- {
				r[j] += k[j] * T[j] * this[i]
				T[j] *= -x
				k[j] += k[j-1]
			}
			r[0] += k[0] * T[0] * this[i]
			T[0] *= -x
		}
	} else {
		copy(r, this)
	}
	r[0] += y
	return r
}

// 求该多项式值为零时的x的解集。次数越高，求解越慢。
func SolveUnary(p Unary) []float64 {
	var s []float64
	i, j := 0, len(p)-1
	for i <= j && p[i] == 0 {
		i++
	}
	for j >= i && p[j] == 0 {
		j--
	}
	j++
	p = p[i:j]
	switch len(p) {
	case 0:
		return nil
	case 1:
		break
	case 2:
		if p[1] != 0 {
			s = []float64{-p[0] / p[1]}
		}
	case 3:
		a, b, c := p[2], p[1], p[0]
		switch k := b*b - 4*a*c; {
		case k > 0:
			k = math.Sqrt(k)
			if a > 0 {
				s = []float64{(-b - k) / (2 * a), (-b + k) / (2 * a)}
			} else {
				s = []float64{(-b + k) / (2 * a), (-b - k) / (2 * a)}
			}
		case k == 0:
			s = []float64{-b / (2 * a)}
		}
	default:
		q := p.Reduce()
		n := SolveUnary(q)
		if l := len(n); l == 0 {
			if x, e := Monotone(p.Compute, q.Compute); e == nil {
				s = append(s, x)
			}
		} else {
			if x, e := Tangent(p.Compute, q.Compute, math.Inf(-1), n[0], n[0]-1); e == nil {
				s = append(s, x)
			}
			for i, l := 0, l-1; i < l; i++ {
				if x, e := Region(p.Compute, n[i], n[i+1]); e == nil {
					s = append(s, x)
				}
			}
			if x, e := Tangent(p.Compute, q.Compute, math.Inf(+1), n[l], n[l]+1); e == nil {
				s = append(s, x)
			}
		}
	}
	if i > 0 {
		s = append(s, 0)
	}
	sort.Float64s(s)
	j, t := 1, s[0]
	for i := 1; i < len(s); i++ {
		if k := s[i]; k != t {
			s[j], t, j = k, k, j+1
		}
	}
	return s[:j]
}

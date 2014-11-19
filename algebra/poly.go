package algebra

import "math"

// 一元非负整次数多项式
type Poly []float64

// 计算多项式函数的值
func (p Poly) Compute(x float64) (y float64) {
	for i := len(p) - 1; i >= 0; i-- {
		y = y*x + p[i]
	}
	return y
}

// 求微分多项式
func (p Poly) Reduce() (q Poly) {
	if len(p) <= 1 {
		return Poly{0}
	}
	q = make(Poly, len(p)-1)
	for i, e := range p[1:] {
		q[i] = float64(i+1) * e
	}
	return q
}

// 求积分多项式
func (p Poly) Integral() (q Poly) {
	if len(p) <= 1 {
		return Poly{0}
	}
	q = make(Poly, len(p)+1)
	for i, e := range p {
		q[i+1] = e / float64(i+1)
	}
	return q
}

// 加上一个实数
func (p Poly) AddFloat(k float64) Poly {
	r := make(Poly, len(p))
	copy(r, p)
	r[0] += k
	return r
}

// 减去一个实数
func (p Poly) DecFloat(k float64) Poly {
	r := make(Poly, len(p))
	copy(r, p)
	r[0] -= k
	return r
}

// 乘以一个系数
func (p Poly) MulFloat(k float64) Poly {
	x := len(p)
	r := make(Poly, x)
	for i := 0; i < x; i++ {
		r[i] = p[i] * k
	}
	return r
}

// 除以一个系数
func (p Poly) DivFloat(k float64) Poly {
	x := len(p)
	r := make(Poly, x)
	for i := 0; i < x; i++ {
		r[i] = p[i] / k
	}
	return r
}

// 多项式相加
func (p Poly) Add(q Poly) Poly {
	x, y := len(p), len(q)
	if x < y {
		p, x, q, y = q, y, p, x
	}
	i, r := 0, make(Poly, x)
	for ; i < y; i++ {
		r[i] = p[i] + q[i]
	}
	for ; i < x; i++ {
		r[i] = p[i]
	}
	return r
}

// 多项式相减
func (p Poly) Dec(q Poly) Poly {
	x, y := len(p), len(q)
	if x >= y {
		i, r := 0, make(Poly, x)
		for ; i < y; i++ {
			r[i] = p[i] - q[i]
		}
		for ; i < x; i++ {
			r[i] = p[i]
		}
		return r
	} else {
		i, r := 0, make(Poly, y)
		for ; i < x; i++ {
			r[i] = p[i] - q[i]
		}
		for ; i < y; i++ {
			r[i] = -q[i]
		}
		return r
	}
}

// 多项式相乘
func (p Poly) Mul(q Poly) Poly {
	x, y := len(p), len(q)
	r := make(Poly, x+y-1)
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			r[i+j] += p[i] * q[j]
		}
	}
	return r
}

// 乘以一个系数
func (p Poly) Pow(n uint) Poly {
	if n == 0 {
		return Poly{1}
	}
	r := make(Poly, len(p))
	copy(r, p)
	for ; n > 1; n-- {
		r = r.Mul(p)
	}
	return r
}

// 多项式相除
func (p Poly) Div(q Poly) Poly {
	x, y := len(p), len(q)
	if y > x {
		return Poly{0}
	}
	u := make(Poly, x)
	r := make(Poly, x-y+1)
	copy(u, p)
	for i, j := x-y, x-1; i >= 0; i, j = i-1, j-1 {
		k := u[j] / q[y-1]
		for a, b := j, y-1; b >= 0; a, b = a-1, b-1 {
			u[a] -= k * q[b]
		}
		r[i] = k
	}
	return r
}

// 多项式取余
func (p Poly) Mod(q Poly) Poly {
	x, y := len(p), len(q)
	if y > x {
		return p
	}
	u := make(Poly, x)
	copy(u, p)
	for i := x - 1; i >= y-1; i-- {
		k := u[i] / q[y-1]
		for a, b := i, y-1; b >= 0; a, b = a-1, b-1 {
			u[a] -= k * q[b]
		}
	}
	for i := y - 2; i >= 0; i-- {
		if u[i] != 0 {
			return u[:i+1]
		}
	}
	return Poly{0}
}

// 多项式相除并取余
func (p Poly) DivMod(q Poly) (Poly, Poly) {
	x, y := len(p), len(q)
	if y > x {
		return Poly{0}, p
	}
	u := make(Poly, x)
	r := make(Poly, x-y+1)
	copy(u, p)
	for i, j := x-y, x-1; i >= 0; i, j = i-1, j-1 {
		k := u[j] / q[y-1]
		for a, b := j, y-1; b >= 0; a, b = a-1, b-1 {
			u[a] -= k * q[b]
		}
		r[i] = k
	}
	for i := y - 2; i >= 0; i-- {
		if u[i] != 0 {
			return r, u[:i+1]
		}
	}
	return r, Poly{0}
}

// 求该多项式值为零时的x的取值集合。次数越高，求解越慢。
// 对一元三次方程，求解速度约为马丹诺/范盛金公式的1/3。
func (p Poly) Solve() (s []float64) {
	i, j := 0, len(p)-1
	for i <= j && p[i] == 0 {
		i++
	}
	for j >= 0 && p[j] == 0 {
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
		s = []float64{-p[0] / p[1]}
	case 3:
		a, b, c := p[2], p[1], p[0]
		if k := b*b - 4*a*c; k == 0 {
			s = []float64{-b / (2 * a)}
		} else if k > 0 {
			k = math.Sqrt(k)
			if a > 0 {
				s = []float64{(-b - k) / (2 * a), (-b + k) / (2 * a)}
			} else {
				s = []float64{(-b + k) / (2 * a), (-b - k) / (2 * a)}
			}
		}
	default:
		var a, b float64
		q := p.Reduce()
		n := q.Solve()
		if n == nil {
			a, b = -1, 1
			if q[0] > 0 {
				for ; p.Compute(a) > 0; a *= 10 {
					if math.IsInf(a, -1) {
						goto exit
					}
				}
				for ; p.Compute(b) < 0; b *= 10 {
					if math.IsInf(b, +1) {
						goto exit
					}
				}
			} else {
				for ; p.Compute(a) < 0; a *= 10 {
					if math.IsInf(a, -1) {
						goto exit
					}
				}
				for ; p.Compute(b) > 0; b *= 10 {
					if math.IsInf(b, +1) {
						goto exit
					}
				}
			}
			x, _ := BDsearch(a, b, p.Compute)
			s = append(s, x)
			goto exit
		}
		m := p[len(p)-1] > 0
		l := len(n) - 1
		a, b = n[0], 1
		if m == (len(p)%2 != 0) {
			if p.Compute(a) <= 0 {
				a -= 1
				for p.Compute(a) <= 0 {
					b *= 10
					a -= b
				}
				x, _ := BDsearch(a, n[0], p.Compute)
				s = append(s, x)
			}
		} else {
			if p.Compute(a) >= 0 {
				a -= 1
				for p.Compute(a) >= 0 {
					b *= 10
					a -= b
				}
				x, _ := BDsearch(a, n[0], p.Compute)
				s = append(s, x)
			}
		}
		for i := 0; i < l; i++ {
			if x, e := BDsearch(n[i], n[i+1], p.Compute); e == nil {
				s = append(s, x)
			}
		}
		a, b = 1, n[l]
		if m {
			if p.Compute(b) <= 0 {
				b += 1
				for p.Compute(b) <= 0 {
					a *= 10
					b += a
				}
				x, _ := BDsearch(n[l], b, p.Compute)
				s = append(s, x)
			}
		} else {
			if p.Compute(b) >= 0 {
				b += 1
				for p.Compute(b) >= 0 {
					a *= 10
					b += a
				}
				x, _ := BDsearch(n[l], b, p.Compute)
				s = append(s, x)
			}
		}
		if len(s) <= 1 {
			goto exit
		}
		i := 0
		for _, x := range s[1:] {
			if s[i] != x {
				i++
				s[i] = x
			}
		}
		s = s[:i+1]
	}
exit:
	if i > 0 {
		if s == nil {
			return []float64{0}
		}
		s = append(s, 0)
		for i = len(s) - 2; i >= 0 && s[i] > 0; i-- {
			s[i+1] = s[i]
		}
		s[i+1] = 0
	}
	return s
}

package algebra

import (
	"errors"
	"math"
)

// 切线法求一元方程解，要求（求解区间内：a->b方向）函数单调凸或单调凹
func NTsearch(a, b float64, f, d func(float64) float64) (float64, error) {
	var i, j, x, y float64
	if f(a) == 0 {
		return a, nil
	}
	if y = f(b); y == 0 {
		return b, nil
	}
	for i, j = a-b, 0; ; b, j = x, y {
		if x = b - y/d(b); x == b {
			return x, nil
		}
		if math.IsInf(x, 0) || i*(a-x) < 0 {
			return 0, errors.New("can't find a result at a -> b side")
		}
		if y = f(x); y == 0 {
			return x, nil
		} else if j*y < 0 {
			break
		}
	}
	return BDsearch(b, x, f)
}

// 二分法求一元方程解，最稳定，效率较低
func BNsearch(a, b float64, f func(float64) float64) (float64, error) {
	var x, y float64
	var l, r bool
	if a > b {
		a, b = b, a
	}
	if y = f(a); y == 0 {
		return a, nil
	} else {
		l = (y > 0)
	}
	if y = f(b); y == 0 {
		return b, nil
	} else {
		r = (y > 0)
	}
	if l != r {
		for {
			if x = (a + b) / 2; x == a || x == b {
				return x, nil
			} else if y = f(x); y == 0 {
				return x, nil
			}
			if l == (y > 0) {
				a = x
			} else {
				b = x
			}
		}
	}
	return 0, errors.New("f(a) and f(b) should have different signs")
}

// 折线法求一元方程解，某些情况下效率极低
func DGsearch(a, b float64, f func(float64) float64) (float64, error) {
	var l, r, x, y float64
	if a > b {
		a, b = b, a
	}
	if l = f(a); l == 0 {
		return a, nil
	}
	if r = f(b); r == 0 {
		return b, nil
	}
	if (l > 0) != (r > 0) {
		for {
			if x = (b*l - a*r) / (l - r); x == a || x == b {
				return x, nil
			} else if y = f(x); y == 0 {
				return x, nil
			}
			if (l > 0) == (y > 0) {
				if a < x {
					a, l = x, y
				}
			} else {
				if b > x {
					b, r = x, y
				}
			}
		}
	}
	return 0, errors.New("f(a) and f(b) should have different signs")
}

// 二分法+折线法联合求一元方程解，综合效率最好
func BDsearch(a, b float64, f func(float64) float64) (x float64, e error) {
	var l, r, p, q, y, z float64
	if a > b {
		a, b = b, a
	}
	if l = f(a); l == 0 {
		return a, nil
	}
	if r = f(b); r == 0 {
		return b, nil
	}
	if (l > 0) != (r > 0) {
		for {
			if p = (a + b) / 2; p == a || p == b {
				return p, nil
			}
			if y = f(p); y == 0 {
				return p, nil
			}
			if q = (b*l - a*r) / (l - r); q == a || q == b {
				return q, nil
			}
			if z = f(q); z == 0 {
				return q, nil
			}
			if q <= l || q >= r {
				if (l > 0) == (y > 0) {
					a, l = p, y
				} else {
					b, r = p, y
				}
				continue
			}
			if p > q {
				p, y, q, z = q, z, p, y
			}
			if i := y > 0; i == (z > 0) {
				if i == (l > 0) {
					a, l = q, z
				} else {
					b, r = p, y
				}
			} else {
				a, l, b, r = p, y, q, z
			}
		}
	}
	return 0, errors.New("f(a) and f(b) should have different signs")
}

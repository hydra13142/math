package algebra

import (
	"errors"
	"math"
)

// 二分法+折线法联合求解一元方程，如果区间内有多个解，只能求出其中一个。
func Region(f func(float64) float64, a, b float64) (float64, error) {
	var l, r, p, q, x, y float64
	if a > b {
		a, b = b, a
	}
	if l = f(a); l == 0 {
		return a, nil
	}
	if r = f(b); r == 0 {
		return b, nil
	}
	if m, n := l > 0, r > 0; m != n {
		for {
			if p = (a + b) / 2; p <= a || p >= b {
				return p, nil
			}
			if x = f(p); x == 0 {
				return p, nil
			}
			if q = l/(l-r)*(b-a) + a; q <= a || q >= b {
				return q, nil
			}
			if y = f(q); y == 0 {
				return q, nil
			}
			if p > q {
				p, x, q, y = q, y, p, x
			}
			switch {
			case m != (x > 0):
				b, r = p, x
			case (y > 0) != n:
				a, l = q, y
			default:
				a, l, b, r = p, x, q, y
			}
		}
	}
	return 0, errors.New("Value f(a) & f(b) have the same sign")
}

// 切线法求一元方程解，f为求解函数，k为f的导函数，求解区间为[p,q]，x为迭代起始点。
func Tangent(f, k func(float64) float64, p, q float64, x float64) (float64, error) {
	var y float64
	if p > q {
		p, q = q, p
	}
	for {
		if y = f(x); y == 0 {
			return x, nil
		}
		if x -= y / k(x); x < p || x > q {
			break
		}
	}
	return 0, errors.New("Focal point of tangent line and x-axis outside the region")
}

// 用切线法生成求解区间，再用Region函数求解，要求函数单调递增/递减，且无导数为0的点
func Monotone(f, k func(float64) float64) (float64, error) {
	var x, y float64 = 0, f(0)
	if y == 0 {
		return 0, nil
	}
	sp := y > 0
	for {
		x -= y / k(x)
		if math.IsInf(x, -1) || math.IsInf(x, +1) {
			return 0, errors.New("Failed to make search region")
		}
		if y = f(x); (y > 0) != sp {
			break
		}
	}
	x, e := Region(f, 0, x)
	return x, e
}


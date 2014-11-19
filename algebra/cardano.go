package algebra

import (
	"math"
	"sort"
)

// 卡丹诺公式解一元三次方程
func Cardano(a, b, c, d float64) (ans []float64) {
	defer func() {
		sort.Float64s(ans)
	}()
	if a == 0 {
		if b == 0 {
			return []float64{-d / c}
		}
		switch k := c*c - 4*b*d; {
		case k > 0:
			k = math.Sqrt(k)
			return []float64{(-c - k) / (2 * b), (-c + k) / (2 * b)}
		case k == 0:
			return []float64{-c / (2 * b)}
		default:
			return nil
		}
	}
	k := b / (3 * a)
	p := (c - b*k) / a
	q := (k*(2*b*k-3*c) + 3*d) / (3 * a)
	m := (q * q / 4) + (p * p * p / 27)
	switch {
	case m > 0:
		h := math.Sqrt(m)
		return []float64{math.Cbrt(-q/2-h) + math.Cbrt(-q/2+h) - k}
	case m == 0:
		h := math.Cbrt(-q / 2)
		return []float64{-h - k, 2*h - k}
	default:
		h := math.Sqrt(-m)
		r := math.Cbrt(math.Hypot(h, -q/2))
		g := math.Atan2(h, -q/2) / 3
		i := r * math.Cos(g)
		j := r * math.Sin(g) * math.Sqrt(3)
		return []float64{i + i - k, j - i - k, -j - i - k}
	}
}

// 范盛金公式解一元三次方程
func Shengjin(a, b, c, d float64) (ans []float64) {
	defer func() {
		sort.Float64s(ans)
	}()
	A := b*b - 3*a*c
	B := b*c - 9*a*d
	C := c*c - 3*b*d
	if A == 0 && B == 0 {
		if b == 0 {
			return nil
		} else {
			return []float64{-c / b}
		}
	}
	switch D := B*B - 4*A*C; {
	case D == 0:
		K := B / A
		return []float64{K - b/a, -K / 2}
	case D > 0:
		K := math.Sqrt(D)
		Y1 := A*b + 1.5*a*(-B+K)
		Y2 := A*b + 1.5*a*(-B-K)
		return []float64{-(b + math.Cbrt(Y1) + math.Cbrt(Y2)) / (a * 3)}
	default:
		K := math.Sqrt(A)
		G := math.Acos((2*A*b-3*a*B)/(2*K*A)) / 3
		i := math.Cos(G)
		j := math.Sin(G) * math.Sqrt(3)
		return []float64{-(b + 2*K*i) / (3 * a), (K*(i+j) - b) / (3 * a), (K*(i-j) - b) / (3 * a)}
	}
}

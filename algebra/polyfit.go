package algebra

import (
	"errors"
	"fmt"
)

// 多项式拟合，n为拟合多项式的次数
func UnaryFit(x, y []float64, n int) (Unary, error) {
	i, j, l := len(x), len(y), 0
	if i < j {
		l = i
	} else {
		l = j
	}
	if n++; n <= 0 {
		return nil, errors.New("Illegal input n")
	}
	if n > l {
		return nil, errors.New("Data-set too small")
	}
	group := make([]Linear, n)
	for i := 0; i < n; i++ {
		group[i] = make(Linear, n+1)
	}
	for t := 0; t < l; t++ {
		X, Y := x[t], y[t]
		for i, p := 0, 1.0; i < n; i, p = i+1, p*X {
			for j, q := 0, p; j < n; j, q = j+1, q*X {
				group[i][j] += q
			}
			group[i][n] -= p * Y
		}
	}
	for i:=0;i<n;i++{
		fmt.Println(group[i])
	}
	ans, err := SolveLinearGroup(group...)
	if err != nil {
		return nil, err
	}
	return Unary(ans), nil
}
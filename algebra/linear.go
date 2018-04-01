package algebra

import "errors"

// 多元一次多项式，即多元线性多项式
type Linear []float64

// 提供变量的值以计算多项式的值，a[0]*x+a[1]*y+a[2]*z+...+a[n]
func (this Linear) Compute(x ...float64) (ans float64, err error) {
	n := len(x)
	if n != len(this)-1 {
		return 0, errors.New("Mismatched number of variables")
	}
	for i, t := range x {
		ans += this[i] * t
	}
	return ans + this[n], nil
}

// 求解多元一次方程组
func SolveLinearGroup(p ...Linear) ([]float64, error) {
	var i, j, k, t int
	if len(p) == 0 {
		return nil, errors.New("Find no ploynomials")
	}
	n, l := len(p), len(p[0])
	if n+1 != l {
		return nil, errors.New("Mismatched number of ploynomials")
	}
	for i = 1; i < n; i++ {
		if len(p[i]) != l {
			return nil, errors.New("Mismatched number of variables")
		}
	}
	M := make([]Linear, n)
	S := make([]bool, n)
	for i = 0; i < n; i++ {
		M[i] = make(Linear, l)
		copy(M[i], p[i])
	}
	for i = l - 2; i >= 0; i-- {
		for j = 0; j < n; j++ {
			if M[j][i] != 0 && !S[j] {
				S[j] = true
				break
			}
		}
		if j == n {
			return nil, errors.New("No feasible solution")
		}
		K := M[j][i]
		for k = 0; k < l; k++ {
			M[j][k] /= K
		}
		for k = 0; k < n; k++ {
			if k != j {
				K = M[k][i]
				for t = 0; t < l; t++ {
					M[k][t] -= M[j][t] * K
				}
			}
		}
	}
	ans := make([]float64, n)
	for i = 0; i < n; i++ {
		for j = 0; j < n; j++ {
			if M[j][i] != 0 {
				break
			}
		}
		ans[i] = -M[j][n]
	}
	return ans, nil
}

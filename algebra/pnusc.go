package algebra

import "fmt"

type Pnusc []float64

func (p Pnusc) Compute(x ...float64) (ans float64, err error) {
	if len(x) != len(p)-1 {
		return 0, fmt.Errorf("number of variables is uncorrect")
	}
	for i, t := range x {
		ans += p[i] * t
	}
	return ans + p[len(x)], nil
}

func SolvePnuscGroup(p ...Pnusc) (ans []float64, err error) {
	if len(p) == 0 {
		return nil, fmt.Errorf("there is no ploynomial")
	}
	n, l := len(p), len(p[0])
	i, j, k, t := 0, 0, 0, 0
	if n+1 != l {
		return nil, fmt.Errorf("number of ploynomials is uncorrect")
	}
	for i = 1; i < n; i++ {
		if len(p[i]) != l {
			return nil, fmt.Errorf("ploynomials are of different variables")
		}
	}
	M := make([]Pnusc, n)
	S := make([]bool, n)
	for i = 0; i < n; i++ {
		M[i] = make(Pnusc, l)
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
			return nil, fmt.Errorf("cannot find answer")
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
	ans = make([]float64, n)
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

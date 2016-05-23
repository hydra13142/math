// 假设数据均是已排序为单调增的
package statistic

import (
	"math"
	"sort"
)

var (
	Sort   = sort.Float64s
	Sorted = sort.Float64sAreSorted
)

// 极小值
func Min(list []float64) float64 {
	return list[0]
}

// 极大值
func Max(list []float64) float64 {
	return list[len(list)-1]
}

// 全距，又称极差
func Range(list []float64) float64 {
	return list[len(list)-1] - list[0]
}

// 均值，即平均数
func Average(list []float64) float64 {
	l := len(list)
	s := 0.0
	for j := 0; j < l; j++ {
		s += list[j]
	}
	return s / float64(l)
}

// 方差
func Variance(list []float64, avr float64) float64 {
	l := len(list)
	s := 0.0
	for j := 0; j < l; j++ {
		d := list[j] - avr
		s += d * d
	}
	return s / float64(l)
}

// 标准差
func Standard(list []float64, avr float64) float64 {
	return math.Sqrt(Variance(list, avr))
}

// 中位数
func Middle(list []float64) (mid float64) {
	var n, t, s int
	s = len(list)
	n = s - 1
	t = n / 2
	switch n & 1 {
	case 0:
		mid = list[t]
	case 1:
		mid = (list[t] + list[t+1]) / 2
	}
	return
}

// 四分位数
func Quartile(list []float64) (lft, mid, rgt float64) {
	var n, t, s int
	s = len(list)
	n = s - 3
	t = n / 4
	switch n & 3 {
	case 0:
		lft = list[t]
	case 1:
		lft = (list[t]*1 + list[t+1]*3) / 4
	case 2:
		lft = (list[t]*1 + list[t+1]*1) / 2
	case 3:
		lft = (list[t]*3 + list[t+1]*1) / 4
	}
	n = s - 1
	t = n / 2
	switch n & 1 {
	case 0:
		mid = list[t]
	case 1:
		mid = (list[t] + list[t+1]) / 2
	}
	n = s*3 - 1
	t = n / 4
	switch n & 3 {
	case 0:
		rgt = list[t]
	case 1:
		rgt = (list[t]*1 + list[t+1]*3) / 4
	case 2:
		rgt = (list[t]*1 + list[t+1]*1) / 2
	case 3:
		rgt = (list[t]*3 + list[t+1]*1) / 4
	}
	return
}

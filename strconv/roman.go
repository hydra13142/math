package strconv

import "errors"

// 罗马数字的各位的三元组表示
// 依次为个、十、百位的1、5、10
// 作为一个变量存在，以便用户定制
var Units = [][3]rune{{'I', 'V', 'X'}, {'X', 'L', 'C'}, {'C', 'D', 'M'}}

// 当提供的整数超出表示范围时，FormatRoman会返回本错误
var OutofRange = errors.New("Number Out Of Range")

// 解析一个罗马数字字符串，返回该数字和解析的字节数，如果字节数为0则表示解析失败
func RtoI(s string) (uint, int) {
	sr := []rune(s)
	ls, lu := len(sr), len(Units)
	jw := make([]int, lu+1)
	i, j, p := 0, 0, 0
	if cur := Units[lu-1][2]; sr[p] == cur {
		for ; p < ls && sr[p] == cur; p++ {
			j++
		}
		jw[0] = j
	}
	if p < ls {
		for t := 1; t <= lu; t++ {
			unit := Units[lu-t]
			cur, hlf, abv := unit[0], unit[1], unit[2]
			if sr[p] == cur {
				if p++; p >= ls {
					jw[t] = 1
					break
				}
				if sr[p] == hlf {
					p, jw[t] = p+1, 4
				} else if sr[p] == abv {
					p, jw[t] = p+1, 9
				} else {
					for i, j = 0, 1; i < 2 && p < ls && sr[p] == cur; i, p = i+1, p+1 {
						j++
					}
					jw[t] = j
				}
			} else if sr[p] == hlf {
				if p++; p >= ls {
					jw[t] = 5
					break
				}
				for i, j = 0, 5; i < 3 && p < ls && sr[p] == cur; i, p = i+1, p+1 {
					j++
				}
				jw[t] = j
			}
			if p >= ls {
				break
			}
		}
	}
	if p == 0 {
		return 0, p
	}
	j = 0
	for _, i = range jw {
		j = j*10 + i
	}
	return uint(j), p
}

// 返回一个无符号整数的罗马数字表示和可能的错误
// 如果参数i等于0或者超出表示范围（一般为4000），
// 会返回空字符串和一个错误。
func ItoR(i uint) (string, error) {
	lu := len(Units)
	j, t := 1, 0
	for ; t < lu; t++ {
		j *= 10
	}
	if i >= 4*uint(j) || i == 0 {
		return "", OutofRange
	}
	sr := make([]rune, 0, 4*lu+3)
	cur := Units[lu-1][2]
	for t = int(i) / j; t > 0; t-- {
		sr = append(sr, cur)
	}
	for t, j = 1, j/10; t <= lu; t, j = t+1, j/10 {
		unit := Units[lu-t]
		cur, hlf, abv := unit[0], unit[1], unit[2]
		switch x := int(i) / j % 10; {
		case x == 9:
			sr = append(sr, cur)
			sr = append(sr, abv)
		case x >= 5:
			sr = append(sr, hlf)
			for x -= 5; x > 0; x-- {
				sr = append(sr, cur)
			}
		case x == 4:
			sr = append(sr, cur)
			sr = append(sr, hlf)
		case x <= 3:
			for ; x > 0; x-- {
				sr = append(sr, cur)
			}
		}
	}
	return string(sr), nil
}

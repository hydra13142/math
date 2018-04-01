package strconv

// 罗马数字的各位的三元组表示，依次为个、十、百位的1、5、10
var unit = [][3]byte{{'I', 'V', 'X'}, {'X', 'L', 'C'}, {'C', 'D', 'M'}}

// 返回一个罗马数字字符串表示的数字和可能的错误
func Rtoi(s string) (int, error) {
	i, j, k, l := 0, 0, 0, len(s)
	if l == 0 {
		return 0, ErrEmpty
	}
	t := [4]int{}
	for ; s[i] == 'M'; i++ {
	}
	if i > 3 {
		return 0, ErrSyntax
	}
	t[0] = i
	for k = 2; k >= 0 && i < l; i, k = j, k-1 {
		L, M, H := unit[k][0], unit[k][1], unit[k][2]
		switch s[i] {
		case L:
			for j = i + 1; j < l && s[j] == L; j++ {
			}
			if j-i > 3 {
				return 0, ErrSyntax
			}
			if j-i == 1 && j < l {
				switch s[j] {
				case H:
					t[3-k], j = 9, j+1
				case M:
					t[3-k], j = 4, j+1
				default:
					t[3-k] = 1
				}
			} else {
				t[3-k] = j - i
			}
		case M:
			for j = i + 1; j < l && s[j] == L; j++ {
			}
			if j-i > 4 {
				return 0, ErrSyntax
			}
			t[3-k] = j - i + 4
		default:
			return 0, ErrRune
		}
	}
	if i < l {
		return 0, ErrSyntax
	}
	return t[0]*1000 + t[1]*100 + t[2]*10 + t[3], nil
}

// 返回一个整数的罗马数字字符串表示和可能的错误
func Itor(n int) (string, error) {
	if n >= 4000 || n <= 0 {
		return "", ErrRange
	}
	s := make([]byte, 0, 15)
	for t := int(n) / 1000; t > 0; t-- {
		s = append(s, 'M')
	}
	for t, k := 2, 100; t >= 0; t, k = t-1, k/10 {
		L, M, H := unit[t][0], unit[t][1], unit[t][2]
		switch x := int(n) / k % 10; x {
		case 9:
			s = append(s, L, H)
		case 4:
			s = append(s, L, M)
		case 5, 6, 7, 8:
			s, x = append(s, M), x-5
			fallthrough
		case 1, 2, 3:
			for ; x > 0; x-- {
				s = append(s, L)
			}
		}
	}
	return string(s), nil
}

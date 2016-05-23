// 本包提供汉语数字字符串和整型数值的转换，只支持转换自然数
// 可以与账簿“大写”和日常“小写”两个版本的汉语数字相互转换
package ic

/*
0000 0001 个 1
0000 0010 十 2
0000 0100 百 4
0000 1000 千 8
0001 0000 万 16
0010 0000 亿 32
0011 0000 零 '0'
0011 0001 一 '1'
0011 0010 二 '2'
0011 0011 三 '3'
0011 0100 四 '4'
0011 0101 五 '5'
0011 0110 六 '6'
0011 0111 七 '7'
0011 1000 八 '8'
0011 1001 九 '9'
*/

// 进行千以下的汉字数字与阿拉伯数字的换算，会考虑是位于开头、顺序（亿万、万千）还是跳跃（亿千）的情况
func ctoi_qyx(m byte, str []byte) (v, i int) {
	var q, b, s, g, t, l int
	if len(str) == 0 {
		return 0, 0
	}
	switch m {
	case 0:
		t = 15
	case 1:
		t = 8
	case 2:
		t = 16
	}

	i, l = 0, len(str)
	for t > 0 && i < l {
		switch c := int(str[i]); {
		case c == '0':
			if t&1 != 0 {
				return -1, i
			}
			i, t = i+1, t-1
		case c >= 48:
			i, g = i+1, c&15
			if i < l {
				if c = int(str[i]); c > 8 || c&t == 0 {
					return -1, i
				}
				switch c {
				case 8:
					q, t = g, 4
				case 4:
					b, t = g, 2
				case 2:
					s, t = g, 1
				}
				i, g = i+1, 0
			}
		default:
			return -1, i
		}
	}
	if g == 0 || t&1 != 0 {
		return q*1000 + b*100 + s*10 + g, i
	}
	return -1, i
}

// 进行千万以下的汉字数字与阿拉伯数字的换算，会考虑是位于开头还是顺序（亿万）的情况
func ctoi_wjb(m byte, str []byte) (v, i int) {
	i, l := 0, len(str)
	for i < l && str[i] != 16 {
		i++
	}
	if i < l {
		h, a := ctoi_qyx(m, str[:i])
		if h < 0 {
			return h, a
		}
		i++
		l, b := ctoi_qyx(1, str[i:])
		if l < 0 {
			return l, a + b + 1
		}
		return h*1e4 + l, a + b + 1
	} else if m != 0 {
		return ctoi_qyx(2, str)
	} else {
		return ctoi_qyx(0, str)
	}
}

// 进行万亿以下的汉字数字与阿拉伯数字的换算，只考虑位于开头的情况
func ctoi_yjb(str []byte) (V, L int) {
	i, l := 0, len(str)
	for i < l && str[i] != 32 {
		i++
	}
	if i < l {
		h, a := ctoi_wjb(0, str[:i])
		if h < 0 {
			return h, a
		}
		i++
		l, b := ctoi_wjb(1, str[i:])
		if l < 0 {
			return l, a + b + 1
		}
		return h*1e8 + l, a + b + 1
	} else {
		return ctoi_wjb(0, str)
	}
}

// 汉字数字转换为阿拉伯数字，考虑有单位和无单位的两种情况
// 返回两个整型值，第一个为解析到的数字，另一个为读取的rune数
func CtoI(num string) (int, int) {
	str := make([]byte, 1, 20)
loop:
	for _, r := range num {
		switch r {
		case '〇', '零':
			str = append(str, '0')
		case '一', '壹':
			str = append(str, '1')
		case '二', '贰', '两':
			str = append(str, '2')
		case '三', '叁':
			str = append(str, '3')
		case '四', '肆':
			str = append(str, '4')
		case '五', '伍':
			str = append(str, '5')
		case '六', '陆':
			str = append(str, '6')
		case '七', '柒':
			str = append(str, '7')
		case '八', '捌':
			str = append(str, '8')
		case '九', '玖':
			str = append(str, '9')
		case '十', '拾':
			str = append(str, 2)
		case '百', '佰':
			str = append(str, 4)
		case '千', '仟':
			str = append(str, 8)
		case '万':
			str = append(str, 16)
		case '亿':
			str = append(str, 32)
		default:
			break loop
		}
	}
	for _, e := range str[1:] {
		if e < 48 {
			if str[1] == 2 {
				str[0] = '1'
			} else {
				str = str[1:]
			}
			if i := len(str) - 1; i > 0 {
				if a := str[i]; a >= '1' && a <= '9' {
					if b := str[i-1]; b != '0' && b != 2 {
						str = append(str, str[i-1]>>1)
					}
				}
			}
			return ctoi_yjb(str)
		}
	}
	n := 0
	for _, e := range str {
		n = n*10 + int(e&15)
	}
	return n, len(str)
}

// 进行千以下的阿拉伯数字与汉字数字的换算，会考虑是位于开头、顺序（亿万、万千）还是跳跃（亿千）的情况，返回中间序列
func itoc_qyx(m byte, n int) []byte {
	if n == 0 {
		return []byte{'0'}
	}
	str := make([]byte, 1, 20)
	ws := make([]byte, 4)
	for i := 0; i < 4; i++ {
		ws[i] = byte(n % 10)
		n /= 10
	}
	for i := 3; i >= 0; i-- {
		if ws[i] == 0 {
			str = append(str, '0')
		} else {
			str = append(str, ws[i]+'0', 1<<byte(i))
		}
	}
	l := len(str)
	for i := l - 1; i >= 0 && str[i] == '0'; i-- {
		str[i] = 0
	}
	a := str[0]
	for i := 1; i < l; i++ {
		b := str[i]
		if b == '0' && a == '0' {
			str[i] = 0
		}
		a = b
	}
	switch m {
	case 0:
		i := 1
		for ; str[i] == '0'; i++ {
			str[i] = 0
		}
		if i+1 < len(str) {
			if str[i] == '1' && str[i+1] == 2 {
				str[i] = 0
			}
		}
	case 1:
	case 2:
		if str[1] != '0' {
			str[0] = '0'
		}
	}
	return str
}

// 进行千万以下的阿拉伯数字与汉字数字的换算，会考虑是位于开头还是顺序（亿万）的情况，返回中间序列
func itoc_wjb(m byte, n int) []byte {
	h, l := n/1e4, n%1e4
	if m == 0 {
		if h == 0 {
			return itoc_qyx(0, l)
		} else {
			return append(append(itoc_qyx(0, h), 16), itoc_qyx(1, l)...)
		}
	} else {
		if h == 0 {
			return itoc_qyx(2, l)
		} else {
			return append(append(itoc_qyx(1, h), 16), itoc_qyx(1, l)...)
		}
	}
}

// 进行万亿以下的阿拉伯数字与汉字数字的换算，只考虑位于开头的情况，返回中间序列
func itoc_yjb(n int) []byte {
	h, l := n/1e8, n%1e8
	if h == 0 {
		return itoc_wjb(0, l)
	} else {
		return append(append(itoc_wjb(0, h), 32), itoc_wjb(1, l)...)
	}
}

// 阿拉伯数字转换为汉字数字，考虑普通“小写”和账簿“大写”的两种情况
// 返回一个表示数字的汉语数字字符和该字符串的rune数
func ItoC(n int, big bool) (string, int) {
	if n < 0 || n >= 1e16 {
		return "", 0
	}
	str := itoc_yjb(n)
	num := make([]rune, 0, 40)
	if big {
		for _, c := range str {
			switch c {
			case 0:
			case 2:
				num = append(num, '拾')
			case 4:
				num = append(num, '佰')
			case 8:
				num = append(num, '仟')
			case 16:
				num = append(num, '万')
			case 32:
				num = append(num, '亿')
			case '0':
				num = append(num, '零')
			case '1':
				num = append(num, '壹')
			case '2':
				num = append(num, '贰')
			case '3':
				num = append(num, '叁')
			case '4':
				num = append(num, '肆')
			case '5':
				num = append(num, '伍')
			case '6':
				num = append(num, '陆')
			case '7':
				num = append(num, '柒')
			case '8':
				num = append(num, '捌')
			case '9':
				num = append(num, '玖')
			}
		}
	} else {
		for _, c := range str {
			switch c {
			case 0:
			case 2:
				num = append(num, '十')
			case 4:
				num = append(num, '百')
			case 8:
				num = append(num, '千')
			case 16:
				num = append(num, '万')
			case 32:
				num = append(num, '亿')
			case '0':
				num = append(num, '零')
			case '1':
				num = append(num, '一')
			case '2':
				num = append(num, '二')
			case '3':
				num = append(num, '三')
			case '4':
				num = append(num, '四')
			case '5':
				num = append(num, '五')
			case '6':
				num = append(num, '六')
			case '7':
				num = append(num, '七')
			case '8':
				num = append(num, '八')
			case '9':
				num = append(num, '九')
			}
		}
	}
	return string(num), len(num)
}
